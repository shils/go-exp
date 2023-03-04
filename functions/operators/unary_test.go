package operators

import (
	"gotest.tools/v3/assert"
	"testing"
)

func Test_IncrP(t *testing.T) {
	var i int
	IncrP(&i)
	assert.Equal(t, i, 1)
}

func Test_DecrP(t *testing.T) {
	var i int
	DecrP(&i)
	assert.Equal(t, i, -1)
}

func Test_NotP(t *testing.T) {
	var b bool
	NotP(&b)
	assert.Equal(t, b, true)
}
