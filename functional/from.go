package functional

import (
	"iter"
)

func FromSlice[E any](slice []E) *Functioner[E] {
	hints := newHints[E]()
	hints.slice = slice
	hints.lengthEstimate = len(slice)

	return newF(func(yield func(E) bool) {
		for i := range slice {
			if !yield(slice[i]) {
				return
			}
		}
	}, nil, nil, hints)
}

func FromSeq[E any](seq iter.Seq[E]) *Functioner[E] {
	return newF(seq, nil, nil, newHints[E]())
}

func FromSeq2[K, V any](seq2 iter.Seq2[K, V]) *Functioner[Pair[K, V]] {
	return newF(func(yield func(Pair[K, V]) bool) {
		for k, v := range seq2 {
			if !yield(Pair[K, V]{First: k, Second: v}) {
				return
			}
		}
	}, nil, nil, newHints[Pair[K, V]]())
}

func FromMap[K comparable, V any](m map[K]V) *Functioner[Pair[K, V]] {
	return newF(func(yield func(Pair[K, V]) bool) {
		for k, v := range m {
			if !yield(Pair[K, V]{First: k, Second: v}) {
				return
			}
		}
	}, nil, nil, newHints[Pair[K, V]]())
}
