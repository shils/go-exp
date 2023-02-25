package partials

import "go-exp/constraints"

func Add[T constraints.Real](x T) func(T) T {
	return func(y T) T {
		return x + y
	}
}

func Sub[T constraints.Real](x T) func(T) T {
	return func(y T) T {
		return y - x
	}
}

func Mul[T constraints.Real](x T) func(T) T {
	return func(y T) T {
		return x * y
	}
}

func Div[T constraints.Real](x T) func(T) T {
	return func(y T) T {
		return y / x
	}
}

// degenerate cases: partial unary operators

func Incr[T constraints.Real](x T) T {
	return x + 1
}

func Decr[T constraints.Real](x T) T {
	return x - 1
}
