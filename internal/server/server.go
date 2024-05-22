package server

import (
	"fmt"
	"github.com/askrella/go-mockserver/internal/config"
	"github.com/askrella/go-mockserver/internal/store"
	"log"
	"net/http"
	"net/http/httputil"
	_ "net/http/pprof"
	"net/url"
	"strconv"
	"strings"
)

func InitializeServer() {
	go func() {
		err := http.ListenAndServe("0.0.0.0:8080", nil)
		if err != nil {
			log.Panicln("Cannot start server for pprof: ", err)
		}
	}()

	// Parse the target URL
	target, err := url.Parse(config.GetTargetURI())
	if err != nil {
		log.Fatal("Error parsing target URL: " + err.Error())
	}

	mux := http.NewServeMux()

	// Create a reverse proxy to forward requests to the target host
	reverseProxy := httputil.NewSingleHostReverseProxy(target)

	// Create a custom handler function that uses the reverse proxy
	reverseProxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the request for our cache
		store.CaptureRequest(r)
		prepareRequest(r)

		// Forward the request to the target host
		reverseProxy.ServeHTTP(w, r)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reverseProxyHandler.ServeHTTP(w, r)
	})

	handler := CacheCompressionHandler(mux)

	// Start the Go server on port 8080
	port := config.GetHost() + ":" + strconv.Itoa(config.GetServerPort())
	fmt.Println("Server listening on port", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}

func prepareRequest(req *http.Request) {
	req.Header.Del("Postman-Token")

	// Let's pretend we are chrome
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	)

	host, _ := strings.CutPrefix(config.GetTargetURI(), "https://")
	host, _ = strings.CutPrefix(host, "http://")

	req.Header.Set("Host", host)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	req.Host = host
}
