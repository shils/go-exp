package operators

import "go-exp/constraints"

func Incr[T constraints.Real](x T) T {
	return x + 1
}

func IncrP[T constraints.Real](x *T) {
	*x++
}

func Decr[T constraints.Real](x T) T {
	return x - 1
}

func DecrP[T constraints.Real](x *T) {
	*x--
}

func Not(b bool) bool {
	return !b
}

func NotP(b *bool) {
	*b = !(*b)
}
