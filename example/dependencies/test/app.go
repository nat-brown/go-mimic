package test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gobuffalo/buffalo"

	"github.com/nat-brown/go-mimic/example/dependencies"
)

// Do executes a request much like http.Do.
func Do(req *http.Request) *http.Response {
	app := dependencies.New()
	middleware := app.Muxer()
	w := httptest.NewRecorder()
	resp := &buffalo.Response{
		ResponseWriter: w,
	}
	middleware.ServeHTTP(resp, req)

	return w.Result()
}
