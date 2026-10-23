package functional

func (f *Functioner[E]) DistinctBy[K comparable](keyer func(e E) K) *Functioner[E] {
	return newF[E](
		func(yield func(E) bool) {
			seen := make(map[K]struct{})
			for e := range f.sequencer(sequencerOptions{IsOrdered: true}) {
				k := keyer(e)
				if _, exists := seen[k]; exists {
					continue
				}
				seen[k] = struct{}{}
				if !yield(e) {
					return
				}
			}
		}, nil, nil, newHints[E]())
}

func (f *Functioner[E]) GroupBy[K comparable, V ~[]E](keyer func(e E) K) *Functioner[Pair[K, V]] {
	return newF[Pair[K, V]](func(yield func(Pair[K, V]) bool) {
		const sizeGuess = 4 // Arbitrary guess for size to control reallocation
		groups := make([]V, 0, sizeGuess)
		m := make(map[K]int, sizeGuess)
		keys := make([]K, 0, sizeGuess)
		if f.shouldCompare() {
			s := f.slice(true, true)
			for i := range s {
				k := keyer(s[i])
				groupIndex, exists := m[k]
				if !exists {
					m[k] = len(groups)
					keys = append(keys, k)
					groupIndex = len(groups)
					groups = append(groups, make(V, 0))
				}
				groups[groupIndex] = append(groups[groupIndex], s[i])
			}
		} else {
			s := f.sequencer(sequencerOptions{IsOrdered: true})
			for e := range s {
				k := keyer(e)
				groupIndex, exists := m[k]
				if !exists {
					m[k] = len(groups)
					keys = append(keys, k)
					groupIndex = len(groups)
					groups = append(groups, make(V, 0))
				}
				groups[groupIndex] = append(groups[groupIndex], e)
			}
		}

		for i := range keys {
			k := keys[i]
			if !yield(Pair[K, V]{
				First:  k,
				Second: groups[m[k]],
			}) {
				return
			}
		}
	}, nil, nil, newHints[Pair[K, V]]())
}
