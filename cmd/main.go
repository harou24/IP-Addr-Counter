package main

import (
	"IP-Addr-Counter/ipcounter/assembly"
	"IP-Addr-Counter/ipcounter/bitset"
	"IP-Addr-Counter/ipcounter/concurrent"
	"IP-Addr-Counter/ipcounter/naive"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ip-addr-counter <implementation> <filename>")
		fmt.Println("Implementations: naive, bitset, concurrent, assembly")
		os.Exit(1)
	}

	impl := os.Args[1]
	filename := os.Args[2]

	var counter interface {
		CountUniqueIPs(filename string) (int64, error)
	}

	switch impl {
	case "naive":
		counter = &naive.NaiveCounter{}
	case "bitset":
		counter = bitset.New()
	case "concurrent":
		counter = concurrent.New()
	case "asm":
		counter = assembly.New()
	default:
		fmt.Printf("Unknown implementation: %s\n", impl)
		fmt.Println("Implementations: naive, bitset, concurrent, assembly")
		os.Exit(1)
	}

	start := time.Now()
	count, err := counter.CountUniqueIPs(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Unique IPs: %d\n", count)
	fmt.Printf("Time taken: %v\n", time.Since(start))

	if os.Getenv("PPROF") != "" {
		fmt.Println("Profiling enabled; check cpu.prof, mem.prof, or goroutine.prof")
	}
}
