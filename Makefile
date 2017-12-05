all: test
	go build 

plugins:
	go build -buildmode=plugin -o uppercase_keys_filter.so filter/uppercase/uppercase_keys_filter.go

test:
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic

dep:
	dep ensure