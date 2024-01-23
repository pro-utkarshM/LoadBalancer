package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type server struct {
	URL           *url.URL
	ReverseProxy  *httputil.ReverseProxy
	Health        bool
	HealthMu      sync.RWMutex
	HealthTimeout time.Duration
}

func newServer(name, addr string) *server {
	u, _ := url.Parse(addr)
	return &server{
		URL:           u,
		ReverseProxy:  httputil.NewSingleHostReverseProxy(u),
		Health:        true,
		HealthTimeout: 5 * time.Second,
	}
}

var (
	serverList = []*server{
		newServer("server-1", "http://127.0.0.1:5001"),
		newServer("server-2", "http://127.0.0.1:5002"),
		newServer("server-3", "http://127.0.0.1:5003"),
		newServer("server-4", "http://127.0.0.1:5004"),
		newServer("server-5", "http://127.0.0.1:5005"),
	}
	lastServedIndex = 0
)

func main() {
	http.HandleFunc("/", forwardRequest)
	go startHealthCheck()
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func forwardRequest(res http.ResponseWriter, req *http.Request) {
	server, err := getHealthyServer()
	if err != nil {
		http.Error(res, "Couldn't process request: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	server.ReverseProxy.ServeHTTP(res, req)
}

func getHealthyServer() (*server, error) {
	for i := 0; i < len(serverList); i++ {
		server := getServer()
		if server.IsHealthy() {
			return server, nil
		}
	}
	return nil, fmt.Errorf("No healthy hosts")
}

func getServer() *server {
	nextIndex := (lastServedIndex + 1) % len(serverList)
	server := serverList[nextIndex]
	lastServedIndex = nextIndex
	return server
}

func startHealthCheck() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, server := range serverList {
				go checkHealth(server)
			}
		}
	}
}

func checkHealth(s *server) {
	// Implement health check logic here
	// For simplicity, let's assume the server is always healthy
	s.HealthMu.Lock()
	s.Health = true
	s.HealthMu.Unlock()
}

// IsHealthy checks if the server is healthy.
func (s *server) IsHealthy() bool {
	s.HealthMu.RLock()
	defer s.HealthMu.RUnlock()
	return s.Health
}
