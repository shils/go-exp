package iterator

import (
	"iter"
)

func Map[T, V any](it iter.Seq[T], fn func(T) V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for t := range it {
			if !yield(fn(t)) {
				break
			}
		}
	}
}

func Map2[T, V, W any](it iter.Seq[T], fn func(T) (V, W)) iter.Seq2[V, W] {
	return func(yield func(V, W) bool) {
		for t := range it {
			if !yield(fn(t)) {
				break
			}
		}
	}
}

func Filter[T any](it iter.Seq[T], fn func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range it {
			if fn(t) && !yield(t) {
				break
			}
		}
	}
}

func Reduce[T, V any](it iter.Seq[T], fn func(V, T) V, init V) V {
	acc := init
	for t := range it {
		acc = fn(acc, t)
	}

	return acc
}

// Zip returns an iterator that yields pairs of elements from the two input iterators, terminating when either
// iterator is exhausted
func Zip[T, U any](it iter.Seq[T], other iter.Seq[U]) iter.Seq2[T, U] {
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
func ZipWith[T, U, V any](it iter.Seq[T], other iter.Seq[U], fn func(T, U) V) iter.Seq[V] {
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

func Indexed[T any](it iter.Seq[T]) iter.Seq2[int, T] {
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

func Concat[T any](its ...iter.Seq[T]) iter.Seq[T] {
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

func Compact[T comparable](it iter.Seq[T]) iter.Seq[T] {
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

func CompactFunc[T any](it iter.Seq[T], eq func(T, T) bool) iter.Seq[T] {
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

func GroupBy[T, K comparable](it iter.Seq[T], keyFn func(T) K) map[K][]T {
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
