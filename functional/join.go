package functional

func estimateSeqSize[E any](f *Functioner[E]) (int, bool) {
	if !f.shouldFilter() && f.hints.slice != nil {
		return len(f.hints.slice), true
	}
	if f.hints.lengthEstimate > 0 {
		return f.hints.lengthEstimate, true
	}
	return -1, false
}

func (f *Functioner[E1]) Join[K comparable, E2, V any](
	f2 *Functioner[E2],
	keyer1 func(E1) K,
	keyer2 func(E2) K,
	projector func(E1, E2) V,
) *Functioner[V] {
	// To preserve outer-major-inner-minor ordering, we can't choose the (estimated) smaller sequence to hash and probe,
	// so we always hash the inner sequence (f2) and probe using the outer sequence (f1).
	return newF(func(yield func(V) bool) {
		f2m := make(map[K][]E2, 4) // Arbitrary guess
		processE2 := func(e2 E2) bool {
			k2 := keyer2(e2)
			f2m[k2] = append(f2m[k2], e2)
			return true
		}
		processE1 := func(e1 E1) bool {
			k1 := keyer1(e1)
			e2s, isInJoin := f2m[k1]
			if isInJoin {
				for i := range e2s {
					if !yield(projector(e1, e2s[i])) {
						return false
					}
				}
				return true
			}
			return true
		}
		if f.hints.slice != nil {
			s2 := f2.slice(true, true)
			for i := range s2 {
				if !processE2(s2[i]) {
					return
				}
			}
			s1 := f.slice(true, true)
			for i := range s1 {
				if !processE1(s1[i]) {
					return
				}
			}
		} else {
			f2.sequencer(sequencerOptions{IsOrdered: true})(processE2)
			f.sequencer(sequencerOptions{IsOrdered: true})(processE1)
		}
	}, nil, nil, newHints[V]())
}

func (f *Functioner[E1]) JoinBy[E2, V any](
	f2 *Functioner[E2],
	joiner func(E1, E2) bool,
	projector func(E1, E2) V,
) *Functioner[V] {
	return newF(func(yield func(V) bool) {
		for e1 := range f.sequencer(sequencerOptions{IsOrdered: true}) {
			for e2 := range f2.sequencer(sequencerOptions{IsOrdered: true}) {
				if joiner(e1, e2) {
					if !yield(projector(e1, e2)) {
						return
					}
				}
			}
		}
	}, nil, nil, newHints[V]())
}

func (f *Functioner[E1]) LeftJoin[K comparable, E2, V any](
	f2 *Functioner[E2],
	keyer1 func(E1) K,
	keyer2 func(E2) K,
	projector func(E1, *E2) V,
) *Functioner[V] {
	return newF(func(yield func(V) bool) {
		f2es := make(map[K][]E2)
		if !f2.shouldCompare() && f2.hints.slice != nil {
			shouldFilter := f2.shouldFilter()
			filter2 := f2.filter()
			for _, e2 := range f2.hints.slice {
				if shouldFilter {
					if !filter2(e2) {
						continue
					}
				}
				k2 := keyer2(e2)
				f2es[k2] = append(f2es[k2], e2)
			}
		} else {
			for e2 := range f2.sequencer(sequencerOptions{IsOrdered: true}) {
				k2 := keyer2(e2)
				f2es[k2] = append(f2es[k2], e2)
			}
		}

		processE1 := func(e1 E1) bool {
			k1 := keyer1(e1)
			e2s, matched := f2es[k1]
			if matched {
				for i := range e2s {
					e2 := e2s[i]
					if !yield(projector(e1, &e2)) {
						return false
					}
				}
				return true
			}
			return yield(projector(e1, nil))
		}
		if !f.shouldCompare() && f.hints.slice != nil {
			shouldFilter := f.shouldFilter()
			filter := f.filter()
			for _, e1 := range f.hints.slice {
				if shouldFilter {
					if !filter(e1) {
						continue
					}
				}
				if !processE1(e1) {
					return
				}
			}
		} else {
			f.sequencer(sequencerOptions{IsOrdered: true})(processE1)
		}
	}, nil, nil, newHints[V]())
}

func (f *Functioner[E1]) RightJoin[K comparable, E2, V any](
	f2 *Functioner[E2],
	keyer1 func(E1) K,
	keyer2 func(E2) K,
	projector func(*E1, E2) V,
) *Functioner[V] {
	return f2.LeftJoin(f, keyer2, keyer1, func(e2 E2, e1p *E1) V { return projector(e1p, e2) })
}

func (f *Functioner[E1]) GroupJoin[K comparable, E2, V any](
	f2 *Functioner[E2],
	keyer1 func(E1) K,
	keyer2 func(E2) K,
	projector func(E1, []E2) V,
) *Functioner[V] {
	return newF(func(yield func(V) bool) {
		f2es := make(map[K][]E2)
		if !f2.shouldCompare() && f2.hints.slice != nil {
			shouldFilter := f2.shouldFilter()
			filter2 := f2.filter()
			for _, e2 := range f2.hints.slice {
				if shouldFilter {
					if !filter2(e2) {
						continue
					}
				}
				k2 := keyer2(e2)
				f2es[k2] = append(f2es[k2], e2)
			}
		} else {
			f2.sequencer(sequencerOptions{IsOrdered: true})(func(e2 E2) bool {
				k2 := keyer2(e2)
				f2es[k2] = append(f2es[k2], e2)
				return true
			})
		}

		processElt := func(e1 E1) bool {
			k1 := keyer1(e1)
			return yield(projector(e1, f2es[k1]))
		}
		if !f.shouldCompare() && f.hints.slice != nil {
			shouldFilter := f.shouldFilter()
			filter := f.filter()
			for _, e1 := range f.hints.slice {
				if shouldFilter {
					if !filter(e1) {
						continue
					}
				}
				if !processElt(e1) {
					return
				}
			}
		} else {
			f.sequencer(sequencerOptions{IsOrdered: true})(processElt)
		}
	}, nil, nil, newHints[V]())
}
