package functional

import "slices"

func (f *Functioner[E]) Skip(i int) *Functioner[E] {
	if i <= 0 {
		return newF(f.sequence, slices.Clone(f.filters), slices.Clone(f.comparers), f.hints)
	}

	hints := f.hints.clone()
	hints.slice = nil
	if !f.shouldCompare() && !f.shouldFilter() && f.hints.slice != nil && 0 <= i && i < len(f.hints.slice) {
		hints.slice = f.hints.slice[i:]
	}
	if hints.lengthEstimate >= 0 {
		hints.lengthEstimate = max(0, hints.lengthEstimate-i)
	}

	// Only materialise the entire sequence if we need ordering
	if f.shouldCompare() {
		return newF(func(yield func(E) bool) {
			counter := i
			s := f.slice(true, true)
			for j := range s {
				counter--
				if counter >= 0 {
					continue
				}
				if !yield(s[j]) {
					return
				}
			}
		}, nil, nil, hints)
	}

	return newF(func(yield func(E) bool) {
		counter := i
		for e := range f.sequencer(sequencerOptions{IsOrdered: false}) {
			counter--
			if counter >= 0 {
				continue
			}
			if !yield(e) {
				return
			}
		}
	}, nil, nil, hints)
}

func (f *Functioner[E]) SkipWhile(predicate func(E) bool) *Functioner[E] {
	hints := f.hints.clone()
	hints.slice = nil
	hints.lengthEstimate = -1

	// Only materialise the entire sequence if ordered
	if f.shouldCompare() {
		return newF(func(yield func(E) bool) {
			skip := true
			s := f.slice(true, true)
			for _, e := range s {
				if skip {
					if predicate(e) {
						continue
					}
					skip = false
				}
				if !yield(e) {
					return
				}
			}
		}, nil, nil, hints)
	}

	return newF(func(yield func(E) bool) {
		skip := true
		f.sequencer(sequencerOptions{IsOrdered: false})(func(e E) bool {
			if skip {
				if predicate(e) {
					return true
				}
				skip = false
			}
			return yield(e)
		})
	}, nil, nil, hints)
}

func (f *Functioner[E]) Take(i int) *Functioner[E] {
	if i <= 0 {
		return newF(func(yield func(E) bool) { return }, nil, nil, newHints[E]())
	}

	hints := f.hints.clone()
	hints.slice = nil
	if !f.shouldCompare() && !f.shouldFilter() && f.hints.slice != nil && 0 <= i {
		hints.slice = f.hints.slice[:min(i, len(f.hints.slice))]
	}
	if f.hints.lengthEstimate >= 0 {
		hints.lengthEstimate = min(f.hints.lengthEstimate, max(i, 0))
	} else {
		hints.lengthEstimate = -1
	}

	if i == 0 {
		return newF(func(yield func(E) bool) { /* empty sequence */ }, nil, nil, hints)
	}

	return newF(f.sequencer(sequencerOptions{MaxSize: i}), nil, nil, hints)
}

func (f *Functioner[E]) TakeWhile(predicate func(E) bool) *Functioner[E] {
	hints := f.hints.clone()
	hints.slice = nil
	hints.lengthEstimate = -1

	return newF(func(yield func(E) bool) {
		take := true
		// Keep this as a sequence to avoid materialising the entire slice if we exit early
		f.sequencer(sequencerOptions{IsOrdered: true})(func(e E) bool {
			if take {
				if predicate(e) {
					return yield(e)
				}
				take = false
			}
			return false
		})
	}, nil, nil, hints)
}
