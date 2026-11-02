package utils

import "iter"

func Filter[T any](slice []T, predicate func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, n := range slice {
			if predicate(n) {
				if !yield(n) {
					return
				}
			}
		}
	}
}

func Partition[T any](slice []T, predicate func(T) bool) (truthy []T, falsy []T) {
	truthy = make([]T, 0)
	falsy = make([]T, 0)

	for _, n := range slice {
		if predicate(n) {
			truthy = append(truthy, n)
		} else {
			falsy = append(falsy, n)
		}
	}
	return truthy, falsy
}
