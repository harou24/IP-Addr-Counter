//go:build amd64
// +build amd64

package assembly

//go:noescape
func setBitAsmRaw(ptr uintptr, mask uint32) bool

//go:noescape
func ParseIPv4AsmRaw(b []byte) (uint32, bool)
