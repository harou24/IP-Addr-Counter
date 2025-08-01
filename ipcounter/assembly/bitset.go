/*
Package assembly provides an efficient implementation for counting unique IPv4 addresses.
It reads a file containing one IPv4 address per line, processes the file in chunks using
multiple goroutines, and tracks uniqueness with a sharded bitset to minimize memory usage.
Atomic operations ensure thread-safe bitset updates, eliminating lock contention. A sync.Pool
reuses buffers to reduce memory allocation overhead.

Pros:
- Memory-efficient due to bitset usage (512MB for 2^32 IPs, divided across shards).
- High concurrency with multiple worker goroutines, leveraging CPU cores.
- Scales well for large datasets with low contention due to atomic operations.

Cons:
- Chunk copying and I/O may introduce overhead for very large files.
- More complex than naive implementations due to sharding and concurrency.
*/

package assembly

import (
	"IP-Addr-Counter/ipcounter/utils"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"unsafe"
)

// Constants defining configuration for the concurrent implementation.
const (
	maxIPv4       = 1 << 32          // Total number of possible IPv4 addresses (2^32).
	bytesPerChunk = 16 * 1024 * 1024 // Size of each file chunk (16MB) for reading.
	chunkQueueLen = 128              // Buffered channel size for chunk processing.
	numShards     = 16384            // Number of shards to distribute IP addresses.
)

// shard represents a portion of the bitset for storing unique IPs.
type shard struct {
	bitset []byte // Bitset for storing unique IPs in this shard, updated atomically.
}

// BitsetCounter manages a sharded bitset for counting unique IPs.
type BitsetCounter struct {
	shards []*shard // Array of shards, each covering a subset of the IP space.
}

// New initializes a BitsetCounter with pre-allocated shards.
func New() *BitsetCounter {
	// Calculate size of each shard's bitset (2^32 bits / 8 / numShards).
	shardSize := maxIPv4 / 8 / numShards
	shards := make([]*shard, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = &shard{
			bitset: make([]byte, shardSize), // Allocate bitset for this shard.
		}
	}
	return &BitsetCounter{shards: shards}
}

//go:noescape
func setBitAsm(ptr uintptr, mask uint32) bool

// setBit in Go (in assembly package)
func setBit(s *shard, offset uint32) bool {
	byteIndex := offset / 8
	bitIndex := offset % 8
	mask := byte(1 << bitIndex)
	wordIndex := byteIndex / 4
	byteOffset := byteIndex % 4
	wordMask := uint32(mask) << (byteOffset * 8)
	ptr := uintptr(unsafe.Pointer(&s.bitset[0])) + uintptr(wordIndex)*4
	return setBitAsm(ptr, wordMask)
}

// CountUniqueIPs counts unique IPv4 addresses in the specified file.
// It reads the file in chunks, processes them concurrently using multiple goroutines,
// and aggregates the count of unique IPs using a sharded bitset with atomic updates.
func (b *BitsetCounter) CountUniqueIPs(filename string) (int64, error) {
	// Open the input file for reading.
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a buffered reader for efficient file reading.
	reader := bufio.NewReader(file)

	// Channels for distributing chunks to workers and collecting results.
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
	//runtime.GOMAXPROCS(numWorkers)

	// Start worker goroutines to process chunks concurrently.
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chunkChan {
				// Process chunk and count unique IPs.
				count := processChunk(chunk, b)
				resultChan <- count
				// Return buffer to pool for reuse.
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
			total += c // Sum unique IP counts from all chunks.
		}
	}()

	// Read file in chunks and distribute to workers.
	for {
		// Get a buffer from the pool.
		buf := bufPool.Get().([]byte)
		n, err := io.ReadFull(reader, buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if n > 0 {
				// Read until newline to avoid splitting IP addresses.
				rem, _ := reader.ReadBytes('\n')
				buf = append(buf[:n], rem...)
				chunkChan <- buf // Send chunk to workers without copying.
			} else {
				bufPool.Put(buf) // Return unused buffer.
			}
			break
		}
		if err != nil {
			bufPool.Put(buf)
			return 0, fmt.Errorf("read error: %w", err)
		}

		// Extend chunk to include complete IP addresses (until newline).
		rem, _ := reader.ReadBytes('\n')
		if len(rem) > 0 {
			buf = append(buf[:n], rem...)
		} else {
			buf = buf[:n]
		}

		// Send chunk to workers without copying to avoid allocations.
		chunkChan <- buf
	}

	// Close channels and wait for workers to finish.
	close(chunkChan)
	wg.Wait()
	close(resultChan)
	resultWg.Wait()

	return total, nil
}

// processChunk processes a chunk of the input file, parsing IPv4 addresses
// and updating the sharded bitset to count unique IPs using atomic operations.
// Returns the number of new unique IPs found in the chunk.
// processChunk processes a chunk of the input file, parsing IPv4 addresses
// and updating the sharded bitset to count unique IPs using atomic operations.
// Returns the number of new unique IPs found in the chunk.
func processChunk(chunk []byte, b *BitsetCounter) int64 {
	var count int64
	start := 0
	for {
		i := bytes.IndexByte(chunk[start:], '\n')
		if i == -1 {
			break
		}
		line := chunk[start : start+i]
		start += i + 1
		ipInt, _ := utils.ParseIPv4Asm(line)
		shardIdx := ipInt % numShards
		s := b.shards[shardIdx]
		offset := ipInt / numShards
		if setBit(s, offset) {
			count++
		}
	}
	return count
}
