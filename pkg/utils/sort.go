package utils

//SortBy sort b by order of a, for solving n+1
func SortBy[T comparable, R any](a []T, b []R, f func(r R) T) []R {
	ret := make([]R, len(a))
	for i, t := range a {
		for _, r := range b {
			if f(r) == t {
				ret[i] = r
			}
		}
	}
	return ret
}

func SortById[T comparable, R interface{ GetId() T }](a []T, b []R) []R {
	return SortBy[T, R](a, b, func(r R) T { return r.GetId() })
}
