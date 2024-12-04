package iterator

import (
	"container/heap"
	"golang.org/x/exp/constraints"
	"iter"
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

type pullIterator[T any] struct {
	next func() (T, bool)
	stop func()
}

func MergeOrdered[T constraints.Ordered](its ...func(func(T) bool)) func(func(T) bool) {
	pulls := make([]pullIterator[T], len(its))
	for i, it := range its {
		next, stop := iter.Pull(it)
		pulls[i] = pullIterator[T]{next, stop}
	}
	return mergeOrderedPulls(pulls)
}

func mergeOrderedPulls[T constraints.Ordered](its []pullIterator[T]) func(func(T) bool) {
	return func(yield func(T) bool) {
		defer func() {
			for _, it := range its {
				it.stop()
			}
		}()

		nOpen := len(its)
		h := &indexedHeap[T]{}
		for i := 0; i < len(its); i++ {
			v, ok := its[i].next()
			if ok {
				h.Push(indexedItem[T]{v, i})
			} else {
				nOpen--
			}
		}
		heap.Init(h)

		for nOpen > 0 {
			item := heap.Pop(h).(indexedItem[T])
			if !yield(item.v) {
				return
			}
			next, ok := its[item.i].next()
			if ok {
				heap.Push(h, indexedItem[T]{next, item.i})
			} else {
				nOpen--
			}
		}
	}
}
