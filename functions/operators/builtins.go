package operators

import (
	"golang.org/x/exp/constraints"
	"reflect"
)

func Append[T any](xs []T, x T) []T {
	return append(xs, x)
}

func Len(iv interface{}) int {
	switch x := iv.(type) {
	case string:
		return len(x)
	case []bool:
		return len(x)
	case []complex64:
		return len(x)
	case []complex128:
		return len(x)
	case []error:
		return len(x)
	case []float32:
		return len(x)
	case []float64:
		return len(x)
	case []int:
		return len(x)
	case []int8:
		return len(x)
	case []int16:
		return len(x)
	case []int32:
		return len(x)
	case []int64:
		return len(x)
	case []string:
		return len(x)
	case []uint:
		return len(x)
	case []uint8:
		return len(x)
	case []uint16:
		return len(x)
	case []uint32:
		return len(x)
	case []uint64:
		return len(x)
	case []uintptr:
		return len(x)
	case chan bool:
		return len(x)
	case chan complex64:
		return len(x)
	case chan complex128:
		return len(x)
	case chan error:
		return len(x)
	case chan float32:
		return len(x)
	case chan float64:
		return len(x)
	case chan int:
		return len(x)
	case chan int8:
		return len(x)
	case chan int16:
		return len(x)
	case chan int32:
		return len(x)
	case chan int64:
		return len(x)
	case chan string:
		return len(x)
	case chan uint:
		return len(x)
	case chan uint8:
		return len(x)
	case chan uint16:
		return len(x)
	case chan uint32:
		return len(x)
	case chan uint64:
		return len(x)
	case chan uintptr:
		return len(x)
	default:
		return reflect.ValueOf(x).Len()
	}
}

func Copy[T any](dst, src []T) int {
	return copy(dst, src)
}

func Delete[K comparable, V any](m map[K]V, key K) {
	delete(m, key)
}

func Complex[T constraints.Float](r, i T) complex128 {
	return complex(float64(r), float64(i))
}

func Real[T constraints.Complex](c T) float64 {
	return real(complex128(c))
}

func Imag[T constraints.Complex](c T) float64 {
	return imag(complex128(c))
}
