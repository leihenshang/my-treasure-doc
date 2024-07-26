package utils

func SliceToStrMap[T comparable, U any](s []U, f func(U) T) map[T]U {
	if len(s) == 0 {
		return nil
	}
	res := make(map[T]U, len(s))
	for _, v := range s {
		res[f(v)] = v
	}
	return res
}
