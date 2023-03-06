package operators

import "reflect"

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
