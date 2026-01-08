// Package slices provides generic slice manipulation functions.
package slices

// Map applies the given function to each element of the slice,
// returning a new slice with the results. The callback function
// receives both the index and the element.
func Map[T1, T2 any](s []T1, f func(int, T1) T2) []T2 {
	if s == nil {
		return nil
	}

	result := make([]T2, len(s))
	for i, v := range s {
		result[i] = f(i, v)
	}
	return result
}

// Filter returns a new slice containing only the elements
// for which the predicate function returns true. The callback
// function receives both the index and the element.
func Filter[T any](s []T, f func(int, T) bool) []T {
	if s == nil {
		return nil
	}

	var result []T
	for i, v := range s {
		if f(i, v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce applies the reduction function to the elements of the slice,
// accumulating a single result value.
func Reduce[T1, T2 any](s []T1, initial T2, f func(T2, T1) T2) T2 {
	result := initial
	for _, v := range s {
		result = f(result, v)
	}
	return result
}
