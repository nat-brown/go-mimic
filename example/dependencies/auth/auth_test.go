package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nat-brown/go-mimic/example/dependencies/test"
)

const (
	username = "a user"
	password = "a pass"
	token    = "612075736572612070617373d41d8cd98f00b204e9800998ecf8427e" // md5 hash of "a usera pass"
)

func TestAuthIntegration(t *testing.T) {
	pass := t.Run("authentication", func(t *testing.T) {
		body := bytes.NewReader([]byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)))
		req, err := http.NewRequest(http.MethodPost, "/auth/token/", body)
		require.NoError(t, err)
		resp := test.Do(test.Request{
			Request: req,
			Test:    t,
		})
		require.Equal(t, http.StatusOK, resp.StatusCode)

		actual := struct {
			Token string `json:"token"`
		}{}
		decoder := json.NewDecoder(resp.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&actual)
		require.NoError(t, err)
		assert.Equal(t, token, actual.Token)
	})
	t.Run("validation", func(t *testing.T) {
		if !pass {
			t.Skip("setup failed")
		}

		body := bytes.NewReader([]byte(fmt.Sprintf(`{"token":"%s"}`, token)))
		req, err := http.NewRequest(http.MethodPost, "/auth/verify/", body)
		require.NoError(t, err)
		resp := test.Do(test.Request{
			Request: req,
			Test:    t,
		})
		require.Equal(t, http.StatusOK, resp.StatusCode)

		actual := struct {
			Username string `json:"username"`
		}{}
		decoder := json.NewDecoder(resp.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&actual)
		require.NoError(t, err)
		assert.Equal(t, username, actual.Username)
	})
}

func TestAuthenticationNegative(t *testing.T) {
	tts := []struct {
		name       string
		body       []byte
		statusCode int
		response   string
	}{
		{
			name:       "no body",
			body:       nil,
			statusCode: http.StatusBadRequest,
			response:   "Invalid request format",
		}, {
			name:       "empty username",
			body:       []byte(`{"username":"","password":"pass"}`),
			statusCode: http.StatusBadRequest,
			response:   "Username and Password required.",
		}, {
			name:       "empty password",
			body:       []byte(`{"username":"user","password":""}`),
			statusCode: http.StatusBadRequest,
			response:   "Username and Password required.",
		}, {
			name:       "invalid type",
			body:       []byte(`{"username":"user","password":2}`),
			statusCode: http.StatusBadRequest,
			response:   "Invalid request format",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewReader(tt.body)
			req, err := http.NewRequest(http.MethodPost, "/auth/token/", body)
			require.NoError(t, err)
			resp := test.Do(test.Request{
				Request: req,
				Test:    t,
			})
			assert.Equal(t, tt.statusCode, resp.StatusCode)
			msg, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.response, string(msg))
		})
	}
}

func TestVerificationNegative(t *testing.T) {
	tts := []struct {
		name       string
		body       []byte
		statusCode int
		response   string
	}{
		{
			name:       "no body",
			body:       nil,
			statusCode: http.StatusBadRequest,
			response:   "Invalid request format",
		}, {
			name:       "empty token",
			body:       []byte(`{"token":""}`),
			statusCode: http.StatusNotFound,
			response:   "",
		}, {
			name:       "invalid type",
			body:       []byte(`{"token":9}`),
			statusCode: http.StatusBadRequest,
			response:   "Invalid request format",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewReader(tt.body)
			req, err := http.NewRequest(http.MethodPost, "/auth/verify/", body)
			require.NoError(t, err)
			resp := test.Do(test.Request{
				Request: req,
				Test:    t,
			})
			require.Equal(t, tt.statusCode, resp.StatusCode)
			msg, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.response, string(msg))
		})
	}
}
