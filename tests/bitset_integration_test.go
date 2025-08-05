package tests

import (
	"IP-Addr-Counter/ipcounter/bitset"
	"testing"
)

func BenchmarkBitsetCountUniqueIPs(b *testing.B) {
	file, err := getTestFile("sample_1M.txt")
	if err != nil {
		b.Fatalf("Failed to get test file: %v", err)
	}

	counter := bitset.New()
	for i := 0; i < b.N; i++ {
		_, err := counter.CountUniqueIPs(file)
		if err != nil {
			b.Fatalf("BitsetCounter failed: %v", err)
		}
	}
}

func TestBitsetWithSampleData(t *testing.T) {
	file, err := getTestFile("sample_1M.txt")
	if err != nil {
		t.Fatalf("Failed to get test file: %v", err)
	}

	expected, err := getExpectedUniqueCount(file)
	if err != nil {
		t.Fatalf("Failed to get expected count: %v", err)
	}

	counter := bitset.New()
	actual, err := counter.CountUniqueIPs(file)
	if err != nil {
		t.Fatalf("BitsetCounter failed: %v", err)
	}

	if expected != int64(actual) {
		t.Errorf("Expected %d unique IPs, got %d", expected, actual)
	}
}

func TestBitsetWithDuplicates(t *testing.T) {
	file, err := getTestFile("sample_1M_with_duplicates.txt")
	if err != nil {
		t.Fatalf("Failed to get test file: %v", err)
	}

	expected, err := getExpectedUniqueCount(file)
	if err != nil {
		t.Fatalf("Failed to get expected count: %v", err)
	}

	counter := bitset.New()
	actual, err := counter.CountUniqueIPs(file)
	if err != nil {
		t.Fatalf("BitsetCounter failed: %v", err)
	}

	if expected != int64(actual) {
		t.Errorf("Expected %d unique IPs, got %d", expected, actual)
	}
}
