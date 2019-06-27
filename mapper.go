package mimic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const noPathsToDataGivenError = "No paths to data given"

// NewMapper returns a new mapper.
func NewMapper(dataPaths ...string) (Mapper, error) {
	if len(dataPaths) == 0 {
		return Mapper{}, errors.New(noPathsToDataGivenError)
	}
	c := Mapper{cache: newReader()}
	for _, path := range dataPaths {
		err := c.cache.load(path)
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

// Mapper maps http requests to responses.
type Mapper struct {
	cache reader
}

// Map the request to a response
func (m *Mapper) Map(req *http.Request) (*http.Response, error) {
	item, err := m.getFromRequest(req)
	if err != nil {
		return nil, err
	}
	return item.HTTPResponse()
}

func (m *Mapper) getFromRequest(req *http.Request) (resp *response, err error) {
	body, err := prepReqBody(req)
	if err != nil {
		return resp, err
	}
	cacheReq := request{
		URL:    req.URL.String(),
		Method: req.Method,
		Body:   body,
	}
	item, ok := m.cache.getFromRequest(cacheReq)
	if !ok {
		return resp, fmt.Errorf("no matching call for %s with method: %s and body: %s (headers: %v)",
			req.URL.String(), req.Method, formatBodyForError(body), formatHeadersForError(req.Header))
	}
	return item, nil
}

func formatBodyForError(body interface{}) string {
	j, err := json.Marshal(body) // Ignore error as it's only for printing.
	if err == nil {
		return string(j)
	}
	return fmt.Sprintf("%v", body)
}

func formatHeadersForError(originalHeader http.Header) string {
	newHeaderMap := map[string]string{}
	for key, vals := range originalHeader {
		newHeaderMap[http.CanonicalHeaderKey(key)] = strings.Join(vals, "; ")
	}
	j, err := json.Marshal(newHeaderMap) // Ignore error as it's only for printing.
	if err == nil {
		return string(j)
	}
	return fmt.Sprintf("%v", newHeaderMap)
}

func prepReqBody(req *http.Request) (interface{}, error) {
	if req.Body == nil {
		return nil, nil
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body for mimic key: %v", err)
	}
	var convertedBody interface{}
	err = json.Unmarshal(body, &convertedBody)
	if err != nil {
		convertedBody = string(body)
	}
	return convertedBody, nil
}
