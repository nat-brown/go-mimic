package mimic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type response struct {
	Body       interface{} `json:"body"`
	StatusCode int         `json:"status_code"`
}

type responses struct {
	called int
	list   []response
}

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
	if body != nil {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	return resp, nil
}

func (rs *responses) Get() (resp *response, ok bool) {
	if rs == nil || rs.called >= len(rs.list) {
		return nil, false
	}
	resp = &rs.list[rs.called]
	rs.called++
	return resp, true
}

func (rs *responses) Set(resp response) {
	// Do not check for nil; allow panic to have helpful stack trace.
	rs.list = append(rs.list, resp)
}
