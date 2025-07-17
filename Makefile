.PHONY: test clean

clean:
	go clean -cache -testcache

build:
	go build ./...

test: clean
	go test -race -count=1 ./...
