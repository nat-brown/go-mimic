package client

import "net/http"

// BaseURL is the base url for all api calls.
const BaseURL = "http://localhost:3000"

// New creates a new base client for modification and use by specialty clients.
var New = func() *http.Client {
	return newInnerClient()
}

func newInnerClient() *http.Client {
	return http.DefaultClient
}
