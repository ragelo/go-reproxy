FROM golang:1.21 as builder

WORKDIR /go/src/github.com/ragelo/go-reproxy

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o /go/bin/go-reproxy ./cmd

# Verify the binary was created
RUN ls -la /go/bin/go-reproxy


FROM alpine:3.7

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/go-reproxy /usr/local/bin/go-reproxy

RUN chmod +x /usr/local/bin/go-reproxy

# Verify the binary exists and is executable
RUN ls -la /usr/local/bin/go-reproxy

CMD ["/usr/local/bin/go-reproxy"]
