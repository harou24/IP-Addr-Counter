package assembly

import (
	"sync"
	"testing"
)

func TestParseIPv4Asm(t *testing.T) {
	tests := []struct {
		input   string
		want    uint32
		wantErr bool
	}{
		{"127.0.0.1", 0x7F000001, false},
		{"0.0.0.0", 0x00000000, false},
		{"255.255.255.255", 0xFFFFFFFF, false},
		{"1.2.3.4", 0x01020304, false},
		{"192.168.1.1", 0xC0A80101, false},
		{"001.002.003.004", 0x01020304, false},
	}

	for _, tt := range tests {
		got, err := parseIPv4Asm([]byte(tt.input))
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseIPv4Asm(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("ParseIPv4Asm(%q) = %08X, want %08X", tt.input, got, tt.want)
		}
	}
}

func TestSetBitAsm(t *testing.T) {
	// Initialize a shard with a small bitset (4 bytes = 32 bits)
	s := &shard{bitset: make([]byte, 4)}

	// Test case 1: Set a bit (offset 0, bit 0 in first byte)
	offset := uint32(0) // First bit in first byte
	got := setBitAsm(s, offset)
	if !got {
		t.Errorf("setBitAsm(s, %d) = false, want true (bit should be set)", offset)
	}
	if s.bitset[0] != 0x01 {
		t.Errorf("setBitAsm(s, %d) did not set bit correctly, got bitset[0] = 0x%02X, want 0x01", offset, s.bitset[0])
	}

	// Test case 2: Set the same bit again (should return false)
	got = setBitAsm(s, offset)
	if got {
		t.Errorf("setBitAsm(s, %d) = true, want false (bit already set)", offset)
	}
	if s.bitset[0] != 0x01 {
		t.Errorf("setBitAsm(s, %d) modified bitset incorrectly, got bitset[0] = 0x%02X, want 0x01", offset, s.bitset[0])
	}

	// Test case 3: Set a different bit (offset 9, bit 1 in second byte)
	offset = uint32(9) // Bit 1 in second byte
	got = setBitAsm(s, offset)
	if !got {
		t.Errorf("setBitAsm(s, %d) = false, want true (bit should be set)", offset)
	}
	if s.bitset[1] != 0x02 {
		t.Errorf("setBitAsm(s, %d) did not set bit correctly, got bitset[1] = 0x%02X, want 0x02", offset, s.bitset[1])
	}

	// Test case 4: Concurrent calls to setBitAsm to verify atomicity
	s = &shard{bitset: make([]byte, 4)} // Reset bitset
	offset = uint32(8)                  // Bit 0 in second byte
	var wg sync.WaitGroup
	successCount := 0
	mu := sync.Mutex{}
	n := 5 // Number of concurrent goroutines
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if setBitAsm(s, offset) {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if successCount != 1 {
		t.Errorf("setBitAsm(s, %d) succeeded %d times, want 1 (atomic operation)", offset, successCount)
	}
	if s.bitset[1] != 0x01 {
		t.Errorf("setBitAsm(s, %d) did not set bit correctly, got bitset[1] = 0x%02X, want 0x01", offset, s.bitset[1])
	}
}
