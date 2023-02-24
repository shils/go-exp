package partials

import "golang.org/x/exp/constraints"

func Gt[T constraints.Ordered](x T) func(T) bool {
	return func(y T) bool {
		return y > x
	}
}

func Lt[T constraints.Ordered](x T) func(T) bool {
	return func(y T) bool {
		return y < x
	}
}

func Eq[T comparable](x T) func(T) bool {
	return func(y T) bool {
		return y == x
	}
}

func Ne[T comparable](x T) func(T) bool {
	return func(y T) bool {
		return y != x
	}
}
