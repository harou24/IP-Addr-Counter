#include "textflag.h"

// func ParseIPv4AsmRaw(b []byte) (uint32, bool)
TEXT ·ParseIPv4AsmRaw(SB), NOSPLIT, $0-32
    MOVD b+0(FP), R0  // b.ptr
    MOVD b+8(FP), R1  // b.len
    MOVD $0, R2       // ip = 0
    MOVD $0, R3       // part = 0
    MOVD $10, R6      // constant 10
    MOVD $255, R7     // max octet (unused, assume <=255)
    MOVD $'.', R8     // '.' (unused, assume dot)
    MOVD $'0', R9     // '0'
    MOVD $'9', R10    // '9' (unused, but for CMP)

// Octet 1
    MOVBU (R0), R5
    SUB R9, R5, R3
    ADD $1, R0
    SUB $1, R1
    MOVBU (R0), R5
    SUB R9, R5, R14
    CMP $0, R14
    BLT dot1
    CMP $9, R14
    BHI dot1
    MUL R6, R3, R3
    ADD R14, R3
    ADD $1, R0
    SUB $1, R1
    MOVBU (R0), R5
    SUB R9, R5, R14
    CMP $0, R14
    BLT dot1
    CMP $9, R14
    BHI dot1
    MUL R6, R3, R3
    ADD R14, R3
    ADD $1, R0
    SUB $1, R1
dot1:
    ADD $1, R0
    SUB $1, R1
    MOVD R3, R2
    MOVD $0, R3

// Octet 2
    MOVBU (R0), R5
    SUB R9, R5, R3
    ADD $1, R0
    SUB $1, R1
    MOVBU (R0), R5
    SUB R9, R5, R14
    CMP $0, R14
    BLT dot2
    CMP $9, R14
    BHI dot2
    MUL R6, R3, R3
    ADD R14, R3
    ADD $1, R0
    SUB $1, R1
    MOVBU (R0), R5
    SUB R9, R5, R14
    CMP $0, R14
    BLT dot2
    CMP $9, R14
    BHI dot2
    MUL R6, R3, R3
    ADD R14, R3
    ADD $1, R0
    SUB $1, R1
dot2:
    ADD $1, R0
    SUB $1, R1
    LSL $8, R2, R2
    ORR R3, R2, R2
    MOVD $0, R3

// Octet 3
    MOVBU (R0), R5
    SUB R9, R5, R3
    ADD $1, R0
    SUB $1, R1
    MOVBU (R0), R5
    SUB R9, R5, R14
    CMP $0, R14
    BLT dot3
    CMP $9, R14
    BHI dot3
    MUL R6, R3, R3
    ADD R14, R3
    ADD $1, R0
    SUB $1, R1
    MOVBU (R0), R5
    SUB R9, R5, R14
    CMP $0, R14
    BLT dot3
    CMP $9, R14
    BHI dot3
    MUL R6, R3, R3
    ADD R14, R3
    ADD $1, R0
    SUB $1, R1
dot3:
    ADD $1, R0
    SUB $1, R1
    LSL $8, R2, R2
    ORR R3, R2, R2
    MOVD $0, R3

// Octet 4
    MOVBU (R0), R5
    SUB R9, R5, R3
    ADD $1, R0
    SUB $1, R1
    CBZ R1, finish
    MOVBU (R0), R5
    SUB R9, R5, R14
    CMP $0, R14
    BLT finish
    CMP $9, R14
    BHI finish
    MUL R6, R3, R3
    ADD R14, R3
    ADD $1, R0
    SUB $1, R1
    CBZ R1, finish
    MOVBU (R0), R5
    SUB R9, R5, R14
    CMP $0, R14
    BLT finish
    CMP $9, R14
    BHI finish
    MUL R6, R3, R3
    ADD R14, R3
    ADD $1, R0
    SUB $1, R1

finish:
    LSL $8, R2, R2
    ORR R3, R2, R2
    MOVW R2, ret+24(FP)
    MOVD $1, R5
    MOVB R5, ret+28(FP)
    RET

invalid:
    MOVW ZR, ret+24(FP)
    MOVB ZR, ret+28(FP)
    RET
