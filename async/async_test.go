package async

import (
	"context"
	"errors"
	"gotest.tools/v3/assert"
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

func zero[T any]() T {
	var t T
	return t
}
