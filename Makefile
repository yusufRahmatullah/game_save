.PHONY: build dep compress pretty short-test test

build: pretty
	go build -ldflags="-s -w" -o gamesave gamesave.go

dep:
	go mod tidy

compress: build
	upx --brute gamesave

pretty:
	gofmt -w *.go

short-test: pretty
	go test -short ./...

test: pretty
	go test ./...

