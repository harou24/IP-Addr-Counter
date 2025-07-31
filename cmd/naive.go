package main

import (
	"fmt"
	"os"
	"time"

	"IP-Addr-Counter/ipcounter/naive"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]

	start := time.Now()

	counter := naive.New()
	count, err := counter.CountUniqueIPs(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting unique IPs: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(start)

	fmt.Printf("Unique IP addresses: %d\n", count)
	fmt.Printf("Execution time: %s\n", elapsed)
}
