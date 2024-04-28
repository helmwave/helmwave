package helper

func SlicesMap[T any, S interface{ ~[]T }, R any](s S, f func(T) R) []R {
	res := make([]R, len(s))

	for i, v := range s {
		res[i] = f(v)
	}

	return res
}
