package proxy

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"

	netProxy "golang.org/x/net/proxy"
)

type ProxyAuthConfig struct {
	User string
	Pass string
}

type proxyRequestHandler struct {
	auth         *ProxyAuthConfig
	proxyManager *ProxyManager
}

func NewProxyRequestHandler(proxyManager *ProxyManager, auth ProxyAuthConfig) *proxyRequestHandler {
	return &proxyRequestHandler{proxyManager: proxyManager, auth: &auth}
}

func (h *proxyRequestHandler) handleProxyAuth(w http.ResponseWriter, r *http.Request) bool {
	if h.auth.User == "" || h.auth.Pass == "" {
		return true
	}

	proxyAuth := r.Header.Get("Proxy-Authorization")
	if proxyAuth == "" {
		http.Error(w, "Proxy Authentication Required", http.StatusProxyAuthRequired)
		return false
	}
	proxyAuth = strings.TrimPrefix(proxyAuth, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(proxyAuth)
	if err != nil {
		http.Error(w, "Proxy Authentication Required", http.StatusProxyAuthRequired)
		return false
	}

	pair := strings.SplitN(string(decoded), ":", 2)
	if len(pair) != 2 {
		http.Error(w, "Proxy Authentication Required", http.StatusProxyAuthRequired)
		return false
	}

	if pair[0] != h.auth.User || pair[1] != h.auth.Pass {
		http.Error(w, "Proxy Authentication Required", http.StatusProxyAuthRequired)
		return false
	}

	return true
}

func (h *proxyRequestHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if !h.handleProxyAuth(w, r) {
		return
	}

	dialer, err := netProxy.FromURL(h.proxyManager.GetProxy(), netProxy.Direct)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	dest_conn, err := dialer.Dial("tcp", r.Host)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
