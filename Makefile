.PHONY: build test test-coverage lint fmt lint-fix demo

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
demo: build
	asciinema record --headless --window-size 100x30 --overwrite -i 2 -c "bash demo/demo.sh" demo/demo.cast
	agg --speed 1.4 --idle-time-limit 1.5 demo/demo.cast demo/demo.gif
	rm -f demo/demo.cast
