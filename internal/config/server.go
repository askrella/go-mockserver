package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func Recompress() bool {
	recompress := os.Getenv("RECOMPRESS")
	if len(recompress) == 0 {
		return true
	}

	return recompress == "true"
}

func CacheEnabled() bool {
	cacheEnabled := os.Getenv("CACHE_ENABLED")
	if cacheEnabled == "" {
		return true
	}

	return cacheEnabled == "true"
}

func GetHost() string {
	host := os.Getenv("MOCK_HOST")
	if host == "" {
		log.Println("No mock host set. Using 0.0.0.0")
		return "0.0.0.0"
	}

	return host
}

func GetServerPort() int {
	rawPort := os.Getenv("MOCK_PORT")
	if rawPort == "" {
		log.Println("No MOCK_PORT set. Using :80.")
		return 80
	}

	port, err := strconv.Atoi(rawPort)
	if err != nil {
		log.Fatal("MOCK_PORT is not an integer: " + rawPort)
	}

	return port
}

func GetTargetURI() string {
	target := os.Getenv("MOCK_TARGET")
	if target == "" {
		log.Fatal("No target set. Please define MOCK_TARGET.")
	}

	if !strings.HasPrefix(target, "http") {
		log.Fatal("Missing 'http' (or https) prefix in MOCK_TARGET.")
	}

	return target
}
