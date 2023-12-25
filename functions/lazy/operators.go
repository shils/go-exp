package lazy

import (
	"go-exp/functions/operators"
)

func And(lhs Value[bool], rhs Value[bool]) Value[bool] {
	return Apply2(operators.And, lhs, rhs)
}

func Or(lhs Value[bool], rhs Value[bool]) Value[bool] {
	return Apply2(operators.Or, lhs, rhs)
}

func Xor(lhs Value[bool], rhs Value[bool]) Value[bool] {
	return Apply2(operators.Xor, lhs, rhs)
}

func Append[T any](xs Value[[]T], x Value[T]) Value[[]T] {
	return Apply2(operators.Append[T], xs, x)
}
