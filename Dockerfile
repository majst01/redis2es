FROM golang:1.9 AS builder

WORKDIR /go/src/github.com/majst01/redis2es/

COPY Makefile Gopkg.* *.go /go/src/github.com/majst01/redis2es/
RUN go get -u github.com/golang/dep/cmd/dep \
 && make dep all

FROM alpine

COPY --from=builder /go/src/github.com/majst01/redis2es/redis2es /redis2es
COPY --from=builder /go/src/github.com/majst01/redis2es/redis2es/lib /redis2es/lib

CMD ["/redis2es"]
