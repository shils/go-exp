package pointer

import "errors"

var ErrNilPointer = errors.New("cannot dereference nil pointer")

func Map[T any, V any](t *T, fn func(T) V) *V {
	if t == nil {
		return nil
	}
	v := fn(*t)
	return &v
}

func FlatMap[T any, V any](t *T, fn func(T) *V) *V {
	if t == nil {
		return nil
	}
	return fn(*t)
}

func Equal[T comparable](a *T, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func EqualFunc[T any](a *T, b *T, eq func(T, T) bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return eq(*a, *b)
}

func To[T any](t T) *T {
	return &t
}

// OrElse is a port of Scala's Option.orElse method.
func OrElse[T any](a *T, b *T) *T {
	if a != nil {
		return a
	}
	return b
}

func OrElseFunc[T any](a *T, fn func() *T) *T {
	if a != nil {
		return a
	}
	return fn()
}

// GetOrElse is a port of Scala's Option.getOrElse method.
func GetOrElse[T any](a *T, b T) T {
	if a != nil {
		return *a
	}
	return b
}

func GetOrElseFunc[T any](a *T, fn func() T) T {
	if a != nil {
		return *a
	}
	return fn()
}

func Filter[T any](t *T, fn func(T) bool) *T {
	if t == nil || !fn(*t) {
		return nil
	}
	return t
}

func Get[T any](t *T) (T, error) {
	var zero T
	if t == nil {
		return zero, ErrNilPointer
	}
	return *t, nil
}
