package streams

import (
	"container/heap"
	"container/ring"
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

func GenerateContext[T any](ctx context.Context, buffer int, initial T, fn func(int, T) T) <-chan T {
	out := make(chan T, buffer)
	go func() {
		current := initial
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				close(out)
				return
			default:
			}
			out <- current
			current = fn(i, current)
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
		current := initial
		for i := 0; cond(i, current); i++ {
			out <- current
			current = gen(i, current)
		}
	}()
	return out
}

func TakeWhile[T any](buffer int, ch <-chan T, fn func(int, T) bool) <-chan T {
	out := make(chan T, buffer)
	go func() {
		defer close(out)
		i := 0
		for v := range ch {
			if !fn(i, v) {
				return
			}
			out <- v
			i++
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

func Map[T, V any](buffer int, ch <-chan T, fn func(int, T) V) <-chan V {
	out := make(chan V, buffer)
	go func() {
		defer close(out)
		i := 0
		for x := range ch {
			out <- fn(i, x)
			i++
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

func Tail[T any](ch <-chan T, n int) []T {
	r := ring.New(n)
	h := r
	count := 0
	for x := range ch {
		count++
		r.Value = x
		r = r.Next()
	}
	if count < n {
		return ringToSlice[T](h, count)
	}
	return ringToSlice[T](r, n)
}

func ringToSlice[T any](r *ring.Ring, sLen int) []T {
	out := make([]T, sLen)
	for i := 0; i < sLen; i++ {
		out[i] = r.Value.(T)
		r = r.Next()
	}
	return out
}

func FromSlice[T any](buffer int, s []T) <-chan T {
	out := make(chan T, buffer)
	go func() {
		defer close(out)
		for _, x := range s {
			out <- x
		}
	}()
	return out
}
