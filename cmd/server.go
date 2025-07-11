package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	proxy "github.com/ragelo/go-reproxy"
)

var (
	port            = flag.String("port", "8080", "port to listen on")
	authUser        = flag.String("authUser", "", "username for proxy authentication")
	authPass        = flag.String("authPass", "", "password for proxy authentication")
	proxiesFilePath = flag.String("proxiesFile", "proxies.txt", "path to the file with proxies URLs")
)

func main() {
	flag.Parse()

	manager, err := proxy.NewProxyManagerFromFile(*proxiesFilePath)
	if err != nil {
		log.Fatalf("failed to create proxy manager: %v", err)
	}
	handler := proxy.NewProxyRequestHandler(manager, proxy.ProxyAuthConfig{
		User: *authUser,
		Pass: *authPass,
	})

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: http.HandlerFunc(handler.Handler),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Fatal(server.ListenAndServe())
}
