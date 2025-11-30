BIN_NAME ?= mycli

.PHONY: all build test fmt lint clean integration-test integration-test-root integration-test-configure integration-test-echo

all: test fmt lint build

build:
	go build -o bin/${BIN_NAME} main.go

test:
	go test ./...

fmt:
	gofmt -s -w .

lint:
	golangci-lint run --enable=govet

clean:
	rm -rf bin/

# Integration tests
integration-test:
	@echo "Running integration tests..."
	@$(MAKE) -C integration_test all

integration-test-root:
	@$(MAKE) -C integration_test test-root

integration-test-configure:
	@$(MAKE) -C integration_test test-configure

integration-test-echo:
	@$(MAKE) -C integration_test test-echo