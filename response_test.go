package mimic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNilResponsesGet(t *testing.T) {
	var rs *responses
	require.True(t, nil == rs, "responses was not nil; test requirement failed")
	_, ok := rs.Get()
	assert.False(t, ok)
}

func TestNilResponsesSet(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "panic was incorrectly silenced")
	}()

	var rs *responses
	require.True(t, nil == rs, "responses was not nil; test requirement failed")
	rs.Set(response{StatusCode: 1}) // This line should panic!
}

func TestResponsesGet(t *testing.T) {
	rs := responses{} // Don't initialize list
	_, ok := rs.Get()
	assert.False(t, ok)
	assert.Equal(t, 0, rs.called) // Uninitialized list wasn't technically called.

	rs.list = []response{
		{StatusCode: 1},
	}
	resp, ok := rs.Get()
	require.True(t, ok)
	assert.Equal(t, &rs.list[0], resp)
	assert.Equal(t, 1, rs.called)

	resp, ok = rs.Get()
	assert.False(t, ok)
	assert.Nil(t, resp)
	assert.Equal(t, 1, rs.called)
}
