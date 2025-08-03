package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"IP-Addr-Counter/ipcounter"
	"IP-Addr-Counter/ipcounter/bitset"
	"IP-Addr-Counter/ipcounter/concurrent"
	"IP-Addr-Counter/ipcounter/naive"
)

func getCounter(impl string) (ipcounter.Counter, error) {
	switch strings.ToLower(impl) {
	case "naive":
		return naive.New(), nil
	case "bitset":
		return bitset.New(), nil
	case "concurrent":
		return concurrent.New(), nil
	default:
		return nil, fmt.Errorf("unknown implementation: %s", impl)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <implementation> <filename>")
		fmt.Println("Example: go run main.go naive testdata/sample_100k.txt")
		fmt.Println("Implementations: naive, bitset, concurrent")
		os.Exit(1)
	}

	impl := os.Args[1]
	filename := os.Args[2]

	var cpuFile, memFile, goroutineFile *os.File
	var err error

	if os.Getenv("PPROF") == "1" {
		cpuFile, err = os.Create("cpu.prof")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(cpuFile)
		defer func() {
			pprof.StopCPUProfile()
			cpuFile.Close()
			fmt.Println("CPU profile saved to cpu.prof")
		}()

		memFile, err = os.Create("mem.prof")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create memory profile: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			runtime.GC() // get up-to-date stats
			if err := pprof.WriteHeapProfile(memFile); err != nil {
				fmt.Fprintf(os.Stderr, "could not write memory profile: %v\n", err)
			} else {
				fmt.Println("Memory profile saved to mem.prof")
			}
			memFile.Close()
		}()

		goroutineFile, err = os.Create("goroutine.prof")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create goroutine profile: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			if err := pprof.Lookup("goroutine").WriteTo(goroutineFile, 0); err != nil {
				fmt.Fprintf(os.Stderr, "could not write goroutine profile: %v\n", err)
			} else {
				fmt.Println("Goroutine profile saved to goroutine.prof")
			}
			goroutineFile.Close()
		}()
	}

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
