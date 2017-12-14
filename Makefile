all: clean test plugins
	go build

plugins:
	mkdir -p lib
	go build -buildmode=plugin -o lib/customer_filter.so filter/customer/customer_filter.go
	go build -buildmode=plugin -o lib/lowercase_keys_filter.so filter/lowercase/lowercase_keys_filter.go

test:
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic

clean:
	rm -rf lib redis2es

dep:
	dep ensure