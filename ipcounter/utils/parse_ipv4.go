package utils

/*
#cgo CFLAGS: -O3 -Wall
#include <stdlib.h>       // <--- Add this line for C.free
#include "parse_ipv4.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

func ParseIPv4Cgo(b []byte) (uint32, error) {
	str := C.CString(string(b))
	defer C.free(unsafe.Pointer(str)) // C.free is now defined!

	var ip C.uint32_t
	success := C.parse_ipv4(str, &ip)

	if success == 1 {
		return uint32(ip), nil
	}
	return 0, errors.New("invalid IP")
}
