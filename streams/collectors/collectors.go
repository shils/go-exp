package collectors

func Slice[T any](ch <-chan T) []T {
	s := make([]T, 0)
	for x := range ch {
		s = append(s, x)
	}
	return s
}

func ToSlice[T any](ch <-chan T, s []T) []T {
	for x := range ch {
		s = append(s, x)
	}
	return s
}

var void = struct{}{}

func Set[K comparable](ch <-chan K) map[K]struct{} {
	s := make(map[K]struct{})
	for x := range ch {
		s[x] = void
	}
	return s
}

func ToSet[K comparable](ch <-chan K, s map[K]struct{}) map[K]struct{} {
	for x := range ch {
		s[x] = void
	}
	return s
}

func Map[K comparable, V any](ch <-chan K, fn func(K) V) map[K]V {
	m := make(map[K]V)
	for x := range ch {
		m[x] = fn(x)
	}
	return m
}

func ToMap[K comparable, V any](ch <-chan K, fn func(K) V, m map[K]V) map[K]V {
	for x := range ch {
		m[x] = fn(x)
	}
	return m
}
