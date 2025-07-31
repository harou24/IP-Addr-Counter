# IP-Addr-Counter

A Go tool to count unique IPv4 addresses from large files efficiently.

This project includes multiple implementations, starting with a naive map-based approach, with plans for optimized versions to improve performance and memory usage.

## How to run

### Makefile Commands

- `make naive FILE=<filename>`  
  Build and run the naive implementation on the given file.

- `make test`  
  Run all unit and integration tests.

- `make bench`  
  Run benchmarks for all implementations and Unix command comparison.

- `make clean`  
  Remove the built binary.