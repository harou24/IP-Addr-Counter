#include <stdint.h>
#include <string.h>

int parse_ipv4(const char *str, uint32_t *ip) {
    uint32_t parts[4] = {0};
    int part_count = 0;
    int current = 0;
    int digit_count = 0;

    while (*str && part_count < 4) {
        char c = *str++;
        
        if (c >= '0' && c <= '9') {
            // Convert digit and update current part
            current = current * 10 + (c - '0');
            digit_count++;
            
            // Check for invalid part values
            if (digit_count > 3 || current > 255) {
                return 0;
            }
        } else if (c == '.') {
            // Save current part and reset for next part
            if (digit_count == 0) return 0;  // No digits before dot
            parts[part_count++] = current;
            current = 0;
            digit_count = 0;
        } else {
            // Invalid character
            return 0;
        }
    }

    // Handle last part
    if (digit_count == 0) return 0;  // No digits after last dot
    parts[part_count++] = current;

    // Must have exactly 4 parts
    if (part_count != 4) return 0;

    // Combine parts into 32-bit IP
    *ip = (parts[0] << 24) | (parts[1] << 16) | (parts[2] << 8) | parts[3];
    return 1;
}
