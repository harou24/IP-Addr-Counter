.PHONY: build run naive bitset concurrent asm test bench clean profile nogc fast

BINARY_NAME=ip-addr-counter

MAIN=./cmd/main.go

build:
ifeq ($(IMPL),asm)
	go build -gcflags=all="-B -l=4 -d=checkptr=0 -wb=0" -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN)
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

asm:
	@echo "Running with assembly implementation"
	$(MAKE) IMPL=asm run

fast:
	@echo "Running with assembly implementation, all disables, and GC off"
	GOGC=off GODEBUG="cgocheck=0,asyncpreemptoff=1,invalidptr=0" $(MAKE) IMPL=asm run

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
