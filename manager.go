package proxy

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type ProxyManager struct {
	proxies []*url.URL
}

func NewProxyManager(proxies []*url.URL) (*ProxyManager, error) {
	if len(proxies) == 0 {
		return nil, fmt.Errorf("no proxies provided")
	}

	return &ProxyManager{proxies: proxies}, nil
}

func readProxiesFromFile(fileName string) ([]*url.URL, error) {
	if !strings.HasPrefix(fileName, "/") {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		fileName = filepath.Join(cwd, fileName)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	proxies := []*url.URL{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		url, err := url.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
		}
		proxies = append(proxies, url)
	}

	return proxies, nil
}

func NewProxyManagerFromFile(filePath string) (*ProxyManager, error) {
	proxies, err := readProxiesFromFile(filePath)
	if err != nil {
		return nil, err
	}
	return NewProxyManager(proxies)
}

func (pm *ProxyManager) GetProxy() *url.URL {
	// TODO: implement round robin and least connections algorithms
	return pm.proxies[rand.Intn(len(pm.proxies))]
}
