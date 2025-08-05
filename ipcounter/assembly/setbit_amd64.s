#include "textflag.h"

// func setBitAsmRaw(ptr uintptr, mask uint32) bool
TEXT ·setBitAsmRaw(SB), NOSPLIT, $0-17
    MOVQ ptr+0(FP), DI    // DI = ptr
    MOVL mask+8(FP), BX   // BX = mask
    XORL AX, AX           // AX = old value (for CMPXCHG)

loop:
    MOVL 0(DI), CX        // CX = current value at ptr
    TESTL BX, CX          // Check if bit is already set
    JNZ already_set       // If set, return false
    ORL BX, CX            // CX = current | mask (new value)
    LOCK                  // Atomic operation
    CMPXCHGL CX, 0(DI)    // Compare and swap: if 0(DI) == AX, set 0(DI) = CX
    JNZ loop              // Retry if CMPXCHG failed
    MOVB $1, ret+16(FP)   // Return true
    RET

already_set:
    MOVB $0, ret+16(FP)   // Return false
    RET
