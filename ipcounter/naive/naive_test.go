package naive_test

import (
	"IP-Addr-Counter/ipcounter/naive"
	"os"
	"testing"
)

func TestCountUniqueIPs(t *testing.T) {
	lines := []string{
		"192.168.0.1",
		"192.168.0.2",
		"192.168.0.1", // duplicate
		"10.0.0.1",
		"invalid.ip",
		"",
		"10.0.0.1", // duplicate
	}

	expectedUnique := 3

	filePath := createTempFile(t, lines)
	defer os.Remove(filePath)

	counter := naive.New()
	count, err := counter.CountUniqueIPs(filePath)
	if err != nil {
		t.Fatalf("CountUniqueIPs failed: %v", err)
	}

	if count != expectedUnique {
		t.Errorf("Expected %d unique IPs, got %d", expectedUnique, count)
	}
}

// createTempFile creates a temporary file with given lines and returns its path.
func createTempFile(t *testing.T, lines []string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test_ips_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	for _, line := range lines {
		_, err := tmpFile.WriteString(line + "\n")
		if err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}
