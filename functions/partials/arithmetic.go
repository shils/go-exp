package partials

import (
	"go-exp/constraints"
	goconstraints "golang.org/x/exp/constraints"
)

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

func Mod[T goconstraints.Integer](x T) func(T) T {
	return func(y T) T {
		return y % x
	}
}
