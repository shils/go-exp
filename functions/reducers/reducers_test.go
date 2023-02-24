package reducers

import (
	"gotest.tools/v3/assert"
	"testing"
)

type xSlice[T any] []T

func Test_Append(t *testing.T) {
	s := []int{1, 2, 3}
	assert.DeepEqual(t, Append(s, 4), []int{1, 2, 3, 4})

	xs := xSlice[int]{1, 2, 3}
	assert.DeepEqual(t, Append(xs, 4), xSlice[int]{1, 2, 3, 4})
}
