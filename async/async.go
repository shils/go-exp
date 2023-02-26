package async

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type Future[T any] interface {
	Done() <-chan struct{}
	Get() (T, error)
}

type future[T any] struct {
	done  chan struct{}
	value T
	err   error
}

func (f *future[T]) Done() <-chan struct{} {
	return f.done
}

func (f *future[T]) Get() (T, error) {
	<-f.done
	return f.value, f.err
}

func Compute[V any](fn func() (V, error)) Future[V] {
	fut := &future[V]{done: make(chan struct{})}
	go func() {
		defer close(fut.done)
		v, err := fn()
		fut.value = v
		fut.err = err
	}()
	return fut
}

func ComputeAll[K comparable, V any](ctx context.Context, fn func(context.Context, K) (V, error), keys ...K) Future[map[K]V] {
	fut := &future[map[K]V]{done: make(chan struct{})}

	numKeys := len(keys)
	results := make([]V, numKeys)

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(numKeys)
	go func() {
		defer close(fut.done)
		if err := g.Wait(); err != nil {
			fut.err = err
			return
		}

		fut.value = make(map[K]V)
		for i, v := range results {
			fut.value[keys[i]] = v
		}
	}()

	for i, k := range keys {
		index, key := i, k
		g.Go(func() error {
			v, err := fn(gCtx, key)
			if err != nil {
				return err
			}
			results[index] = v
			return nil
		})
	}
	return fut
}
