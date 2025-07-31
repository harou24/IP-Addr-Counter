.PHONY: build run naive bitset concurrent test bench clean profile

BINARY_NAME=ip-addr-counter
MAIN=./cmd/main.go

build:
ifeq ($(IMPL),cgo)
	CGO_ENABLED=1 go build -o $(BINARY_NAME) $(MAIN)
else
	go build -o $(BINARY_NAME) $(MAIN)
endif

run: build
	@if [ -z "$(FILE)" ]; then \
	    echo "Usage: make run FILE=<filename>"; \
	    exit 1; \
	fi; \
	./$(BINARY_NAME) $(IMPL) $(FILE)

naive:
	@echo "Running with naive implementation"
	$(MAKE) IMPL=naive run

bitset:
	@echo "Running with bitset implementation"
	$(MAKE) IMPL=bitset run

concurrent:
	@echo "Running with concurrent bitset implementation"
	$(MAKE) IMPL=concurrent run

profile:
	@if [ -z "$(FILE)" ]; then \
	    echo "Usage: make profile FILE=<filename>"; \
	    exit 1; \
	fi; \
	PPROF=1 ./$(BINARY_NAME) $(IMPL) $(FILE)

test:
	go test ./...

bench:
	go test -bench=. -benchmem ./tests

clean:
	rm -f $(BINARY_NAME) cpu.prof mem.prof goroutine.prof
