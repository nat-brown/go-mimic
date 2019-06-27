package mimic

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
)

// Request object for hash key.
type request struct {
	Body   interface{} `json:"body"`
	Method string      `json:"method"`
	URL    string      `json:"url"`
}

func (r request) key() ([md5.Size]byte, error) {
	bytes := []byte(r.URL + r.Method)
	if r.Body != nil {
		body, err := json.Marshal(r.Body)
		if err != nil {
			return [md5.Size]byte{}, fmt.Errorf("error marshalling body for key: %v", err)
		}
		bytes = append(bytes, body...)
	}
	return md5.Sum(bytes), nil
}
