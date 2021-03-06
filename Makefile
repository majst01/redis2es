all: clean dep test plugins
	go build

plugins:
	mkdir -p lib
	go build -buildmode=plugin -o lib/catchall_filter.so filter/catchall/catchall_filter.go
	go build -buildmode=plugin -o lib/customer_filter.so filter/customer/customer_filter.go
	go build -buildmode=plugin -o lib/lowercase_keys_filter.so filter/lowercase/lowercase_keys_filter.go

test:
	go test -v -race -cover $(shell go list ./...)

clean:
	rm -rf lib redis2es

dep:
	dep ensure

docker:
	docker build --no-cache -t majst01/redis2es .