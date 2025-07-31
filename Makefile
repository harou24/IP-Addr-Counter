.PHONY: build run naive best test bench clean

BINARY_NAME=ip-addr-counter

MAIN=cmd/naive.go

build:
	go build -o $(BINARY_NAME) $(MAIN)

run:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make run FILE=<filename>"; \
		exit 1; \
	fi
	./$(BINARY_NAME) $(FILE)

naive:
	@echo "Building and running naive implementation"
	$(MAKE) MAIN=cmd/naive.go build
	$(MAKE) run

best:
	@echo "Best implementation not done yet, skipping..."

test:
	go test ./...

bench:
	go test -bench=. -benchmem ./tests

clean:
	rm -f $(BINARY_NAME)
