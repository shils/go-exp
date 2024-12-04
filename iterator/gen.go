package iterator

func Generate[T any](init T, gen func(T) T) func(func(T) bool) {
	return func(yield func(T) bool) {
		for current := init; yield(current); current = gen(current) {
		}
	}
}
