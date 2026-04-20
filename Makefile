.PHONY: tests

tests:
	gotestsum --format=short-verbose ./tests/...

start:
	go run cmd/main.go

lint:
	golangci-lint run