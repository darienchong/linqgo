package functional

import (
	"cmp"
)

type filter[E any] = func(E) bool

type comparer[E any] = func(E, E) int

// copyOnAppend creates a new slice to ensure no slice clobbering during append() calls.
func copyOnAppend[E any](s []E, elts ...E) []E {
	return append(append(make([]E, 0, len(s)+len(elts)+1), s...), elts...)
}

// reverse reverses the order of a comparer.
// Where possible, prioritise folding in the reverse action into a comparer rather than composing reverse() since
// this wraps another function around the comparer.
func reverse[E any](c comparer[E]) comparer[E] {
	return func(a, b E) int {
		return -c(a, b)
	}
}

// keyComparer generates a comparer for a given type E given the key function E -> K.
func keyComparer[K cmp.Ordered, E any](keyer func(E) K, reverse bool) comparer[E] {
	return func(a, b E) int {
		ak, bk := keyer(a), keyer(b)
		if reverse {
			return -cmp.Compare(ak, bk)
		}
		return cmp.Compare(ak, bk)
	}
}

// then concatenates two comparers into a comparer in lexicographical order (elements with the same order are then further compared using the secondary comparer)
func then[E any](cmp1, cmp2 comparer[E]) comparer[E] {
	if cmp1 == nil {
		return cmp2
	}
	return func(a, b E) int {
		if v := cmp1(a, b); v != 0 {
			return v
		}
		return cmp2(a, b)
	}
}

func (f *Functioner[E]) shouldCompare() bool {
	return len(f.comparers) > 0
}

func (f *Functioner[E]) shouldFilter() bool {
	return len(f.filters) > 0
}
