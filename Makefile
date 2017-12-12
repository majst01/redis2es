all: test
	go build 

plugins:
	go build -buildmode=plugin -o lowercase_keys_filter.so filter/lowercase/lowercase_keys_filter.go

test:
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic

dep:
	dep ensure