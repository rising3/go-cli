BIN_NAME ?= mycli

.PHONY: all build test fmt lint clean

all: build test fmt lint

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