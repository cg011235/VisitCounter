package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("RemoteAddr: %s, Method: %s, URL: %s, Protocol: %s, Duration: %dms",
			r.RemoteAddr, r.Method, r.URL, r.Proto, duration.Milliseconds())
	})
}

func main() {
	// Define the backend servers
	backends := []string{
		"http://web:8080",
	}

	var currentBackend uint64

	// Create a reverse proxy
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			backend := backends[atomic.AddUint64(&currentBackend, 1)%uint64(len(backends))]
			backendURL, err := url.Parse(backend)
			if err != nil {
				log.Fatal(err)
			}
			req.URL.Scheme = backendURL.Scheme
			req.URL.Host = backendURL.Host
		},
	}

	// Create a new mux and register the proxy handler with logging middleware
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(proxy))

	// Start the reverse proxy server
	log.Println("Reverse proxy server is running on port 80")
	log.Fatal(http.ListenAndServe(":80", mux))
}
