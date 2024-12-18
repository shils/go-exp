package iterator

import (
	"go-exp/functions/partials"
	"gotest.tools/v3/assert"
	"maps"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	it := slices.Values([]int{1, 2, 3, 4})
	toString := func(n int) string {
		return strconv.Itoa(n)
	}

	result := slices.Collect(Map(it, toString))
	expected := []string{"1", "2", "3", "4"}

	assert.DeepEqual(t, result, expected)

}

func TestMap2(t *testing.T) {
	it := slices.Values([]string{"a", "three whole words", "two tokens"})

	results := maps.Collect(Map2(it, func(str string) (int, int) {
		words := strings.Fields(str)
		return len(words), len(str)
	}))

	expected := map[int]int{1: 1, 2: 10, 3: 17}
	assert.DeepEqual(t, results, expected)
}

func TestFilter(t *testing.T) {
	it := slices.Values([]int{1, 2, 3, 4, 5})
	isEven := func(n int) bool {
		return n%2 == 0
	}

	result := slices.Collect(Filter(it, isEven))
	assert.DeepEqual(t, result, []int{2, 4})
}

func TestReduce(t *testing.T) {
	it := slices.Values([]int{1, 2, 3, 4})
	sum := func(acc, n int) int {
		return acc + n
	}

	result := Reduce(it, sum, 0)
	assert.Equal(t, result, 10)
}

func TestZip(t *testing.T) {
	it1 := slices.Values([]int{1, 2, 3})
	it2 := slices.Values([]string{"a", "b", "c"})

	result := maps.Collect(Zip(it1, it2))
	expected := map[int]string{1: "a", 2: "b", 3: "c"}

	assert.DeepEqual(t, result, expected)
}

func TestZipWith(t *testing.T) {
	it1 := slices.Values([]int{1, 2, 3})
	it2 := slices.Values([]string{"c", "b", "a"})

	result := slices.Collect(ZipWith(it1, it2, func(n int, s string) string {
		return strings.Repeat(s, n)
	}))
	expected := []string{"c", "bb", "aaa"}

	assert.DeepEqual(t, result, expected)
}

func TestIndexed(t *testing.T) {
	it := slices.Values([]string{"c", "b", "a"})
	result := maps.Collect(Indexed(it))

	expected := map[int]string{0: "c", 1: "b", 2: "a"}
	assert.DeepEqual(t, result, expected)
}

func TestConcat(t *testing.T) {
	it1 := slices.Values([]int{1, 2})
	it2 := slices.Values([]int{3, 4})

	result := slices.Collect(Concat(it1, it2))

	expected := []int{1, 2, 3, 4}
	assert.DeepEqual(t, result, expected)
}

func TestCompact(t *testing.T) {
	it := slices.Values([]int{2, 2, 1, 3, 3, 3})
	result := slices.Collect(Compact(it))

	expected := []int{2, 1, 3}
	assert.DeepEqual(t, result, expected)
}

func TestCompactFunc(t *testing.T) {
	it := slices.Values([]int{2, 7, 1, 9, 9, 9})
	result := slices.Collect(CompactFunc(it, func(a, b int) bool {
		return a%5 == b%5
	}))

	expected := []int{2, 1, 9}
	assert.DeepEqual(t, result, expected)
}

func TestGroupBy(t *testing.T) {
	it := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	result := GroupBy(it, func(n int) bool {
		return n%2 == 0
	})

	assert.DeepEqual(t, result, map[bool][]int{false: {1, 3, 5, 7, 9}, true: {2, 4, 6, 8, 10}})
}

func TestTakeWhile(t *testing.T) {
	tests := []struct {
		s    []int
		cond func(int) bool
		want []int
	}{
		{[]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n < 4 }, []int{1, 2, 3}},
		{[]int{5, 1, 2, 3, 4}, func(n int) bool { return n < 5 }, nil},
	}

	for _, tc := range tests {
		it := slices.Values(tc.s)
		result := slices.Collect(TakeWhile(it, tc.cond))
		assert.DeepEqual(t, result, tc.want)
	}
}

func TestAccumulate(t *testing.T) {
	it := slices.Values([]int{1, 2, 3, 4})
	sum := func(acc, n int) int {
		return acc + n
	}

	result := slices.Collect(Accumulate(it, sum))
	expected := []int{1, 3, 6, 10}

	assert.DeepEqual(t, result, expected)
}

func TestChunk(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	tests := []struct {
		size int
		want [][]int
	}{
		{3, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}},
		{4, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9}}},
		{11, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}}},
		{1, [][]int{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}}},
	}

	for _, tc := range tests {
		it := Chunk(slices.Values(nums), tc.size)
		assert.DeepEqual(t, slices.Collect(it), tc.want)
	}
}

func TestFirst(t *testing.T) {
	it := slices.Values([]int{4, 3, 2, 1})
	result, ok := First(it)
	assert.Equal(t, result, 4)
	assert.Assert(t, ok)

	it = slices.Values([]int{})
	result, ok = First(it)
	assert.Equal(t, result, 0)
	assert.Assert(t, !ok)
}

func TestFirstOrElse(t *testing.T) {
	it := slices.Values([]int{4, 3, 2, 1})
	result := FirstOrElse(it, 10)
	assert.Equal(t, result, 4)

	it = slices.Values([]int{})
	result = FirstOrElse(it, 10)
	assert.Equal(t, result, 10)
}

func TestFlatten(t *testing.T) {
	it := slices.Values([][]int{{1, 2}, {3, 4}, {5, 6}})
	result := slices.Collect(Flatten(it))
	expected := []int{1, 2, 3, 4, 5, 6}

	assert.DeepEqual(t, result, expected)
}

func TestTee(t *testing.T) {
	it := slices.Values([]int{4, 3, 2, 1})
	its := Tee(it, 2)

	result1 := slices.Collect(its[0])
	result2 := slices.Collect(its[1])

	expected := []int{4, 3, 2, 1}
	assert.DeepEqual(t, result1, expected)
	assert.DeepEqual(t, result2, expected)
}

func BenchmarkTee(b *testing.B) {
	if b.N > 10000000 {
		b.Skipf("N too large: %d", b.N)
	}

	it := TakeWhile(Generate(b.N, func(x int) int { return x - 1 }), partials.Gt(0))
	its := Tee(it, 10)

	results := make([][]int, 10)

	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func() {
			results[i] = slices.Collect(its[i])
			wg.Done()
		}()
	}

	expected := slices.Collect(TakeWhile(Generate(b.N, func(x int) int { return x - 1 }), partials.Gt(0)))

	wg.Wait()
	for _, result := range results {
		assert.DeepEqual(b, result, expected)
	}
}
