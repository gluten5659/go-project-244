.PHONY: build test lint fmt lint-fix

build:
	go build -o bin/gendiff ./cmd/gendiff
test:
	go test -race ./...
lint:
	go tool golangci-lint run
fmt:
	go tool golangci-lint fmt
lint-fix:
	go tool golangci-lint fmt
	go tool golangci-lint run --fix
