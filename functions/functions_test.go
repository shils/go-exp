package functions

import (
	"gotest.tools/v3/assert"
	"testing"
)

type foo struct{}

func Test_IsNilPtr(t *testing.T) {
	var sp *string
	var ip *int
	var fp *foo
	var f foo
	assert.Check(t, IsNilPtr(sp))
	assert.Check(t, IsNilPtr(ip))
	assert.Check(t, IsNilPtr(fp))
	assert.Check(t, !IsNilPtr(&f))
	assert.Check(t, IsNilPtr[bool](nil))
}

func Test_XorNilPtr(t *testing.T) {
	var sp *string
	var f foo

	assert.Check(t, !XorNilPtr(sp, sp))
	assert.Check(t, XorNilPtr(sp, &f))
	assert.Check(t, !XorNilPtr[string, bool](sp, nil))
}
