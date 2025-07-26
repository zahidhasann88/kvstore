.PHONY: build run clean test server

build:
	go build -o kvstore .

run:
	go run .

server:
	go run . server

clean:
	rm -f kvstore *.json

test:
	go test ./...

install:
	go install .