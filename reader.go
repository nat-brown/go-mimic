package mimic

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type accessor interface {
	Get() (*response, bool)
	Set(response)
}

type reader struct {
	cache map[[md5.Size]byte]accessor
}

type entry struct {
	Request  request  `json:"request"`
	Response response `json:"response"`
}

// NewReader returnes a new reader instance for quick access.
func newReader() reader {
	return reader{
		cache: map[[md5.Size]byte]accessor{},
	}
}

func (r *reader) getFromRequest(req request) (resp *response, ok bool) {
	key, err := req.key()
	if err != nil {
		return resp, false
	}
	_, ok = r.cache[key]
	if !ok {
		return resp, false
	}
	return r.cache[key].Get()
}

// Initialize defends against actions on nil attributes.
// Does not overwrite pre-existing ones.
func (r *reader) initialize() {
	if r.cache == nil {
		r.cache = map[[md5.Size]byte]accessor{}
	}
}

// Insert adds new data to the reader.
func (r *reader) insert(e entry) error {
	key, err := e.Request.key()
	if err != nil {
		return fmt.Errorf("error with item %v: %v", e, err)
	}

	if _, ok := r.cache[key]; !ok {
		r.cache[key] = &responses{
			list: []response{},
		}
	}
	r.cache[key].Set(e.Response)
	return nil
}

// Load takes an absolute path to the json data file.
func (r *reader) load(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening mimic cache file: %v", err)
	}
	defer jsonFile.Close()

	err = json.NewDecoder(jsonFile).Decode(&r)
	if err != nil {
		return fmt.Errorf("error decoding mimic cache file: %v", err)
	}

	return nil
}

// UnmarshalJSON unmarshals file contents into a mimic cache instance
func (r *reader) UnmarshalJSON(data []byte) error {
	if r == nil {
		return errors.New("reader was nil during unmarshaling")
	}
	r.initialize()

	subR := newReader()

	entries := []entry{}
	err := json.Unmarshal(data, &entries)
	if err == nil {
		for _, e := range entries {
			err := subR.insert(e)
			if err != nil {
				return fmt.Errorf("error with item %+v: %v", e, err)
			}
		}
	} else {
		if newErr := json.Unmarshal(data, &[]interface{}{}); newErr == nil {
			return fmt.Errorf("error unmarshaling mock cache: %v", err) // Data was a list, but badly formatted.
		}

		e := entry{}
		err := json.Unmarshal(data, &e)
		if err != nil {
			return fmt.Errorf("error unmarshaling mock cache: %v", err)
		}

		err = subR.insert(e)
		if err != nil {
			return fmt.Errorf("error with item %+v: %v", e, err)
		}
	}

	for key, val := range subR.cache {
		r.cache[key] = val
	}

	return nil
}
