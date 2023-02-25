package streams

import (
	"context"
	"go-exp/functions/hof"
	"go-exp/functions/partials"
	"go-exp/streams/collectors"
	"gotest.tools/v3/assert"
	"testing"
)

func Test_MergeOrdered(t *testing.T) {
	odds := GenerateWhile(0, 1, hof.LiftArity1Left[int](partials.Add(2)), hof.LiftArity1Left[int](partials.Lt(10)))
	evens := GenerateWhile(0, 0, hof.LiftArity1Left[int](partials.Add(2)), hof.LiftArity1Left[int](partials.Lt(10)))
	threes := GenerateWhile(0, 3, hof.LiftArity1Left[int](partials.Add(3)), hof.LiftArity1Left[int](partials.Lt(10)))
	s := collectors.Slice(MergeOrdered(odds, evens, threes))
	assert.DeepEqual(t, s, []int{0, 1, 2, 3, 3, 4, 5, 6, 6, 7, 8, 9, 9})
}

func Test_Generate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	out := GenerateContext(ctx, 0, 1, func(_ int, i int) int { return 2*i + 1 })
	var s []int
	s = Reduce(
		TakeWhile(0, out, hof.IgnoredIndex(partials.Lt(20))),
		s,
		func(acc []int, val int) []int { return append(acc, val) })
	cancel()

	assert.DeepEqual(t, s, []int{1, 3, 7, 15})
}
