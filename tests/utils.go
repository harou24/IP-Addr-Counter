// File: tests/utils.go
package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// getTestFile returns the absolute path to a test file given its name.
func getTestFile(name string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, "..", "testdata", name), nil
}

// getExpectedUniqueCount runs the Unix pipeline: sort | uniq | wc -l
// and returns the number of unique lines (IP addresses) in the file.
func getExpectedUniqueCount(filename string) (int, error) {
	cmd := exec.Command("sh", "-c", "sort "+filename+" | uniq | wc -l")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	result := strings.TrimSpace(string(output))
	return strconv.Atoi(result)
}
