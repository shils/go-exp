package streams

import "go-exp/functions"

type ChannelStream[T, V any] interface {
	Out() <-chan V
	Filter(func(V) bool) ChannelStream[T, V]
	transform(T) (V, bool)
	in() <-chan T
}

type channelStream[T, V any] struct {
	input <-chan T
	fn    func(T) (V, bool)
}

func identity[T any](v T) (T, bool) {
	return v, true
}

func Of[T any](in <-chan T) ChannelStream[T, T] {
	return &channelStream[T, T]{input: in, fn: identity[T]}
}

func (cs *channelStream[T, V]) Out() <-chan V {
	ch := make(chan V)
	go func() {
		for x := range cs.input {
			r, ok := cs.fn(x)
			if ok {
				ch <- r
			}
		}
	}()
	return ch
}

func (cs *channelStream[T, V]) in() <-chan T {
	return cs.input
}

func (cs *channelStream[T, V]) transform(t T) (V, bool) {
	return cs.fn(t)
}

func (cs *channelStream[T, V]) Filter(filterFn func(V) bool) ChannelStream[T, V] {
	transform := cs.transform
	fn := func(t T) (V, bool) {
		res, ok := transform(t)
		if !ok || filterFn(res) {
			return res, ok
		}
		return res, false
	}
	return &channelStream[T, V]{input: cs.in(), fn: fn}
}

func Map[T, U, V any](cs ChannelStream[T, U], mapFn func(U) V) ChannelStream[T, V] {
	fn := func(t T) (V, bool) {
		res, ok := cs.transform(t)
		if !ok {
			return functions.Zero[V](), false
		}
		return mapFn(res), true
	}
	return &channelStream[T, V]{input: cs.in(), fn: fn}
}
