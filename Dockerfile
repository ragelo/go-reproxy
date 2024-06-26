FROM golang:1.21 as builder

WORKDIR /go/src/github.com/ragelo/go-reproxy

COPY . .

RUN go get -d -v ./...

RUN go build -ldflags "-s -w" -o /go/bin/proxy ./cmd


FROM alpine:3.7

COPY --from=builder /go/bin/proxy /proxy

CMD ["/proxy"]
