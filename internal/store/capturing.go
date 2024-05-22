package store

import (
	"net/http"
)

var transactions []*Transaction

type Transaction struct {
	Request
	Response
}

type Response struct {
	Body    []byte
	Headers map[string][]string
}

type Request struct {
	Headers map[string][]string
	URL     string
	Method  string
}

func (r Transaction) Equals(o Transaction) bool {
	return r.URL == o.URL && r.Method == o.Method // Ignore headers, usually not relevant for content.
}

func FindTransaction(r *http.Request) (*Transaction, bool) {
	incoming := toTransaction(r)
	for _, transaction := range transactions {
		if transaction.Equals(*incoming) {
			return transaction, true
		}
	}

	return &Transaction{}, false
}

func CaptureRequest(request *http.Request) {
	transactions = append(transactions, toTransaction(request))
}

func toTransaction(request *http.Request) *Transaction {
	return &Transaction{Request: Request{
		Headers: request.Header,
		URL:     request.URL.RequestURI(),
		Method:  request.Method,
	}, Response: Response{}}
}
