package channels

import (
	"iter"
	"slices"
)

type iterable[T any] interface {
	~func(func(T) bool) | ~[]T | chan T | <-chan T
}

func Send[T any, U iterable[T]](ch chan<- T, in U) int {
	it := iterator[T](in)

	i := 0
	for t := range it {
		ch <- t
		i++
	}
	return i
}

func SendUntilBlocked[T any, U iterable[T]](ch chan<- T, in U) int {
	it := iterator[T](in)

	i := 0
	for t := range it {
		select {
		case ch <- t:
			i++
		default:
			break
		}
	}
	return i
}

func iterator[T any, U iterable[T]](it U) iter.Seq[T] {
	switch ts := any(it).(type) {
	case func(func(T) bool):
		return ts
	case []T:
		return slices.Values(ts)
	case <-chan T:
		return Iterator(ts)
	case chan T:
		return Iterator(ts)
	default:
		panic("unexpected type")
	}
}
