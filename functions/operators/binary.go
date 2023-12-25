package operators

import "golang.org/x/exp/constraints"
import xconstraints "go-exp/constraints"

func Gt[T constraints.Ordered](x T, y T) bool {
	return x > y
}

func Lt[T constraints.Ordered](x T, y T) bool {
	return x < y
}

func Eq[T comparable](x T, y T) bool {
	return x == y
}

func Ne[T comparable](x T, y T) bool {
	return x != y
}

func And(x bool, y bool) bool {
	return x && y
}

func Or(x bool, y bool) bool {
	return x || y
}

func Xor(x bool, y bool) bool {
	return x != y
}

func BwXor[T constraints.Integer](x T, y T) T {
	return x ^ y
}

func Add[T xconstraints.Real](x T, y T) T {
	return x + y
}

func Sub[T xconstraints.Real](x T, y T) T {
	return x - y
}

func Mul[T xconstraints.Real](x T, y T) T {
	return x * y
}

func Div[T xconstraints.Real](x T, y T) T {
	return x / y
}

func Mod[T constraints.Integer](x T, y T) T {
	return x % y
}
