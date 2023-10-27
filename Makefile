.PHONY: all build clean deps format lint test tidy

all: tidy format clean build test lint

build: tidy
	go build ./...

clean:
	rm -rfv bin/
	rm -rfv coverage/

deps:
	brew update
	brew upgrade golangci-lint || brew install golangci-lint

format:
	go fmt ./...

lint:
	golangci-lint run --skip-dirs cmd

test:
	go test ./...

test-cover:
	mkdir -p coverage
	go test ./... -covermode=count -coverprofile=coverage/coverage.out

tidy:
	go mod tidy
