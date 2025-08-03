package utils

import "testing"

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
		got, err := ParseIPv4Asm([]byte(tt.input))
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseIPv4Asm(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("ParseIPv4Asm(%q) = %08X, want %08X", tt.input, got, tt.want)
		}
	}
}
