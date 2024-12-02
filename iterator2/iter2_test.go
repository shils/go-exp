package iterator2

import (
	"go-exp/functions/hof"
	"gotest.tools/v3/assert"
	"maps"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	it := slices.All([]int{1, 2, 3, 4})
	toString := func(i, n int) (string, string) {
		return strconv.Itoa(i), strconv.Itoa(n)
	}

	result := maps.Collect(Map(it, toString))
	expected := map[string]string{"0": "1", "1": "2", "2": "3", "3": "4"}

	assert.DeepEqual(t, result, expected)

}

func TestMap1(t *testing.T) {
	it := slices.All([]int{1, 2, 3, 4})

	result := slices.Collect(Map1(it, func(i, n int) string {
		return strconv.Itoa(n)
	}))
	expected := []string{"1", "2", "3", "4"}

	assert.DeepEqual(t, result, expected)
}

func TestFilter(t *testing.T) {
	it := slices.All([]int{1, 2, 3, 4, 5})
	isEven := func(n int) bool {
		return n%2 == 0
	}

	result := maps.Collect(Filter(it, hof.LiftArity1Left[int, int, bool](isEven)))
	expected := map[int]int{1: 2, 3: 4}

	assert.DeepEqual(t, result, expected)
}

func TestReduce(t *testing.T) {
	it := slices.All([]string{"d", "c", "b", "a"})
	expand := func(acc string, i int, s string) string {
		return acc + strings.Repeat(s, i)
	}

	result := Reduce(it, expand, "")
	assert.Equal(t, result, "cbbaaa")
}
