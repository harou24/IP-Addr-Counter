.PHONY: build run naive bitset test bench clean

BINARY_NAME=ip-addr-counter
MAIN=./cmd/main.go

# Default implementation
IMPL ?= naive

build:
	go build -o $(BINARY_NAME) $(MAIN)

run: build
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make run IMPL=<naive|bitset> FILE=<filename>"; \
		exit 1; \
	fi; \
	./$(BINARY_NAME) $(IMPL) $(FILE)

naive:
	@echo "Running with naive implementation"
	$(MAKE) IMPL=naive run

bitset:
	@echo "Running with bitset implementation"
	$(MAKE) IMPL=bitset run

test:
	go test ./...

bench:
	go test -bench=. -benchmem ./tests

clean:
	rm -f $(BINARY_NAME)
