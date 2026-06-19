.PHONY: build test test-coverage lint fmt lint-fix

build:
	go build -o bin/gendiff ./cmd/gendiff
test:
	go test -race ./...
test-coverage:
	go test -race -coverprofile=coverage.out ./...
lint:
	go tool golangci-lint run
fmt:
	go tool golangci-lint fmt
lint-fix:
	go tool golangci-lint run --fix
