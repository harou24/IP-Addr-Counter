package assembly

import (
	"errors"
	"unsafe"
)

var errInvalidIP = errors.New("invalid IP")

// setBit in Go (in assembly package)
func setBitAsm(s *shard, offset uint32) bool {
	byteIndex := offset / 8
	bitIndex := offset % 8
	mask := byte(1 << bitIndex)
	wordIndex := byteIndex / 4
	byteOffset := byteIndex % 4
	wordMask := uint32(mask) << (byteOffset * 8)
	ptr := uintptr(unsafe.Pointer(&s.bitset[0])) + uintptr(wordIndex)*4
	return setBitAsmRaw(ptr, wordMask)
}

// ParseIPv4Asm parses an IPv4 address from a byte slice using assembly.
func parseIPv4Asm(b []byte) (uint32, error) {
	ip, ok := ParseIPv4AsmRaw(b)
	if !ok {
		return 0, errInvalidIP
	}
	return ip, nil
}
