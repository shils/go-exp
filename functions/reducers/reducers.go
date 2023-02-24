package reducers

import (
	xConstraints "go-exp/constraints"
	"golang.org/x/exp/constraints"
)

func Append[A ~[]T, T any](acc A, v T) A {
	return append(acc, v)
}

func Add[T xConstraints.Real](acc T, v T) T {
	return acc + v
}

func Max[T constraints.Ordered](acc T, v T) T {
	if v > acc {
		return v
	}
	return acc
}

func Min[T constraints.Ordered](acc T, v T) T {
	if v < acc {
		return v
	}
	return acc
}
