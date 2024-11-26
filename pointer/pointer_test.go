package pointer

import (
	"gotest.tools/v3/assert"
	"strings"
	"testing"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		a    *int
		b    *int
		want bool
	}{
		{nil, nil, true},
		{nil, To(0), false},
		{To(1), To(1), true},
		{To(0), To(1), false},
	}

	for _, tc := range tests {
		assert.Assert(t, Equal(tc.a, tc.b) == tc.want, "Expected Equal(%v, %v) to be %v", tc.a, tc.b, tc.want)
	}
}

func TestEqualFunc(t *testing.T) {
	equalIgnoreCase := func(a, b string) bool {
		return strings.ToLower(a) == strings.ToLower(b)
	}
	tests := []struct {
		a    *string
		b    *string
		want bool
	}{
		{nil, nil, true},
		{nil, To("abc"), false},
		{To("abc"), To("Abc"), true},
		{To("abc"), To("def"), false},
	}

	for _, tc := range tests {
		assert.Assert(t, EqualFunc(tc.a, tc.b, equalIgnoreCase) == tc.want, "Expected EqualFunc(%v, %v, equalIgnoreCase) to be %v", tc.a, tc.b, tc.want)
	}
}

func TestMap(t *testing.T) {
	double := func(i int) int {
		return i * 2
	}

	tests := []struct {
		a    *int
		want *int
	}{
		{nil, nil},
		{To(1), To(2)},
	}

	for _, tc := range tests {
		assert.DeepEqual(t, Map(tc.a, double), tc.want)
	}
}

func TestFlatMap(t *testing.T) {
	doubleP := func(i int) *int {
		ret := i * 2
		return &ret
	}

	tests := []struct {
		a    *int
		want *int
	}{
		{nil, nil},
		{To(1), To(2)},
	}

	for _, tc := range tests {
		assert.DeepEqual(t, FlatMap(tc.a, doubleP), tc.want)
	}
}

func TestOrElse(t *testing.T) {
	tests := []struct {
		a    *int
		b    *int
		want *int
	}{
		{nil, nil, nil},
		{nil, To(1), To(1)},
		{To(0), To(1), To(0)},
	}

	for _, tc := range tests {
		assert.DeepEqual(t, OrElse(tc.a, tc.b), tc.want)
	}
}

func TestOrElseFunc(t *testing.T) {
	supplierOf := func(a *int) func() *int {
		return func() *int {
			return a
		}
	}

	tests := []struct {
		a    *int
		fn   func() *int
		want *int
	}{
		{nil, supplierOf(nil), nil},
		{nil, supplierOf(To(0)), To(0)},
		{To(0), supplierOf(To(1)), To(0)},
	}

	for _, tc := range tests {
		assert.DeepEqual(t, OrElseFunc(tc.a, tc.fn), tc.want)
	}

	var called bool
	fn := func() *int {
		called = true
		return To(1)
	}
	assert.DeepEqual(t, OrElseFunc(To(0), fn), To(0))
	assert.Assert(t, !called, "Expected OrElseFunc(To(0), fallback) to not call fallback")
}

func TestGetOrElse(t *testing.T) {
	tests := []struct {
		a    *int
		b    int
		want int
	}{
		{nil, 1, 1},
		{To(0), 1, 0},
	}

	for _, tc := range tests {
		assert.Equal(t, GetOrElse(tc.a, tc.b), tc.want)
	}
}

func TestGetOrElseFunc(t *testing.T) {
	supplierOf := func(a int) func() int {
		return func() int {
			return a
		}
	}

	tests := []struct {
		a    *int
		fn   func() int
		want int
	}{
		{nil, supplierOf(1), 1},
		{To(0), supplierOf(1), 0},
	}

	for _, tc := range tests {
		assert.Equal(t, GetOrElseFunc(tc.a, tc.fn), tc.want)
	}
}

func TestFilter(t *testing.T) {
	isOdd := func(i int) bool {
		return i%2 == 1
	}

	tests := []struct {
		a    *int
		want bool
	}{
		{nil, false},
		{To(1), true},
		{To(0), false},
	}

	for _, tc := range tests {
		assert.Equal(t, Filter(tc.a, isOdd) != nil, tc.want)
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		a       *int
		want    int
		wantErr error
	}{
		{nil, 0, ErrNilPointer},
		{To(1), 1, nil},
	}

	for _, tc := range tests {
		v, err := Get(tc.a)
		if tc.wantErr != nil {
			assert.ErrorIs(t, err, tc.wantErr)
		} else {
			assert.NilError(t, err)
			assert.Equal(t, v, tc.want)
		}
	}
}
