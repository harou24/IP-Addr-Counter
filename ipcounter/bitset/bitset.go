package bitset

import (
	"IP-Addr-Counter/ipcounter/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// BitsetCounter efficiently tracks unique IPv4 addresses using a fixed-size bitset.
type BitsetCounter struct {
	bitset []byte
}

const maxIPv4 = 1 << 32 // 2^32 IPs = 4,294,967,296 bits (512MB)

// New returns a new BitsetCounter with enough space to track every possible IPv4 address.
func New() *BitsetCounter {
	return &BitsetCounter{
		bitset: make([]byte, maxIPv4/8), // 512MB
	}
}

// CountUniqueIPs counts the number of unique IPv4 addresses in the given file.
// It uses a bitset to efficiently track seen addresses.
func (b *BitsetCounter) CountUniqueIPs(filename string) (int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var count int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		ipInt, err := utils.IPToUint32(line)
		if err != nil {
			continue // skip malformed IPs
		}

		byteIndex := ipInt / 8
		bitIndex := ipInt % 8
		mask := byte(1 << bitIndex)

		if b.bitset[byteIndex]&mask == 0 {
			b.bitset[byteIndex] |= mask
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file: %w", err)
	}

	return int64(count), nil
}
