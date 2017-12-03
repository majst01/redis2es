FROM golang:1.9 AS builder

WORKDIR /go/src/github.com/majst01/redis-to-elastic/

COPY Makefile Gopkg.* *.go /go/src/github.com/majst01/redis-to-elastic/
RUN go get -u github.com/golang/dep/cmd/dep \
 && make dep all

FROM alpine

COPY --from=builder /go/src/github.com/majst01/redis-to-elastic/redis-to-elastic /redis-to-elastic

CMD ["/redis-to-elastic"]
