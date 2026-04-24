.PHONY: tests
include .env
export

tests:
	gotestsum --format=short-verbose ./tests/...

start:
	go run cmd/main.go

lint:
	golangci-lint run

build:
	GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) ./cmd

deploy: build
	scp $(APP_NAME) $(USER)@$(HOST):/home/user/

run: deploy
	ssh $(USER)@$(HOST) "./$(APP_NAME)"