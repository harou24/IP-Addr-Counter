// File: tests/naive_integration_test.go
package tests

import (
	"IP-Addr-Counter/ipcounter/naive"
	"testing"
)

func BenchmarkNaiveCountUniqueIPs(b *testing.B) {
	file, err := getTestFile("sample_100k.txt")
	if err != nil {
		b.Fatalf("Failed to get test file: %v", err)
	}

	counter := naive.New()
	for i := 0; i < b.N; i++ {
		_, err := counter.CountUniqueIPs(file)
		if err != nil {
			b.Fatalf("NaiveCounter failed: %v", err)
		}
	}
}

func TestNaiveIntegrationWithSampleData(t *testing.T) {
	file, err := getTestFile("sample_100k.txt")
	if err != nil {
		t.Fatalf("Failed to get test file: %v", err)
	}

	expected, err := getExpectedUniqueCount(file)
	if err != nil {
		t.Fatalf("Failed to get expected count: %v", err)
	}

	counter := naive.New()
	actual, err := counter.CountUniqueIPs(file)
	if err != nil {
		t.Fatalf("NaiveCounter failed: %v", err)
	}

	if expected != actual {
		t.Errorf("Mismatch: expected %d, got %d", expected, actual)
	}
}

func TestNaiveIntegrationWithDuplicates(t *testing.T) {
	file, err := getTestFile("sample_100k_with_duplicates.txt")
	if err != nil {
		t.Fatalf("Failed to get test file: %v", err)
	}

	expected, err := getExpectedUniqueCount(file)
	if err != nil {
		t.Fatalf("Failed to get expected count: %v", err)
	}

	counter := naive.New()
	actual, err := counter.CountUniqueIPs(file)
	if err != nil {
		t.Fatalf("NaiveCounter failed: %v", err)
	}

	if expected != actual {
		t.Errorf("Mismatch (with duplicates): expected %d, got %d", expected, actual)
	}
}
