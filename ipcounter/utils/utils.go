package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// IPToUint32 converts a string IPv4 address into a 32-bit integer.
func IPToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(strings.TrimSpace(ipStr)).To4()
	if ip == nil {
		return 0, fmt.Errorf("invalid IPv4 address: %s", ipStr)
	}
	// Shift each byte to its position and combine into one uint32
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3]), nil
}

// ParseIPv4 parses each octet directly, avoiding string conversions and net.ParseIP overhead.
// Returns an error for invalid inputs
func ParseIPv4(b []byte) (uint32, error) {
	var ip, part uint32
	count := 0
	for _, c := range b {
		switch {
		case c >= '0' && c <= '9':
			part = part*10 + uint32(c-'0')
			if part > 255 {
				return 0, errors.New("invalid octet")
			}
		case c == '.':
			if count == 3 {
				return 0, errors.New("too many dots")
			}
			ip = (ip << 8) | part
			part = 0
			count++
		default:
			return 0, errors.New("invalid character")
		}
	}
	if count != 3 {
		return 0, errors.New("too few dots")
	}
	return (ip << 8) | part, nil
}
