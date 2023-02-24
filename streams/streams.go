package streams

import (
	"container/heap"
	"context"
	"golang.org/x/exp/constraints"
)

type indexedItem[T constraints.Ordered] struct {
	v T
	i int
}

type indexedHeap[T constraints.Ordered] []indexedItem[T]

// sort

func (h indexedHeap[T]) Len() int           { return len(h) }
func (h indexedHeap[T]) Less(i, j int) bool { return h[i].v < h[j].v }
func (h indexedHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// heap

func (h *indexedHeap[T]) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(indexedItem[T]))
}

func (h *indexedHeap[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func MergeOrdered[T constraints.Ordered](streams ...<-chan T) <-chan T {
	out := make(chan T, len(streams))
	go func() {
		defer close(out)
		nOpen := len(streams)

		h := &indexedHeap[T]{}
		heap.Init(h)
		for i := 0; i < len(streams); i++ {
			v, ok := <-streams[i]
			if ok {
				heap.Push(h, indexedItem[T]{v, i})
			} else {
				nOpen--
			}
		}

		for nOpen > 0 {
			item := heap.Pop(h).(indexedItem[T])
			out <- item.v
			next, ok := <-streams[item.i]
			if ok {
				heap.Push(h, indexedItem[T]{next, item.i})
			} else {
				nOpen--
			}
		}
	}()
	return out
}

func GenerateContext[T any](ctx context.Context, buffer int, initial T, fn func(T) T) <-chan T {
	out := make(chan T, buffer)
	go func() {
		current := initial
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			default:
			}
			out <- current
			current = fn(current)
		}
	}()
	return out
}

func Buffered[T any](buffer int, ch <-chan T) <-chan T {
	out := make(chan T, buffer)
	go func() {
		defer close(out)
		for v := range ch {
			out <- v
		}
	}()
	return out
}

func GenerateWhile[T any](buffer int, initial T, gen func(int, T) T, cond func(int, T) bool) <-chan T {
	out := make(chan T, buffer)
	go func() {
		defer close(out)
		count := 0
		current := initial
		for cond(count, current) {
			out <- current
			current = gen(count, current)
			count++
		}
	}()
	return out
}

func TakeWhile[T any](buffer int, ch <-chan T, fn func(T) bool) <-chan T {
	out := make(chan T, buffer)
	go func() {
		defer close(out)
		for v := range ch {
			if !fn(v) {
				return
			}
			out <- v
		}
	}()
	return out
}

func Reduce[A, V any](ch <-chan V, initial A, fn func(A, V) A) A {
	acc := initial
	for v := range ch {
		acc = fn(acc, v)
	}
	return acc
}

func ReduceContext[A, V any](ctx context.Context, ch <-chan V, initial A, fn func(A, V) A) (A, error) {
	acc := initial
	for v := range ch {
		select {
		case <-ctx.Done():
			return acc, ctx.Err()
		default:
		}
		acc = fn(acc, v)
	}
	return acc, nil
}

func Map[T, V any](buffer int, ch <-chan T, fn func(T) V) <-chan V {
	out := make(chan V, buffer)
	go func() {
		defer close(out)
		for x := range ch {
			out <- fn(x)
		}
	}()
	return out
}

func Tee[T any](in <-chan T, outs ...chan<- T) {
	n := len(outs)
	defer func() {
		for i := 0; i < n; i++ {
			close(outs[i])
		}
	}()
	for x := range in {
		for i := 0; i < n; i++ {
			outs[i] <- x
		}
	}
}
