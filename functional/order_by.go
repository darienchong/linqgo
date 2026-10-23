package functional

import "cmp"

func (f *Functioner[E]) OrderBy(comparer comparer[E]) *Functioner[E] {
	return newF(f.sequence, f.filters, copyOnAppend(f.comparers, comparer), f.hints.clone())
}

func (f *Functioner[E]) OrderByKey[K cmp.Ordered](keyer func(E) K) *Functioner[E] {
	return f.OrderBy(keyComparer(keyer, false))
}

func (f *Functioner[E]) OrderByDescending(cmp func(E, E) int) *Functioner[E] {
	return f.OrderBy(reverse(cmp))
}

func (f *Functioner[E]) OrderByKeyDescending[K cmp.Ordered](keyer func(E) K) *Functioner[E] {
	return f.OrderBy(keyComparer(keyer, true))
}

func (f *Functioner[E]) ThenBy(cmp comparer[E]) *Functioner[E] {
	if len(f.comparers) == 0 {
		return f.OrderBy(cmp)
	}
	return newF(f.sequence, f.filters, copyOnAppend(f.comparers[:len(f.comparers)-1], then(f.comparers[len(f.comparers)-1], cmp)), f.hints.clone())
}

func (f *Functioner[E]) ThenByDescending(cmp comparer[E]) *Functioner[E] {
	return f.ThenBy(reverse(cmp))
}

func (f *Functioner[E]) ThenByKey[K cmp.Ordered](keyer func(E) K) *Functioner[E] {
	return f.ThenBy(keyComparer(keyer, false))
}

func (f *Functioner[E]) ThenByKeyDescending[K cmp.Ordered](keyer func(E) K) *Functioner[E] {
	return f.ThenBy(keyComparer(keyer, true))
}

func (f *Functioner[E]) Unordered() *Functioner[E] {
	return newF(f.sequence, f.filters, nil, f.hints.clone())
}
