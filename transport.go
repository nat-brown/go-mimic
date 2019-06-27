package mimic

import "net/http"

type mapper interface {
	Map(*http.Request) (*http.Response, error)
}

// Transport satisfies the RoundTripper interface.
type Transport struct {
	Mapper mapper
}

// RoundTrip is a passthrough to Map and satisfies the RoundTripper interface for the net/http client.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.Mapper.Map(req)
}
