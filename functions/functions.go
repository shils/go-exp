package functions

import "unsafe"

func Identity[T any](t T) T {
	return t
}

func IsNilPtr[T any](t *T) bool {
	return uintptr(unsafe.Pointer(t)) == 0
}

func XorNilPtr[A any, B any](a *A, b *B) bool {
	return IsNilPtr(a) != IsNilPtr(b)
}

func Zero[T any]() T {
	var ret T
	return ret
}
