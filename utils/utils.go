package utils

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Comparable interface {
	rune | string | Number
}

// IsAnyOf returns true if a is equal to any element of v.
func IsAnyOf[T Comparable](a T, v []T) bool {
	for _, x := range v {
		if a == x {
			return true
		}
	}
	return false
}

func Equal[T Number](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
