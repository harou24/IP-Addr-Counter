# IP-Addr-Counter
[![Ubuntu CI](https://github.com/harou24/IP-Addr-Counter/actions/workflows/ubuntu.yml/badge.svg)](https://github.com/harou24/IP-Addr-Counter/actions/workflows/ubuntu.yml)
[![macOS CI](https://github.com/harou24/IP-Addr-Counter/actions/workflows/macos.yml/badge.svg)](https://github.com/harou24/IP-Addr-Counter/actions/workflows/macos.yml)

A Go tool to count unique IPv4 addresses from large files efficiently.

This project includes multiple implementations with varying levels of optimization for performance and memory usage:
- **naive**: A basic implementation using a map to track unique IPs. Simple but memory-intensive for large datasets (uses ~GBs for 1B IPs).
- **bitset**: An efficient single-threaded implementation using a fixed-size bitset (512MB for all possible IPv4 addresses) to mark seen IPs, reducing memory compared to maps.
- **concurrent**: A multi-threaded version of bitset with sharding (divides the bitset into 16384 shards) and atomic updates for thread-safe concurrency, leveraging multiple CPU cores for faster processing on large files.
- **asm**: The most optimized implementation, building on concurrent with assembly-optimized IP parsing and bit operations for lower-level efficiency. Includes compiler flags to disable bounds checks (-B), enable aggressive inlining (-l=4), disable pointer checks (-d=checkptr=0), and disable write barriers (-wb=0) for speed.

Optimizations in "asm" and variants focus on reducing runtime overheads like bounds checking and GC pauses, but they assume well-formed input and can increase crash risk if misused—use for benchmarking.

## How to Run

### Prerequisites
- Go 1.24+.
- Test data file.

# Makefile Commands


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


## Examples

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
