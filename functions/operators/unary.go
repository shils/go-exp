package operators

import (
	xconstraints "go-exp/constraints"
	"golang.org/x/exp/constraints"
)

func Incr[T xconstraints.Real](x T) T {
	return x + 1
}

func IncrP[T xconstraints.Real](x *T) {
	*x++
}

func Decr[T xconstraints.Real](x T) T {
	return x - 1
}

func DecrP[T xconstraints.Real](x *T) {
	*x--
}

func Not(b bool) bool {
	return !b
}

func NotP(b *bool) {
	*b = !(*b)
}

func BwNot[T constraints.Integer](x T) T {
	return ^x
}
