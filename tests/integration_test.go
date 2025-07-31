package tests

import (
	"IP-Addr-Counter/ipcounter/naive"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// BenchmarkNaiveCountUniqueIPs benchmarks the naive implementation
// by counting unique IPs on the sample 100k dataset.
func BenchmarkNaiveCountUniqueIPs(b *testing.B) {
	wd, err := os.Getwd()
	if err != nil {
		b.Fatalf("Failed to get working directory: %v", err)
	}
	file := filepath.Join(wd, "..", "testdata", "sample_100k.txt")

	counter := naive.New()

	// Run the benchmark b.N times for more stable results
	for i := 0; i < b.N; i++ {
		_, err := counter.CountUniqueIPs(file)
		if err != nil {
			b.Fatalf("NaiveCounter failed: %v", err)
		}
	}
}

// BenchmarkUnixUniqueCount benchmarks the Unix pipeline: sort | uniq | wc -l
// on the sample 100k dataset.
func BenchmarkUnixUniqueCount(b *testing.B) {
	wd, err := os.Getwd()
	if err != nil {
		b.Fatalf("Failed to get working directory: %v", err)
	}
	file := filepath.Join(wd, "..", "testdata", "sample_100k.txt")

	cmdStr := "sort " + file + " | uniq | wc -l"

	for i := 0; i < b.N; i++ {
		cmd := exec.Command("sh", "-c", cmdStr)
		err := cmd.Run()
		if err != nil {
			b.Fatalf("Unix command failed: %v", err)
		}
	}
}

// getExpectedUniqueCount runs the Unix command: sort | uniq | wc -l
// to calculate the number of unique lines (IP addresses) in the given file.
// This serves as the expected value to verify our Go implementation against.
func getExpectedUniqueCount(filename string) (int, error) {
	cmd := exec.Command("sh", "-c", "sort "+filename+" | uniq | wc -l")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	result := strings.TrimSpace(string(output))
	return strconv.Atoi(result)
}

// TestNaiveIntegrationWithSampleData verifies the naive implementation
// counts unique IPs correctly by comparing it against the Unix command output
// on a sample file containing 100k IP addresses without duplicates.
func TestNaiveIntegrationWithSampleData(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	file := filepath.Join(wd, "..", "testdata", "sample_100k.txt")

	expected, err := getExpectedUniqueCount(file)
	if err != nil {
		t.Fatalf("Failed to get expected count using Unix command: %v", err)
	}

	counter := naive.New()
	actual, err := counter.CountUniqueIPs(file)
	if err != nil {
		t.Fatalf("NaiveCounter failed: %v", err)
	}

	if expected != actual {
		t.Errorf("Mismatch: expected %d unique IPs, got %d", expected, actual)
	}
}

// TestNaiveIntegrationWithDuplicates verifies the naive implementation
// correctly handles files containing duplicate IP addresses by comparing
// its output to the Unix command result on a sample file with duplicates.
func TestNaiveIntegrationWithDuplicates(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	file := filepath.Join(wd, "..", "testdata", "sample_100k_with_duplicates.txt")

	expected, err := getExpectedUniqueCount(file)
	if err != nil {
		t.Fatalf("Failed to get expected count using Unix command: %v", err)
	}

	counter := naive.New()
	actual, err := counter.CountUniqueIPs(file)
	if err != nil {
		t.Fatalf("NaiveCounter failed: %v", err)
	}

	if expected != actual {
		t.Errorf("Mismatch (with duplicates): expected %d unique IPs, got %d", expected, actual)
	}
}
