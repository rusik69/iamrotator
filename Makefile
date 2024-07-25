.PHONY: build get-deps

build:
	go build -o bin/iamrotator cmd/iamrotator/*.go

get:
	go get -v -d ./...

default: get build