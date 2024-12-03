package iterator

import (
	"gotest.tools/v3/assert"
	"maps"
	"slices"
	"strconv"
	"strings"
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
