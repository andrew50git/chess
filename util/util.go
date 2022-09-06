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

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
