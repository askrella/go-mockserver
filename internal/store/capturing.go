package store

import "net/http"

var transactions []Transaction

type Transaction struct {
	Request
}

type Request struct {
	Headers map[string][]string
	URL     string
	Method  string
}

func CaptureRequest(request *http.Request) {
	transactions = append(transactions, Transaction{Request{
		Headers: request.Header,
		URL:     request.URL.RequestURI(),
		Method:  request.Method,
	}})
}
