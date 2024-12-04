package iterator

import (
	"gotest.tools/v3/assert"
	"slices"
	"testing"
)

func TestGenerate(t *testing.T) {
	it := Generate(0, func(x int) int { return x + 1 })

	result := slices.Collect(TakeWhile(it, func(x int) bool { return x < 5 }))
	expected := []int{0, 1, 2, 3, 4}

	assert.DeepEqual(t, result, expected)
}
