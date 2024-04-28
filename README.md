# go-reproxy

Forward proxy service to route requests from origin through one of random known third-party proxies.

## Build

```bash

go build -ldflags "-s -w" -o ./out/proxy ./cmd

```

## Run

Create `proxies.txt` in your workdir:
```text
http://<user>:<pass>:<host>:<port>/
socks5://<user>:<pass>:<host>:<port>/
```

Run service:
```bash

./out/proxy --port 8080 --proxiesFile proxies.txt
```
