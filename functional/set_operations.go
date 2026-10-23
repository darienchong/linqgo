package functional

func (f *Functioner[E]) Concat(f2 *Functioner[E]) *Functioner[E] {
	hints := newHints[E]()
	if f.hints.lengthEstimate >= 0 && f2.hints.lengthEstimate >= 0 {
		hints.lengthEstimate = f.hints.lengthEstimate + f2.hints.lengthEstimate
	}

	return newF(func(yield func(E) bool) {
		shouldContinue := true
		if f.shouldCompare() {
			s := f.slice(true, true)
			for i := range s {
				if !yield(s[i]) {
					return
				}
			}
		} else {
			f.sequencer(sequencerOptions{IsOrdered: false})(func(e E) bool {
				if !yield(e) {
					shouldContinue = false
					return false
				}
				return true
			})
		}

		if !shouldContinue {
			return
		}

		if f2.shouldCompare() {
			s2 := f2.slice(true, true)
			for i := range s2 {
				if !yield(s2[i]) {
					return
				}
			}
		} else {
			f2.sequencer(sequencerOptions{IsOrdered: false})(yield)
		}

	}, nil, nil, hints)
}

func (f *Functioner[E]) IntersectBy[K comparable](f2 *Functioner[E], keyer func(E) K) *Functioner[E] {
	hints := newHints[E]()
	hints.lengthEstimate = -1

	return newF(func(yield func(E) bool) {
		seen := make(map[K]struct{})
		emitted := make(map[K]struct{})
		for e := range f2.sequencer(sequencerOptions{IsOrdered: true}) {
			seen[keyer(e)] = struct{}{}
		}
		for e := range f.sequencer(sequencerOptions{IsOrdered: true}) {
			k := keyer(e)
			if _, isInIntersection := seen[k]; isInIntersection {
				if _, dedupe := emitted[k]; dedupe {
					continue
				}
				emitted[k] = struct{}{}
				if !yield(e) {
					return
				}
			}
		}
	}, nil, nil, hints)
}

func (f *Functioner[E]) UnionBy[K comparable](f2 *Functioner[E], keyer func(E) K) *Functioner[E] {
	hints := newHints[E]()
	hints.lengthEstimate = -1

	return newF(func(yield func(E) bool) {
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
		for e := range f2.sequencer(sequencerOptions{IsOrdered: true}) {
			k := keyer(e)
			if _, exists := seen[k]; exists {
				continue
			}
			seen[k] = struct{}{}
			if !yield(e) {
				return
			}
		}
	}, nil, nil, hints)
}

func (f *Functioner[E]) ExceptBy[K comparable](f2 *Functioner[E], keyer func(E) K) *Functioner[E] {
	hints := newHints[E]()
	hints.lengthEstimate = -1

	return newF(func(yield func(E) bool) {
		seen := make(map[K]struct{})
		for e := range f2.sequencer(sequencerOptions{IsOrdered: true}) {
			seen[keyer(e)] = struct{}{}
		}
		for e := range f.sequencer(sequencerOptions{IsOrdered: true}) {
			k := keyer(e)
			if _, exists := seen[k]; exists {
				continue
			}
			if !yield(e) {
				return
			}
		}
	}, nil, nil, hints)
}
