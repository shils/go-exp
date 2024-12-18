package iterator

import (
	"iter"
	"slices"
	"sync"
)

func Map[T, V any](it func(func(T) bool), fn func(T) V) func(func(V) bool) {
	return func(yield func(V) bool) {
		for t := range it {
			if !yield(fn(t)) {
				break
			}
		}
	}
}

func Map2[T, V, W any](it func(func(T) bool), fn func(T) (V, W)) func(func(V, W) bool) {
	return func(yield func(V, W) bool) {
		for t := range it {
			if !yield(fn(t)) {
				break
			}
		}
	}
}

func Filter[T any](it func(func(T) bool), fn func(T) bool) func(func(T) bool) {
	return func(yield func(T) bool) {
		for t := range it {
			if fn(t) && !yield(t) {
				break
			}
		}
	}
}

func Reduce[T, V any](it func(func(T) bool), fn func(V, T) V, init V) V {
	acc := init
	for t := range it {
		acc = fn(acc, t)
	}

	return acc
}

// Zip returns an iterator that yields pairs of elements from the two input iterators, terminating when either
// iterator is exhausted
func Zip[T, U any](it func(func(T) bool), other func(func(U) bool)) func(func(T, U) bool) {
	return func(yield func(T, U) bool) {
		next, stop := iter.Pull(other)
		defer stop()
		for t := range it {
			u, ok := next()
			if !ok || !yield(t, u) {
				break
			}
		}
	}
}

// ZipWith returns an iterator that yields the results of applying a function to pairs of elements from the two input
// iterators, terminating when either iterator is exhausted
func ZipWith[T, U, V any](it func(func(T) bool), other func(func(U) bool), fn func(T, U) V) func(func(V) bool) {
	return func(yield func(V) bool) {
		next, stop := iter.Pull(other)
		defer stop()
		for t := range it {
			u, ok := next()
			if !ok || !yield(fn(t, u)) {
				break
			}
		}
	}
}

func Indexed[T any](it func(func(T) bool)) func(func(int, T) bool) {
	return func(yield func(int, T) bool) {
		i := 0
		for t := range it {
			if !yield(i, t) {
				break
			}
			i++
		}
	}
}

func Concat[T any](its ...func(func(T) bool)) func(func(T) bool) {
	return func(yield func(T) bool) {
		for _, it := range its {
			for t := range it {
				if !yield(t) {
					break
				}
			}
		}
	}
}

func Compact[T comparable](it func(func(T) bool)) func(func(T) bool) {
	return func(yield func(T) bool) {
		var last T
		var started bool
		for t := range it {
			if started && t == last {
				continue
			}

			if !yield(t) {
				break
			}
			last = t
			started = true
		}
	}
}

func CompactFunc[T any](it func(func(T) bool), eq func(T, T) bool) func(func(T) bool) {
	return func(yield func(T) bool) {
		var last T
		var started bool
		for t := range it {
			if started && eq(t, last) {
				continue
			}

			if !yield(t) {
				break
			}
			last = t
			started = true
		}
	}
}

func GroupBy[T, K comparable](it func(func(T) bool), keyFn func(T) K) map[K][]T {
	m := make(map[K][]T)
	for t := range it {
		k := keyFn(t)
		m[k] = append(m[k], t)
	}
	return m
}

func TakeWhile[T any](it func(func(T) bool), cond func(T) bool) func(func(T) bool) {
	return func(yield func(T) bool) {
		for t := range it {
			if !cond(t) || !yield(t) {
				return
			}
		}
	}
}

func Accumulate[T any](it func(func(T) bool), fn func(T, T) T) func(func(T) bool) {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(it)
		defer stop()

		v, ok := next()
		if !ok {
			return
		}

		acc := v
		for yield(acc) {
			if v, ok = next(); !ok {
				return
			} else {
				acc = fn(acc, v)
			}
		}
	}
}

func Chunk[T any](it func(func(T) bool), size int) func(func([]T) bool) {
	return func(yield func([]T) bool) {
		batch := make([]T, size)
		i := 0
		for t := range it {
			batch[i] = t
			i++
			if i == size {
				if !yield(batch) {
					return
				}
				i = 0
				batch = make([]T, size)
			}
		}
		if i > 0 {
			yield(batch[:i])
		}
	}
}

func First[T any](it func(func(T) bool)) (T, bool) {
	for t := range it {
		return t, true
	}
	var zero T
	return zero, false
}

func FirstOrElse[T any](it func(func(T) bool), defaultVal T) T {
	if t, ok := First(it); ok {
		return t
	}
	return defaultVal
}

func Flatten[T any, S ~[]T](it func(func(S) bool)) func(func(T) bool) {
	return func(yield func(T) bool) {
		for s := range it {
			for _, t := range s {
				if !yield(t) {
					return
				}
			}
		}
	}
}

type teeState[T any] struct {
	next      func() (T, bool)
	mu        sync.Mutex
	buf       []T
	positions []int
}

func (s *teeState[T]) advanceOne(i int) (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.positions[i] == len(s.buf) {
		t, ok := s.next()
		if !ok {
			return t, ok
		}
		s.buf = append(s.buf, t)
		minPos := slices.Min(s.positions)
		if minPos > 0 {
			s.buf = s.buf[minPos:]
			for j := range s.positions {
				s.positions[j] -= minPos
			}
		}
	}
	pos := s.positions[i]
	s.positions[i]++
	return s.buf[pos], true
}

func Tee[T any](it func(func(T) bool), n int) []func(func(T) bool) {
	next, stop := iter.Pull(it)

	state := &teeState[T]{
		next:      next,
		positions: make([]int, n),
	}
	stopped := 0

	outs := make([]func(func(T) bool), n)
	for i := range n {
		i := i
		outs[i] = func(yield func(T) bool) {
			for {
				t, ok := state.advanceOne(i)
				if !ok {
					return
				}
				if !yield(t) {
					stopped++
					break
				}
			}
			if stopped == n {
				stop()
			}
		}
	}
	return outs
}
