package partials

import (
	"golang.org/x/exp/constraints"
)

type Predicate[T any] func(T) bool

func Gt[T constraints.Ordered](x T) Predicate[T] {
	return func(y T) bool {
		return y > x
	}
}

func Lt[T constraints.Ordered](x T) Predicate[T] {
	return func(y T) bool {
		return y < x
	}
}

func Eq[T comparable](x T) Predicate[T] {
	return func(y T) bool {
		return y == x
	}
}

func Ne[T comparable](x T) Predicate[T] {
	return func(y T) bool {
		return y != x
	}
}
