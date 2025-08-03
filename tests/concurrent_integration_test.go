package tests

import (
	"IP-Addr-Counter/ipcounter/concurrent"
	"testing"
)

func BenchmarkConcurrentCountUniqueIPs(b *testing.B) {
	file, err := getTestFile("sample_100k.txt")
	if err != nil {
		b.Fatalf("Failed to get test file: %v", err)
	}
	counter := concurrent.New()
	for i := 0; i < b.N; i++ {
		_, err := counter.CountUniqueIPs(file)
		if err != nil {
			b.Fatalf("ConcurrentCounter failed: %v", err)
		}
	}
}

func TestConcurrentWithSampleData(t *testing.T) {
	file, err := getTestFile("sample_100k.txt")
	if err != nil {
		t.Fatalf("Failed to get test file: %v", err)
	}
	expected, err := getExpectedUniqueCount(file)
	if err != nil {
		t.Fatalf("Failed to get expected count: %v", err)
	}
	counter := concurrent.New()
	actual, err := counter.CountUniqueIPs(file)
	if err != nil {
		t.Fatalf("ConcurrentCounter failed: %v", err)
	}
	if expected != int64(actual) {
		t.Errorf("Expected %d unique IPs, got %d", expected, actual)
	}
}

func TestConcurrentWithDuplicates(t *testing.T) {
	file, err := getTestFile("sample_100k_with_duplicates.txt")
	if err != nil {
		t.Fatalf("Failed to get test file: %v", err)
	}
	expected, err := getExpectedUniqueCount(file)
	if err != nil {
		t.Fatalf("Failed to get expected count: %v", err)
	}
	counter := concurrent.New()
	actual, err := counter.CountUniqueIPs(file)
	if err != nil {
		t.Fatalf("ConcurrentCounter failed: %v", err)
	}
	if expected != int64(actual) {
		t.Errorf("Expected %d unique IPs, got %d", expected, actual)
	}
}
