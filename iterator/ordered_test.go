package iterator

import (
	"go-exp/functions/partials"
	"gotest.tools/v3/assert"
	"slices"
	"testing"
)

func Test_MergeOrdered(t *testing.T) {
	odds := TakeWhile(Generate(1, partials.Add(2)), partials.Lt(10))
	evens := TakeWhile(Generate(0, partials.Add(2)), partials.Lt(10))
	threes := TakeWhile(Generate(3, partials.Add(3)), partials.Lt(10))

	s := slices.Collect(MergeOrdered(odds, evens, threes))
	assert.DeepEqual(t, s, []int{0, 1, 2, 3, 3, 4, 5, 6, 6, 7, 8, 9, 9})
}
