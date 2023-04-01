package analyze

import "golang.org/x/exp/constraints"

// https://gosamples.dev/generics-min-max/
func max[T constraints.Ordered](s ...T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m < v {
			m = v
		}
	}
	return m
}

func min[T constraints.Ordered](s ...T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m > v {
			m = v
		}
	}
	return m
}
