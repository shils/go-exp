package iterator2

import "iter"

func Map[T, U, V, W any](it iter.Seq2[T, U], fn func(T, U) (V, W)) iter.Seq2[V, W] {
	return func(yield func(V, W) bool) {
		for t, u := range it {
			if !yield(fn(t, u)) {
				break
			}
		}
	}
}

func Map1[T, U, V any](it iter.Seq2[T, U], fn func(T, U) V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for t, u := range it {
			if !yield(fn(t, u)) {
				break
			}
		}
	}
}

func Filter[T, U any](it iter.Seq2[T, U], fn func(T, U) bool) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		for t, u := range it {
			if fn(t, u) && !yield(t, u) {
				break
			}
		}
	}
}

func Reduce[T, U, V any](it iter.Seq2[T, U], fn func(V, T, U) V, init V) V {
	acc := init
	for t, u := range it {
		acc = fn(acc, t, u)
	}

	return acc
}
