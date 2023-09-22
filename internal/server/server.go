package server

import (
	"fmt"
	"github.com/askrella/go-mockserver/internal/config"
	"github.com/askrella/go-mockserver/internal/store"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func InitializeServer() {
	// Parse the target URL
	target, err := url.Parse(config.GetTargetURI())
	if err != nil {
		log.Fatal("Error parsing target URL: " + err.Error())
	}

	// Create a reverse proxy to forward requests to the target host
	reverseProxy := httputil.NewSingleHostReverseProxy(target)

	// Create a custom handler function that uses the reverse proxy
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request:", r.Method, r.URL)

		store.CaptureRequest(r)

		// Forward the request to the target host
		reverseProxy.ServeHTTP(w, r)
	})

	// Start the Go server on port 8080
	port := config.GetHost() + ":" + strconv.Itoa(config.GetServerPort())
	fmt.Println("Server listening on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}
