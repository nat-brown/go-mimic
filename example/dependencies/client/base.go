package client

import "net/http"

// BaseURL is the base url for all api calls.
const BaseURL = "http://localhost:3000"

// New returns a new client generator
var New = func() *http.Client {
	return http.DefaultClient
}
