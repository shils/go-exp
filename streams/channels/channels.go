package channels

import (
	"container/heap"
	"context"
	"go-exp/functions/reducers"
	"golang.org/x/exp/constraints"
	"sync"
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
		for i := 0; i < len(streams); i++ {
			v, ok := <-streams[i]
			if ok {
				h.Push(indexedItem[T]{v, i})
			} else {
				nOpen--
			}
		}
		heap.Init(h)

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

func DropWhile[T any](buffer int, ch <-chan T, fn func(int, T) bool) <-chan T {
	out := make(chan T, buffer)
	go func() {
		defer close(out)
		i := -1
		for v := range ch {
			i++
			if !fn(i, v) {
				out <- v
				break
			}
		}
		for v := range ch {
			i++
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

func Filter[T any](buffer int, ch <-chan T, fn func(int, T) bool) <-chan T {
	out := make(chan T, buffer)
	go func() {
		defer close(out)
		i := 0
		for x := range ch {
			if fn(i, x) {
				out <- x
			}
			i++
		}
	}()
	return out
}

func DoWithIndex[T, V any](ch <-chan T, fn func(int, T)) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		i := 0
		for x := range ch {
			fn(i, x)
			i++
		}
	}()
	return done
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

func BufferedTee[T any](bufLen int, in <-chan T, outs ...chan<- T) {
	n := len(outs)
	var buf []T
	lags := make([]int, n)

	closed := false
	for !closed {
		// catch up as much as possible without blocking
		for i, out := range outs {
			if lags[i] == 0 {
				continue
			}
			sent := sendUntilBlocked(out, buf[len(buf)-lags[i]:])
			lags[i] -= sent
		}

		maxLag := 0
		for _, lag := range lags {
			maxLag = reducers.Max(maxLag, lag)
		}

		select {
		case x, more := <-in:
			if !more {
				closed = true
				break
			}

			for i, out := range outs {
				// block on out channels that are `bufLen` items behind
				if lags[i] == bufLen {
					out <- buf[0]
					lags[i]--
				}

				// blocked out channels can't receive this new element yet
				if lags[i] != 0 {
					lags[i] += 1
					maxLag = reducers.Max(maxLag, lags[i])
					continue
				}
				select {
				case out <- x:
				default:
					lags[i] = 1
					maxLag = reducers.Max(maxLag, 1)
				}
			}
			// if any out channel is blocked, add the latest in element to our buffer
			if maxLag > 0 {
				buf = append(buf, x)
			}
		default:
		}
		// remove unnecessary elements from the buffer
		buf = buf[len(buf)-maxLag:]
	}

	// wait for all blocked out channels to catch up in separate goroutines
	var wg sync.WaitGroup
	for i, lag := range lags {
		out := outs[i]
		if lag == 0 {
			close(out)
			continue
		}

		wg.Add(1)
		// only use the portion of the buffer needed by this out channel
		s := buf[len(buf)-lag:]
		go func() {
			defer wg.Done()
			for _, x := range s {
				out <- x
			}
			close(out)
		}()
	}
	wg.Wait()
}

func sendUntilBlocked[T any](out chan<- T, vs []T) int {
	for i, v := range vs {
		select {
		case out <- v:
		default:
			return i
		}
	}
	return len(vs)
}

func Tail[T any](ch <-chan T, n int) []T {
	buf := make([]T, n)
	count := 0
	for x := range ch {
		buf[count%n] = x
		count++
	}
	if count <= n {
		return buf[0:count:count]
	}
	return append(buf[count%n:], buf[:count%n]...)
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

func Find[T any](ch <-chan T, fn func(T) bool) (result T, found bool) {
	for x := range ch {
		if fn(x) {
			return x, true
		}
	}
	return
}

func FindLast[T any](ch <-chan T, fn func(T) bool) (result T, found bool) {
	for x := range ch {
		if fn(x) {
			result = x
			found = true
		}
	}
	return
}
