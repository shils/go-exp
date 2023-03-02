package partials

import (
	"gotest.tools/v3/assert"
	"testing"
)

func Test_Lt(t *testing.T) {
	lt7 := Lt(7)
	assert.Check(t, !lt7(7))
	assert.Check(t, lt7(5))
}

func Test_Gt(t *testing.T) {
	gt7 := Gt(7)
	assert.Check(t, !gt7(7))
	assert.Check(t, gt7(10))
}
