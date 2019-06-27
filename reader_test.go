package mimic

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func absPathFinder(dirName string, otherPaths ...string) (string, error) {
	fp, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error with with getting working directory for absolute path: %v", err)
	}
	if !strings.Contains(fp, dirName) {
		return "", fmt.Errorf("filepath %s does not contain '%s'; cannot get path to test data", fp, dirName)
	}
	split := strings.SplitAfterN(fp, dirName, 2)
	return filepath.Join(append(split[:1], otherPaths...)...), nil
}

func TestLoad(t *testing.T) {
	var r reader
	err := r.load(filepath.Join(dataPath, "basic_test.json"))
	require.NoError(t, err)

	key, err := request{
		Method: http.MethodGet,
		URL:    "http://url.test",
	}.key()
	assert.Equal(t,
		reader{
			cache: map[[md5.Size]byte]accessor{
				key: &responses{
					list: []response{{StatusCode: 200}},
				},
			},
		}, r)
}

func TestLoadBody(t *testing.T) {
	var r reader
	err := r.load(filepath.Join(dataPath, "body.json"))
	require.NoError(t, err)

	key, err := request{
		Method: http.MethodPost,
		URL:    "http://test.body",
		Body: map[string]interface{}{
			"boolean": true,
			"key":     "value",
			"list":    []int{1, 2, 4},
		},
	}.key()
	assert.Equal(t,
		reader{
			cache: map[[md5.Size]byte]accessor{
				key: &responses{
					list: []response{{
						StatusCode: 200,
						Body: map[string]interface{}{
							"status": "success",
						},
					}},
				},
			},
		}, r)
}

func TestLoadMultipleElements(t *testing.T) {
	var r reader
	err := r.load(filepath.Join(dataPath, "multiple_elements.json"))
	require.NoError(t, err)

	key, err := request{
		Method: http.MethodGet,
		URL:    "http://multi",
	}.key()
	otherKey, err := request{
		Method: http.MethodGet,
		URL:    "http://something.else",
	}.key()
	assert.Equal(t,
		reader{
			cache: map[[md5.Size]byte]accessor{
				key: &responses{
					list: []response{{
						StatusCode: 200,
					}, {
						StatusCode: 500,
					}},
				},
				otherKey: &responses{
					list: []response{{StatusCode: 404}},
				},
			},
		}, r)
}

func TestLayerFiles(t *testing.T) {
	var r reader
	err := r.load(filepath.Join(dataPath, "multiple_elements.json"))
	require.NoError(t, err)
	err = r.load(filepath.Join(dataPath, "layering.json"))
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
	assert.Equal(t,
		reader{
			cache: map[[md5.Size]byte]accessor{
				key: &responses{
					list: []response{{StatusCode: 201}},
				},
				otherKey: &responses{
					list: []response{{StatusCode: 404}},
				},
			},
		}, r)
}

func TestLoadBadPath(t *testing.T) {
	var r reader
	err := r.load(filepath.Join(dataPath, "non_existant_file.json"))
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "error opening mimic cache file"),
		fmt.Sprintf("error was masked: %s", err.Error()))
}

func TestBadFormatting(t *testing.T) {
	tts := []struct {
		errorMessage, filename, name string
	}{
		{
			errorMessage: "error decoding mimic cache file: error unmarshaling mock cache: json: cannot unmarshal string into Go struct field response.status_code of type int",
			filename:     "bad_nested_formatting.json",
			name:         "nested object error correctly reported",
		}, {
			errorMessage: "error decoding mimic cache file: error unmarshaling mock cache: json: cannot unmarshal number into Go struct field request.url of type string",
			filename:     "bad_object_formatting.json",
			name:         "single object error correctly reported",
		}, {
			errorMessage: "error decoding mimic cache file: error unmarshaling mock cache: json: cannot unmarshal string into Go value of type mimic.entry",
			filename:     "bad_list_formatting.json",
			name:         "list error correctly reported",
		}, {
			errorMessage: "error decoding mimic cache file: invalid character 'l' looking for beginning of value",
			filename:     "not_json.txt",
			name:         "completely invalid json correctly reports",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			var r reader
			err := r.load(filepath.Join(dataPath, tt.filename))
			require.Error(t, err)
			assert.Equal(t, tt.errorMessage, err.Error())
		})
	}
}
