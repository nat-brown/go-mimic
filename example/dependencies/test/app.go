package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gobuffalo/buffalo"

	"github.com/nat-brown/go-mimic/example/dependencies"
)

// Request is a request specifically for testing
type Request struct {
	Request *http.Request
	Test    *testing.T
}

// Do executes a request much like http.Do.
func Do(req Request) *http.Response {
	var h *hook
	options := dependencies.Options()
	options.Logger, h = setupBuffaloLogger(req.Test)
	app := dependencies.Construct(options)

	middleware := app.Muxer()
	w := httptest.NewRecorder()
	resp := &buffalo.Response{
		ResponseWriter: w,
	}
	middleware.ServeHTTP(resp, req.Request)

	h.Print()
	return w.Result()
}
