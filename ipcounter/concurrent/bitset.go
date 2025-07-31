package concurrent

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
	maxIPv4       = 1 << 32
	bytesPerChunk = 4 * 1024 * 1024
	chunkQueueLen = 128
	numShards     = 4096 // Fine-grained shards for minimal lock contention
)

type shard struct {
	mu     sync.Mutex
	bitset []byte
}

type BitsetCounter struct {
	shards []*shard
}

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

func (b *BitsetCounter) CountUniqueIPs(filename string) (int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	chunkChan := make(chan []byte, chunkQueueLen)
	resultChan := make(chan int64, chunkQueueLen)

	bufPool := sync.Pool{
		New: func() interface{} {
			return make([]byte, bytesPerChunk)
		},
	}

	numWorkers := runtime.NumCPU()
	runtime.GOMAXPROCS(numWorkers)

	var wg sync.WaitGroup
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

	var resultWg sync.WaitGroup
	resultWg.Add(1)
	var total int64
	go func() {
		defer resultWg.Done()
		for c := range resultChan {
			total += c
		}
	}()

	// Feed chunks
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
				chunkChan <- chunkCopy
			} else {
				bufPool.Put(buf)
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
			ipInt, err := utils.IPToUint32(string(line))
			if err != nil {
				continue
			}

			shardIdx := ipInt % numShards
			s := b.shards[shardIdx]
			offset := ipInt / numShards
			byteIndex := offset / 8
			bitIndex := offset % 8
			mask := byte(1 << bitIndex)

			s.mu.Lock()
			if s.bitset[byteIndex]&mask == 0 {
				s.bitset[byteIndex] |= mask
				count++
			}
			s.mu.Unlock()
		}
	}
	return count
}
