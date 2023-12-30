package streams

import (
	"go-exp/functions/partials"
	"go-exp/streams/channels"
	"go-exp/streams/collectors"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
	"strconv"
	"testing"
)

func Test_Filter(t *testing.T) {
	tcs := []struct {
		in       []int
		expected []int
	}{
		{[]int{2, 3, 7, 5, 10}, []int{7, 10}},
		{[]int{2, 1, 0, 5}, []int{}},
	}

	for i, tc := range tcs {
		outCh := Of(channels.FromSlice(0, tc.in)).
			Filter(partials.Gt(5)).
			Out()
		_ = assert.Check(t, cmp.DeepEqual(collectors.Slice(outCh), tc.expected), "test case %d", i)
	}
}

func Test_Map(t *testing.T) {
	in := make(chan int)
	go func() {
		in <- 2
		in <- 3
		in <- 7
		in <- 5
		in <- 10
		close(in)
	}()

	result := Map(
		Of(in).Filter(partials.Gt(5)),
		func(i int) string {
			return strconv.Itoa(2 * i)
		},
	)
	assert.DeepEqual(t, collectors.Slice(result.Out()), []string{"14", "20"})
}
