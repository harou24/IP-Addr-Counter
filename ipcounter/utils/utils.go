package utils

import (
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
