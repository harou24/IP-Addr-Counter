package naive

/*
Naive implementation of unique IPv4 address counting.

This program reads a file containing one IPv4 address per line,
parses each IP into a 32-bit integer, and stores it in a map to track uniqueness.

Pros:
- Simple and easy to understand.
- Correctly counts unique IPs for small to medium-sized files.

Cons:
- Uses a map that grows with the number of unique IPs,
  which can consume a large amount of memory for very big files.
- Not suitable for extremely large datasets.
*/

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type NaiveCounter struct{}

func New() *NaiveCounter {
	return &NaiveCounter{}
}

func (c *NaiveCounter) CountUniqueIPs(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	uniqueIPs := make(map[uint32]struct{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		ipInt, err := ipToUint32(line)
		if err != nil {
			continue
		}
		uniqueIPs[ipInt] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file: %w", err)
	}

	return len(uniqueIPs), nil
}

// ipToUint32 converts a string IPv4 address
// into its 32-bit unsigned integer representation.
func ipToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return 0, fmt.Errorf("invalid IPv4 address: %s", ipStr)
	}
	// Shift each byte to its position and combine into one uint32
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3]), nil
}
