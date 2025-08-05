# IP-Addr-Counter
[![Ubuntu CI](https://github.com/harou24/IP-Addr-Counter/actions/workflows/ubuntu.yml/badge.svg)](https://github.com/harou24/IP-Addr-Counter/actions/workflows/ubuntu.yml)
[![macOS CI](https://github.com/harou24/IP-Addr-Counter/actions/workflows/macos.yml/badge.svg)](https://github.com/harou24/IP-Addr-Counter/actions/workflows/macos.yml)

A Go tool to count unique IPv4 addresses from large files efficiently.

This project includes multiple implementations with varying levels of optimization for performance and memory usage:
- **naive**: A basic implementation using a map to track unique IPs as a starting point.
- **bitset**: An efficient single-threaded implementation using a fixed-size bitset (512MB for all possible IPv4 addresses) to mark seen IPs, reducing memory compared to maps.
- **concurrent**: A multi-threaded version of bitset with sharding (divides the bitset into 16384 shards) and atomic updates for thread-safe concurrency, leveraging multiple CPU cores for faster processing on large files.
- **asm**: The most optimized implementation, building on concurrent with assembly-optimized IP parsing and bit operations for lower-level efficiency. Includes compiler flags to disable bounds checks (-B), enable aggressive inlining (-l=4), disable pointer checks (-d=checkptr=0), and disable write barriers (-wb=0) for speed.

Optimizations in "asm" and variants focus on reducing runtime overheads like bounds checking and GC pauses, but they assume well-formed input.


## How to Run

### Prerequisites
- Go 1.24+.
- Test data file.


### Test Data
Sample test data is provided in the repository under the `testdata` folder, including `sample_1M.txt` (100,000 IP addresses), `sample_1M_with_duplicates.txt`, and `sample_35M.txt` (35 million IP addresses). To access the full 120GB dataset, download the archive (`ip_addresses.zip`) as described in the assignment file: [IP-Addr-Counter-GO.md](https://github.com/harou24/IP-Addr-Counter/blob/unsafe/assignment/IP-Addr-Counter-GO.md). Extract the archive and place the contents in the `testdata` folder for testing.


### Makefile Commands

| Command | Description |
|---------|-------------|
| `make naive FILE=<filename>` | Build and run the naive implementation on the given file. |
| `make bitset FILE=<filename>` | Build and run the bitset implementation on the given file. |
| `make concurrent FILE=<filename>` | Build and run the concurrent sharded bitset implementation on the given file. |
| `make asm FILE=<filename>` | Build and run the assembly-optimized implementation (with compiler flags for speed) on the given file. |
| `make fast FILE=<filename>` | Build and run the assembly implementation with maximum disables: GC off, no cgo checks, no async preemption, and no invalid pointer checks (via GODEBUG). Highest risk but potentially fastest for benchmarking. |
| `make profile FILE=<filename>` | Build and run with profiling enabled (generates cpu.prof, mem.prof, goroutine.prof for analysis with `go tool pprof`). |
| `make test` | Run all unit and integration tests. |
| `make bench` | Run benchmarks for all implementations and Unix command comparison. |
| `make clean` | Remove the built binary and profile files. |


## Benchmark Results

This section summarizes performance benchmarks for counting unique IP addresses using different implementations: **Asm** (assembly-optimized), **Bitset** (bitset-based), **Concurrent** (multi-threaded), and **Naive** (baseline map/set approach). Tests were run on macOS with ARM64 architecture and Apple M3 Max CPU, using input files with 1M, 10M, 35M, and 50M lines.

### Metrics
- **Time (ns/op)**: Execution time per operation (lower is better).
- **Memory (B/op)**: Bytes allocated per operation (lower is better).
- **Allocs (allocs/op)**: Heap allocations per operation (lower is better).

> **Note**: The Naive method was tested only on ~1M input for baseline comparison.

### Key Findings
- **Asm** excels in speed across all sizes, with low allocations and balanced memory usage—ideal for high-performance needs.
- **Concurrent** is competitive, leveraging parallelism but slightly slower (20-50% more time) than Asm.
- **Bitset** is memory-efficient for large inputs but scales poorly in time (up to 10x slower for 50M), with high allocations.
- **Naive** is inefficient, suitable only for small datasets due to high memory and allocations.
- Time scales linearly for Asm and Concurrent; Bitset degrades for larger inputs.
- Results may vary on different hardware; profile for your use case.

## Detailed Benchmarks

### 1M Input
| Method      | Time (ns/op) | Memory (B/op) | Allocs/op   |
|-------------|--------------|---------------|-------------|
| Asm         | 47,641,423   | 16,792,178    | 44          |
| Bitset      | 56,415,604   | 16,001,688    | 1,000,004   |
| Concurrent  | 68,283,444   | 16,792,187    | 44          |
| Naive       | 102,052,233  | 54,903,222    | 1,008,221   |

### 10M Input
| Method      | Time (ns/op) | Memory (B/op) | Allocs/op   |
|-------------|--------------|---------------|-------------|
| Asm         | 82,926,061   | 318,849,109   | 88          |
| Concurrent  | 106,129,971  | 318,846,496   | 75          |
| Bitset      | 558,396,854  | 159,980,328   | 10,000,004  |

### 35M Input
| Method      | Time (ns/op) | Memory (B/op) | Allocs/op   |
|-------------|--------------|---------------|-------------|
| Asm         | 196,118,153  | 1,005,468,178 | 197         |
| Concurrent  | 285,219,021  | 1,084,478,862 | 201         |
| Bitset      | 1,952,486,500| 559,933,048   | 35,000,006  |

### 50M Input
| Method      | Time (ns/op) | Memory (B/op) | Allocs/op   |
|-------------|--------------|---------------|-------------|
| Asm         | 259,271,969  | 1,220,824,466 | 249         |
| Concurrent  | 354,116,833  | 1,402,618,557 | 264         |
| Bitset      | 2,800,025,791| 799,928,360   | 50,000,006  |


Commande examples

### Bitset

```
 make bitset FILE=testdata/ip_addresses

Running with bitset implementation
/Library/Developer/CommandLineTools/usr/bin/make IMPL=bitset run
go build -o ip-addr-counter ./cmd/main.go
Implementation: bitset
Unique IP addresses: 1000000000
Execution time: 7m23.139349708s
```

### Concurrent

```
make profile IMPL=concurrent FILE=testdata/ip_addresses 


Implementation: concurrent
Unique IP addresses: 1000000000
Execution time: 1m38.650049584s
Goroutine profile saved to goroutine.prof
Memory profile saved to mem.prof
CPU profile saved to cpu.prof
```

```
make profile IMPL=concurrent FILE=testdata/ip_addresses

Implementation: concurrent
Unique IP addresses: 1000000000
Execution time: 51.644481042s
Goroutine profile saved to goroutine.prof
Memory profile saved to mem.prof
CPU profile saved to cpu.prof
```

### Fastest
```
make profile IMPL=asm FILE=testdata/ip_addresses  
Unique IPs: 1000000000
Time taken: 36.730638416s
```
```
 make asm FILE=testdata/ip_addresses
Running with assembly implementation
/Library/Developer/CommandLineTools/usr/bin/make IMPL=asm run
go build -gcflags=all="-B -l=4" -ldflags="-s -w" -o ip-addr-counter ./cmd/main.go
Unique IPs: 1000000000
Time taken: 34.757187875s
```

```
make asm FILE=testdata/ip_addresses
Running with assembly implementation
/Library/Developer/CommandLineTools/usr/bin/make IMPL=asm run
go build -gcflags=all="-B -l=4 -d=checkptr=0" -ldflags="-s -w" -o ip-addr-counter ./cmd/main.go # Add -pgo=default.pgo if using PGO
Unique IPs: 1000000000
Time taken: 34.599428291s
```

```
make fast FILE=testdata/ip_addresses
Running with assembly implementation, all disables, and GC off
GOGC=off GODEBUG="cgocheck=0,asyncpreemptoff=1,invalidptr=0" /Library/Developer/CommandLineTools/usr/bin/make IMPL=asm run
go build -gcflags=all="-B -l=4 -d=checkptr=0 -wb=0" -ldflags="-s -w" -o ip-addr-counter ./cmd/main.go
Unique IPs: 1000000000
Time taken: 34.262430834s
```
```
make fast FILE=testdata/ip_addresses
Running with assembly implementation, all disables, and GC off
GOGC=off GODEBUG="cgocheck=0,asyncpreemptoff=1,invalidptr=0" /Library/Developer/CommandLineTools/usr/bin/make IMPL=asm run
go build -gcflags=all="-B -l=4 -d=checkptr=0 -wb=0" -ldflags="-s -w" -o ip-addr-counter ./cmd/main.go
Unique IPs: 1000000000
Time taken: 34.128075041s
```

```
make fast FILE=testdata/ip_addresses
Running with assembly implementation, all disables, and GC off
GOGC=off GODEBUG="cgocheck=0,asyncpreemptoff=1,invalidptr=0" /Library/Developer/CommandLineTools/usr/bin/make IMPL=asm run
go build -gcflags=all="-B -l=4 -d=checkptr=0 -wb=0" -ldflags="-s -w" -o ip-addr-counter ./cmd/main.go
Unique IPs: 1000000000
Time taken: 33.885425208s
```
