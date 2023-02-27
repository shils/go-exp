package async

import (
	"context"
	"errors"
	"gotest.tools/v3/assert"
	"strings"
	"testing"
	"time"
)

func Test_ComputeError(t *testing.T) {
	vCh := make(chan int)
	eCh := make(chan error)
	fut := Compute(func() (int, error) {
		select {
		case v := <-vCh:
			return v, nil
		case e := <-eCh:
			return 0, e
		}
	})

	e := errors.New("closed")
	eCh <- e
	res, err := fut.Get()
	assert.Equal(t, res, 0)
	assert.ErrorIs(t, err, e)
}

func Test_ComputeAll(t *testing.T) {
	fut := ComputeAll(context.Background(), func(ctx context.Context, s string) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		time.Sleep(10 * time.Millisecond)
		return len(s), nil
	}, "a", "ab", "abc", "abcd")
	expected := map[string]int{
		"a":    1,
		"ab":   2,
		"abc":  3,
		"abcd": 4,
	}

	result, err := fut.Get()
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, result)
}

func Test_ComputeAllDone(t *testing.T) {
	wait := make(chan struct{})
	fut := ComputeAll(context.Background(), func(ctx context.Context, s string) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-wait:
			return len(s), nil
		}
	}, "a", "ab", "abc", "abcd")

	expected := map[string]int{
		"a":    1,
		"ab":   2,
		"abc":  3,
		"abcd": 4,
	}

	closed := false
	for {
		select {
		case <-fut.Done():
			result, err := fut.Get()
			assert.NilError(t, err)
			assert.DeepEqual(t, expected, result)
			return
		default:
			if !closed {
				close(wait)
				closed = true
			}
		}
	}
}

type wordCounts map[string]int

func (wc wordCounts) Reduce(other wordCounts) wordCounts {
	for k, v1 := range other {
		v0 := wc[k]
		wc[k] = v0 + v1
	}
	return wc
}

func countWords(ctx context.Context, words string) (wordCounts, error) {
	m := make(map[string]int)
	for _, s := range strings.Fields(words) {
		c := m[s]
		m[s] = c + 1
	}
	return m, nil
}

func Test_MapReduce(t *testing.T) {
	keys := []string{"cat dog", "pig cat", "dog dog", "pig mouse pig"}
	fut := MapReduce(context.Background(), countWords, make(map[string]int), keys...)
	expected := wordCounts{
		"cat":   2,
		"dog":   3,
		"pig":   3,
		"mouse": 1,
	}

	res, err := fut.Get()
	assert.DeepEqual(t, expected, res)
	assert.NilError(t, err)
}

func Test_MapReduceError(t *testing.T) {
	tooManyWordsErr := errors.New("too many words")

	keys := []string{"cat dog", "pig cat", "dog dog", "pig mouse pig"}
	badCountWords := func(ctx context.Context, key string) (wordCounts, error) {
		if len(strings.Fields(key)) > 2 {
			return nil, tooManyWordsErr
		}
		return countWords(ctx, key)
	}

	fut := MapReduce(context.Background(), badCountWords, make(map[string]int), keys...)

	res, err := fut.Get()
	assert.Assert(t, res == nil)
	assert.ErrorIs(t, err, tooManyWordsErr)
}

func zero[T any]() T {
	var t T
	return t
}
