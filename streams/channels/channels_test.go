package channels

import (
	"context"
	"go-exp/functions/hof"
	"go-exp/functions/operators"
	"go-exp/functions/partials"
	"go-exp/streams/collectors"
	"gotest.tools/v3/assert"
	"slices"
	"strings"
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

func Test_Tail(t *testing.T) {
	ch := GenerateWhile(10, 0, hof.IgnoredIndex(operators.Incr[int]), hof.IgnoredIndex(partials.Lt(10)))
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)

	go Tee(ch, ch1, ch2)

	assert.DeepEqual(t, Tail(ch1, 4), []int{6, 7, 8, 9})
	// number of channel items < n
	assert.DeepEqual(t, Tail(ch2, 20), []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func Test_BufferedTee(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	in := FromSlice(0, s)
	out1 := make(chan int)
	out2 := make(chan int)
	out3 := make(chan int)

	go BufferedTee(5, in, out1, out2, out3)

	s1 := collectors.Slice(out1)
	s2 := collectors.Slice(out2)
	s3 := collectors.Slice(out3)

	assertAllSlicesEqual(t, s, s1, s2, s3)
}

func Test_Map(t *testing.T) {
	in := make(chan int)
	out := Map(0, in, func(i int, k int) int {
		return i * k
	})

	go func() {
		in <- 1
		in <- 2
		in <- 3
		in <- 4
		close(in)
	}()

	assert.DeepEqual(t, collectors.Slice(out), []int{0, 2, 6, 12})
}

func Test_DropWhile(t *testing.T) {
	in := make(chan string)
	out := DropWhile(0, in, hof.IgnoredIndex(func(s string) bool {
		return strings.HasPrefix(s, "a")
	}))

	go func() {
		in <- "abc"
		in <- "a"
		in <- "b"
		in <- "a"
		defer close(in)
	}()

	assert.DeepEqual(t, collectors.Slice(out), []string{"b", "a"})
}

func Test_Find(t *testing.T) {
	in := make(chan int)
	go func() {
		in <- 2
		in <- 3
		in <- 7
		close(in)
	}()

	result, found := Find(in, partials.Gt(5))

	assert.Equal(t, result, 7)
	assert.Check(t, found)
}

func Test_FindFail(t *testing.T) {
	in := make(chan int)
	go func() {
		in <- 2
		in <- 3
		in <- 1
		close(in)
	}()

	result, found := Find(in, partials.Gt(5))

	assert.Equal(t, result, 0)
	assert.Check(t, !found)
}

func Test_FindLast(t *testing.T) {
	in := make(chan int)

	go func() {
		in <- 2
		in <- 8
		in <- 7
		in <- 1
		close(in)
	}()

	result, found := FindLast(in, partials.Gt(5))

	assert.Equal(t, result, 7)
	assert.Check(t, found)
}

func Test_FindLastFail(t *testing.T) {
	in := make(chan int)

	go func() {
		in <- 2
		in <- 1
		in <- 3
		close(in)
	}()

	result, found := FindLast(in, partials.Gt(5))

	assert.Equal(t, result, 0)
	assert.Check(t, !found)
}

func Test_Iterator(t *testing.T) {
	in := make(chan int, 3)
	go func() {
		in <- 3
		in <- 1
		in <- 2
		close(in)
	}()

	it := Iterator(in)
	assert.DeepEqual(t, slices.Collect(it), []int{3, 1, 2})
}

func assertAllSlicesEqual[T any](t *testing.T, expected []T, tss ...[]T) {
	for _, ts := range tss {
		assert.DeepEqual(t, expected, ts)
	}
}
