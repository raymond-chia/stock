package analyze

import "golang.org/x/exp/constraints"

// https://gosamples.dev/generics-min-max/
func Max[T constraints.Ordered](s ...T) T {
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

func Min[T constraints.Ordered](s ...T) T {
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

func Sum(s []float64) float64 {
	sum := 0.0
	for _, e := range s {
		sum += e
	}
	return sum
}
