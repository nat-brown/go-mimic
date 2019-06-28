package set_test

import (
	"testing"

	"github.com/nat-brown/go-mimic/set"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := set.New()
	assert.Equal(t, set.String{}, s)

	// No false positives.
	const key = "key"
	assert.False(t, s.Contains(key))

	// No false negatives.
	s.Add(key)
	assert.True(t, s.Contains(key))

	// Multiple keys don't interfere.
	const other = "other key"
	s.Add(other)
	assert.True(t, s.Contains(key))
	assert.True(t, s.Contains(other))
	assert.False(t, s.Contains("third"))
}

func TestUnion(t *testing.T) {
	const (
		key   = "key"
		other = "other"
		same  = "same"
	)

	setOne := set.New()
	setTwo := set.New()

	setOne.Add(key)
	setOne.Add(same)

	setTwo.Add(other)
	setTwo.Add(same)

	setThree := set.Union(setOne, setTwo)

	assert.True(t, setOne.Contains(key))
	assert.True(t, setOne.Contains(same))
	assert.False(t, setOne.Contains(other))

	assert.False(t, setTwo.Contains(key))
	assert.True(t, setTwo.Contains(other))
	assert.True(t, setTwo.Contains(same))

	assert.True(t, setThree.Contains(key))
	assert.True(t, setThree.Contains(other))
	assert.True(t, setThree.Contains(same))
}
