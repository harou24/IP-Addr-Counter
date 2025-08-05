package utils

import (
	"errors"
	"testing"
)

func TestParseIPv4(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint32
		wantErr error
	}{
		{name: "valid IP", input: "192.168.1.1", want: 0xC0A80101, wantErr: nil},
		{name: "valid IP - padded", input: "001.002.003.004", want: 0x01020304, wantErr: nil},
		{name: "invalid octet", input: "1.2.256.4", want: 0, wantErr: errors.New("invalid octet")},
		{name: "invalid digit", input: "1.a.3.4", want: 0, wantErr: errors.New("invalid digit")},
		{name: "extra data", input: "1.2.3.4.5", want: 0, wantErr: errors.New("extra data")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIPv4([]byte(tt.input))
			if err != nil && (tt.wantErr == nil || err.Error() != tt.wantErr.Error()) || err == nil && tt.wantErr != nil {
				t.Errorf("ParseIPv4(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseIPv4(%q) = %08X, want %08X", tt.input, got, tt.want)
			}
		})
	}
}
