package util

func RemoveIndex[T any](s []T, index int) []T {
	ret := make([]T, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func RemoveValue[T comparable](s []T, value T) []T {
	for i := range s {
		if s[i] == value {
			return RemoveIndex(s, i)
		}
	}
	return []T{}
}

func Contains[T comparable](s []T, value T) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

func Min[T int | float32](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

func Max[T int | float32](a, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}

func Abs[T int | float32](x T) T {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
