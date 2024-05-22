package server

import (
	"bytes"
	"github.com/askrella/go-mockserver/internal/compression"
	"github.com/askrella/go-mockserver/internal/config"
	"github.com/askrella/go-mockserver/internal/store"
	"io"
	"log"
	"net/http"
	"strconv"
)

// CacheCompressionHandler provides the ability to
// 1. Store and replay requests and
// 2. DecompressGzip upstream and compress with stronger algorithms/levels for downstream
func CacheCompressionHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		if config.CacheEnabled() {
			if cacheUsed := useTransactionResponseCache(respWriter, request); cacheUsed {
				return
			}
		}

		var resultingBody []byte

		// Intercept the response request with a custom request and get the body data.
		var buffer bytes.Buffer
		decompressedResponseWriter := &responseWriterInterceptor{ResponseWriter: respWriter, Writer: &buffer}
		handler.ServeHTTP(decompressedResponseWriter, request)
		body := buffer.Bytes()
		resultingBody = body

		// Check if the content encoding is gzip and in case it is not, skip
		currentEncoding := decompressedResponseWriter.Header().Get("Content-Encoding")
		if currentEncoding != "gzip" {
			log.Println("Skip compression, content-encoding: " + currentEncoding)
			_, err := respWriter.Write(body)
			if err != nil {
				respWriter.WriteHeader(http.StatusInternalServerError)
				log.Println("Error writing original response without compression: ", err)
				return
			}

			return
		}

		if config.Recompress() {
			// DecompressGzip and re-compress. This seems weird at first, but we want to control how the content
			// gets compressed and sent to the downstream. For example: when loading a huge CSV dataset, we saw substantial
			// performance improvements by using level 9 gzip instead of the original gzip compression.
			decompressed := compression.DecompressGzip(body)
			compressed, contentEncoding := compression.CompressGzip(decompressed)
			contentLength := strconv.FormatInt(int64(len(compressed)), 10)

			respWriter.Header().Set("Content-Encoding", contentEncoding)
			respWriter.Header().Set("Content-Length", contentLength)
			_, err := respWriter.Write(compressed)
			if err != nil {
				respWriter.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to write response", err)
				return
			}

			resultingBody = compressed
		}

		if config.CacheEnabled() {
			storeTransactionResponse(respWriter, request, resultingBody)
		}
	})
}

func useTransactionResponseCache(respWriter http.ResponseWriter, request *http.Request) (cacheUsed bool) {
	transaction, found := store.FindTransaction(request)
	// Cache handling: If a transaction is found, use the body from the cached transaction.
	if !found {
		return false
	}
	log.Println("Found transaction, using cache.")
	response := transaction.Response

	// Write headers from cached transaction to downstream
	for key, value := range response.Headers {
		for _, headerValue := range value {
			respWriter.Header().Add(key, headerValue)
		}
	}

	_, err := respWriter.Write(response.Body)
	if err != nil {
		log.Println(err)
		respWriter.WriteHeader(http.StatusInternalServerError)
		return true
	}

	return true
}

func storeTransactionResponse(respWriter http.ResponseWriter, request *http.Request, compressed []byte) {
	// Find transaction and store our compressed content
	transaction, found := store.FindTransaction(request)
	if !found {
		respWriter.WriteHeader(http.StatusInternalServerError)
		log.Panicln("Transaction not found")
		return
	}

	transaction.Response = store.Response{
		Body:    compressed,
		Headers: respWriter.Header(),
	}
}

type responseWriterInterceptor struct {
	http.ResponseWriter
	Body   io.Reader
	Writer io.Writer
}

func (w *responseWriterInterceptor) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
