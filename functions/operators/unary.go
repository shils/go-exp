package operators

import "go-exp/constraints"

func Incr[T constraints.Real](x T) T {
	return x + 1
}

func Decr[T constraints.Real](x T) T {
	return x - 1
}

func Not(b bool) bool {
	return !b
}
