all: test
	go build 

test:
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic

dep:
	dep ensure