#include "textflag.h"

// func setBitAsm(ptr uintptr, mask uint32) bool
TEXT ·setBitAsmRaw(SB), NOSPLIT, $0-17
    MOVD ptr+0(FP), R0
    MOVW mask+8(FP), R1

loop:
    LDXRW (R0), R2  // Load exclusive 32-bit
    TSTW R1, R2     // Check if bit already set
    BNE already_set
    ORRW R1, R2, R3 // Set bit
    STXRW R3, (R0), R4  // Store exclusive, status in R4
    CBNZ R4, loop  // Retry if failed
    MOVD $1, R5
    MOVB R5, ret+16(FP)
    RET

already_set:
    MOVB ZR, ret+16(FP)
    RET
