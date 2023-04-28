package lazy

import (
	"go-exp/functions/hof"
	"sync"
)

type Value[T any] interface {
	Get() T
}

type syncValue[T any] struct {
	v    T
	once *sync.Once
	get  func() T
}

func From[T any](fn func() T) Value[T] {
	sv := &syncValue[T]{
		once: &sync.Once{},
		get:  fn,
	}
	return sv
}

func (sv *syncValue[T]) Get() T {
	sv.once.Do(func() {
		sv.v = sv.get()
	})
	return sv.v
}

type simpleValue[T any] struct {
	v T
}

func (sv simpleValue[T]) Get() T {
	return sv.v
}

func Of[T any](v T) Value[T] {
	return simpleValue[T]{v: v}
}

func getFn[T any](v Value[T]) func() T {
	switch x := v.(type) {
	case *syncValue[T]:
		return x.get
	case simpleValue[T]:
		return hof.Constant(x.v)
	default:
		return x.Get
	}
}

func Apply[T, U any](fn func(T) U, arg Value[T]) Value[U] {
	get := getFn(arg)
	return From(func() U {
		return fn(get())
	})
}

func Apply2[T, U, V any](fn func(T, U) V, arg1 Value[T], arg2 Value[U]) Value[V] {
	get1, get2 := getFn(arg1), getFn(arg2)
	return From(func() V {
		return fn(get1(), get2())
	})
}
