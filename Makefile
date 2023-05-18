build:
	go build -ldflags="-s -w" ./...

format:
	gofmt -s -w -l .

test:
	go test -v ./...
