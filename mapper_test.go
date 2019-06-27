package mimic

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prefixFiles(files ...string) []string {
	ss := make([]string, len(files))
	for i, s := range files {
		ss[i] = filepath.Join(dataPath, s)
	}
	return ss
}

func wrapRequestBody(b []byte) (io.Reader, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(b)
	return body, err
}

func TestBasicMap(t *testing.T) {
	m, err := NewMapper(prefixFiles("basic_test.json")...)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, "http://url.test", nil)
	require.NoError(t, err)

	resp, err := m.Map(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, resp.Body)
}

func TestMapNonJSON(t *testing.T) {
	m, err := NewMapper(prefixFiles("nonjson_body.json")...)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "non.json",
		bytes.NewBuffer([]byte(`--This is a message rather than json--`)))
	require.NoError(t, err)

	resp, err := m.Map(req)
	require.NoError(t, err)

	assert.Equal(t, 200000, resp.StatusCode)
}

func TestMapMiss(t *testing.T) {
	m, err := NewMapper(prefixFiles("basic_test.json")...)
	require.NoError(t, err)

	body, err := wrapRequestBody([]byte(`{"key":"value"}`))
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodGet, "http://url.test", body)
	require.NoError(t, err)

	resp, err := m.Map(req)
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestMultiCallMapMiss(t *testing.T) {
	m, err := NewMapper(prefixFiles("basic_test.json")...)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, "http://url.test", nil)
	require.NoError(t, err)

	resp, err := m.Map(req)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	req, err = http.NewRequest(http.MethodGet, "http://url.test", nil)
	require.NoError(t, err)

	headerKey := "h-key"
	req.Header.Add(headerKey, "one")
	req.Header.Add(headerKey, "two")

	resp, err = m.Map(req)
	require.Error(t, err)
	assert.Equal(t, `no matching call for http://url.test with method: GET and body: null (headers: {"H-Key":"one; two"})`, err.Error())
	assert.Nil(t, resp)
}

func TestNewMapper(t *testing.T) {
	m, err := NewMapper(prefixFiles("multiple_elements.json", "layering.json")...)
	require.NoError(t, err)
	key, err := request{
		Method: http.MethodGet,
		URL:    "http://multi",
	}.key()
	require.NoError(t, err)
	otherKey, err := request{
		Method: http.MethodGet,
		URL:    "http://something.else",
	}.key()
	require.NoError(t, err)
	assert.Equal(t, Mapper{
		cache: reader{
			cache: map[[md5.Size]byte]accessor{
				key: &responses{
					list: []response{{StatusCode: 201}},
				},
				otherKey: &responses{
					list: []response{{StatusCode: 404}},
				},
			},
		},
	}, m)
}

func TestNewMapperNoFiles(t *testing.T) {
	_, err := NewMapper()
	require.Error(t, err)
	assert.Equal(t, noPathsToDataGivenError, err.Error())
}
