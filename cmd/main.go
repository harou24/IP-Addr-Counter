package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"IP-Addr-Counter/ipcounter"
	"IP-Addr-Counter/ipcounter/bitset"
	"IP-Addr-Counter/ipcounter/naive"
)

func getCounter(impl string) (ipcounter.Counter, error) {
	switch strings.ToLower(impl) {
	case "naive":
		return naive.New(), nil
	case "bitset":
		return bitset.New(), nil
	default:
		return nil, fmt.Errorf("unknown implementation: %s", impl)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <implementation> <filename>")
		fmt.Println("Example: go run main.go naive testdata/sample_100k.txt")
		os.Exit(1)
	}
	impl := os.Args[1]
	filename := os.Args[2]

	counter, err := getCounter(impl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	start := time.Now()

	count, err := counter.CountUniqueIPs(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting unique IPs: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(start)

	fmt.Printf("Implementation: %s\n", impl)
	fmt.Printf("Unique IP addresses: %d\n", count)
	fmt.Printf("Execution time: %s\n", elapsed)
}
