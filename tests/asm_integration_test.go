package tests

import (
	"IP-Addr-Counter/ipcounter/assembly"
	"testing"
)

func BenchmarkAsmCountUniqueIPs(b *testing.B) {
	file, err := getTestFile("sample_1M.txt")
	if err != nil {
		b.Fatalf("Failed to get test file: %v", err)
	}
	counter := assembly.New()
	for i := 0; i < b.N; i++ {
		_, err := counter.CountUniqueIPs(file)
		if err != nil {
			b.Fatalf("AsmCounter failed: %v", err)
		}
	}
}

func TestAsmWithSampleData(t *testing.T) {
	file, err := getTestFile("sample_1M.txt")
	if err != nil {
		t.Fatalf("Failed to get test file: %v", err)
	}
	expected, err := getExpectedUniqueCount(file)
	if err != nil {
		t.Fatalf("Failed to get expected count: %v", err)
	}
	counter := assembly.New()
	actual, err := counter.CountUniqueIPs(file)
	if err != nil {
		t.Fatalf("AsmCounter failed: %v", err)
	}
	if expected != int64(actual) {
		t.Errorf("Expected %d unique IPs, got %d", expected, actual)
	}
}

func TestAsmWithDuplicates(t *testing.T) {
	file, err := getTestFile("sample_1M_with_duplicates.txt")
	if err != nil {
		t.Fatalf("Failed to get test file: %v", err)
	}
	expected, err := getExpectedUniqueCount(file)
	if err != nil {
		t.Fatalf("Failed to get expected count: %v", err)
	}
	counter := assembly.New()
	actual, err := counter.CountUniqueIPs(file)
	if err != nil {
		t.Fatalf("AsmCounter failed: %v", err)
	}
	if expected != int64(actual) {
		t.Errorf("Expected %d unique IPs, got %d", expected, actual)
	}
}
