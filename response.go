package mimic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type response struct {
	Body       interface{} `json:"body"`
	StatusCode int         `json:"status_code"`
	stream     TrackCloser
}

type responses struct {
	called int
	list   []response
}

// TrackCloser tracks if Closed was called on a ReadCloser.
// It does not track if closing was successful.
type TrackCloser interface {
	io.ReadCloser
	WasClosed() bool
}

type trackBody struct {
	closed bool
	io.ReadCloser
}

func (b trackBody) Close() error {
	b.closed = true
	return b.ReadCloser.Close()
}

func (b trackBody) WasClosed() bool { return b.closed }

func (r *response) httpBody() ([]byte, error) {
	if r == nil || r.Body == nil {
		return nil, nil
	}
	b, err := json.Marshal(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling cache body: %v -- %v", r.Body, err)
	}
	return b, nil
}

func (r *response) HTTPResponse() (*http.Response, error) {
	body, err := r.httpBody()
	if err != nil {
		return nil, fmt.Errorf("error creating http response: %v", err)
	}

	resp := &http.Response{
		StatusCode: r.StatusCode,
	}
	tracker := trackBody{ReadCloser: ioutil.NopCloser(bytes.NewBuffer(body))}
	resp.Body = tracker
	r.stream = tracker

	return resp, nil
}

func (r *response) WasClosed() bool { return r.stream.WasClosed() }

func (rs *responses) Get() (resp *response, ok bool) {
	if rs == nil || rs.called >= len(rs.list) {
		return nil, false
	}
	resp = &rs.list[rs.called]
	rs.called++
	return resp, true
}

// Open returns if any of the returned responses were left open.
func (rs *responses) Open() bool {
	for i := 0; i < rs.called; i++ {
		if !rs.list[i].WasClosed() {
			return true
		}
	}
	return false
}

func (rs *responses) Set(resp response) {
	// Do not check for nil; allow panic to have helpful stack trace.
	rs.list = append(rs.list, resp)
}
