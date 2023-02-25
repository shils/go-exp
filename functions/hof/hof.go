package hof

func ConstantV[T any](v T) func(...interface{}) T {
	return func(...interface{}) T { return v }
}

func Constant[T any](v T) func() T {
	return func() T { return v }
}

func Constant1[T, A any](v T) func(A) T {
	return func(a A) T { return v }
}

func Constant2[T, A, B any](v T) func(A, B) T {
	return func(A, B) T { return v }
}

func LiftArity[T, V any](fn func() V) func(T) V {
	return func(_ T) V {
		return fn()
	}
}

func LiftArity1Left[T, U, V any](fn func(U) V) func(T, U) V {
	return func(_ T, u U) V {
		return fn(u)
	}
}

func LiftArity2Left[S, T, U, V any](fn func(T, U) V) func(S, T, U) V {
	return func(_ S, t T, u U) V {
		return fn(t, u)
	}
}

func LiftArity1Right[T, U, V any](fn func(U) V) func(U, T) V {
	return func(u U, _ T) V {
		return fn(u)
	}
}

func LiftArity2Right[S, T, U, V any](fn func(T, U) V) func(T, U, S) V {
	return func(t T, u U, _ S) V {
		return fn(t, u)
	}
}

func IgnoredIndex[T, V any](fn func(T) V) func(int, T) V {
	return func(_ int, t T) V {
		return fn(t)
	}
}
