package concurrent

/*
Package concurrent provides an efficient implementation for counting unique IPv4 addresses.

It reads a file containing one IPv4 address per line, processes the file in chunks using
multiple goroutines, and tracks uniqueness with a sharded bitset to minimize memory usage
and lock contention. Each shard uses a mutex for thread-safe updates, and a sync.Pool
reuses buffers to reduce memory allocation overhead.

Pros:
- Memory-efficient due to bitset usage (512MB for 2^32 IPs, divided across shards).
- High concurrency with multiple worker goroutines, leveraging CPU cores.
- Scales well for large datasets.

Cons:
- Mutex contention across shards can bottleneck performance for high worker counts.
- Chunk copying and I/O may introduce overhead for very large files.
*/

import (
	"IP-Addr-Counter/ipcounter/utils"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

const (
	maxIPv4       = 1 << 32         // Total number of possible IPv4 addresses (2^32).
	bytesPerChunk = 4 * 1024 * 1024 // Size of each file chunk (4MB) for reading.
	chunkQueueLen = 128             // Buffered channel size for chunk processing.
	numShards     = 4096            // Number of shards to reduce mutex contention.
)

// shard represents a portion of the bitset with a mutex for thread-safe updates.
type shard struct {
	mu     sync.Mutex // Protects bitset from concurrent access.
	bitset []byte     // Bitset for storing unique IPs in this shard.
}

// BitsetCounter manages a sharded bitset for counting unique IPs.
type BitsetCounter struct {
	shards []*shard // Array of shards, each covering a subset of the IP space.
}

// New initializes a BitsetCounter with pre-allocated shards.
func New() *BitsetCounter {
	shards := make([]*shard, numShards)
	shardSize := maxIPv4 / 8 / numShards
	for i := 0; i < numShards; i++ {
		shards[i] = &shard{
			bitset: make([]byte, shardSize),
		}
	}
	return &BitsetCounter{shards: shards}
}

// CountUniqueIPs counts unique IPv4 addresses in the specified file.
// It reads the file in chunks, processes them concurrently using multiple goroutines,
// and aggregates the count of unique IPs using a sharded bitset.
func (b *BitsetCounter) CountUniqueIPs(filename string) (int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	chunkChan := make(chan []byte, chunkQueueLen)
	resultChan := make(chan int64, chunkQueueLen)

	// Initialize a sync.Pool to reuse chunk buffers and reduce allocations.
	bufPool := sync.Pool{
		New: func() interface{} {
			return make([]byte, bytesPerChunk)
		},
	}

	// Set number of workers to CPU core count for optimal parallelism.
	numWorkers := runtime.NumCPU()
	runtime.GOMAXPROCS(numWorkers)

	var wg sync.WaitGroup
	// Start worker goroutines to process chunks concurrently.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chunkChan {
				count := processChunk(chunk, b)
				resultChan <- count
				bufPool.Put(chunk)
			}
		}()
	}

	// Start a goroutine to aggregate results from workers.
	var resultWg sync.WaitGroup
	resultWg.Add(1)
	var total int64
	go func() {
		defer resultWg.Done()
		for c := range resultChan {
			total += c
		}
	}()

	// Read file in chunks and distribute to workers.
	for {
		buf := bufPool.Get().([]byte)
		n, err := io.ReadFull(reader, buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if n > 0 {
				// Read until newline to complete last line
				rem, _ := reader.ReadBytes('\n')
				buf = append(buf[:n], rem...)
				chunkCopy := make([]byte, len(buf))
				copy(chunkCopy, buf)
				chunkChan <- chunkCopy // Send chunk to workers.
			} else {
				bufPool.Put(buf) // Return unused buffer.
			}
			break
		}
		if err != nil {
			bufPool.Put(buf)
			return 0, fmt.Errorf("read error: %w", err)
		}

		// Extend chunk until newline to not split IPs
		rem, _ := reader.ReadBytes('\n')
		if len(rem) > 0 {
			buf = append(buf[:n], rem...)
		} else {
			buf = buf[:n]
		}

		// Create a copy of the chunk to avoid data races in workers.
		chunkCopy := make([]byte, len(buf))
		copy(chunkCopy, buf)
		chunkChan <- chunkCopy
	}

	close(chunkChan)
	wg.Wait()
	close(resultChan)
	resultWg.Wait()

	return total, nil
}

// processChunk processes a chunk of the input file, parsing IPv4 addresses
// and updating the sharded bitset to count unique IPs.
// Returns the number of new unique IPs found in the chunk.
func processChunk(chunk []byte, b *BitsetCounter) int64 {
	var count int64
	start := 0
	for i, c := range chunk {
		if c == '\n' {
			line := bytes.TrimSpace(chunk[start:i])
			start = i + 1
			if len(line) == 0 {
				continue
			}
			// Parse IP address to uint32 using optimized byte-based parser.
			ipInt, err := utils.ParseIPv4(line)
			if err != nil {
				continue
			}

			// Determine shard and bit position for the IP.
			shardIdx := ipInt % numShards
			s := b.shards[shardIdx]
			offset := ipInt / numShards
			byteIndex := offset / 8
			bitIndex := offset % 8
			mask := byte(1 << bitIndex)

			// Update bitset under lock to ensure thread safety.
			s.mu.Lock()
			if s.bitset[byteIndex]&mask == 0 {
				s.bitset[byteIndex] |= mask // Mark IP as seen.
				count++
			}
			s.mu.Unlock()
		}
	}
	return count
}
