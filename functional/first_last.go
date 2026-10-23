package functional

import "iter"

func (f *Functioner[E]) First(defaultValue E) E {
	var s iter.Seq[E]
	if !f.shouldCompare() && !f.shouldFilter() {
		s = f.sequence
	} else {
		s = f.sequencer(sequencerOptions{
			IsOrdered: true,
			MaxSize:   1,
		})
	}
	for e := range s {
		return e
	}
	return defaultValue
}

func (f *Functioner[E]) FirstBy(predicate func(E) bool, defaultValue E) E {
	if !f.shouldCompare() && !f.shouldFilter() && f.hints.slice != nil {
		for i := range f.hints.slice {
			e := f.hints.slice[i]
			if !predicate(e) {
				continue
			}
			return e
		}
		return defaultValue
	}

	for e := range newF(
		f.sequence,
		copyOnAppend(f.filters, predicate),
		f.comparers,
		f.hints.clone(),
	).sequencer(sequencerOptions{MaxSize: 1}) {
		return e
	}
	return defaultValue
}

func (f *Functioner[E]) Last(defaultValue E) E {
	if !f.shouldCompare() && !f.shouldFilter() && f.hints.slice != nil {
		if len(f.hints.slice) == 0 {
			return defaultValue
		}
		return f.hints.slice[len(f.hints.slice)-1]
	}

	for e := range f.sequencer(sequencerOptions{
		MaxSize:                     1,
		IsOrdered:                   true,
		ShouldIterateInReverseOrder: true,
	}) {
		return e
	}
	return defaultValue
}

func (f *Functioner[E]) LastBy(predicate func(E) bool, defaultValue E) E {
	shouldCompare := f.shouldCompare()
	shouldFilter := f.shouldFilter()

	if f.hints.slice != nil && !shouldCompare {
		p := f.filter()
		//goland:noinspection GoFixSlicesBackward
		for i := len(f.hints.slice) - 1; i >= 0; i-- {
			e := f.hints.slice[i]
			if (shouldFilter && !p(e)) || !predicate(e) {
				continue
			}
			return e
		}
		return defaultValue
	}

	var c comparer[E]
	if shouldCompare {
		c = f.comparer(comparerForward)
	}

	var last E
	init := false
	for e := range f.sequencer(sequencerOptions{IsOrdered: false}) {
		if !predicate(e) {
			continue
		}
		if !init {
			init = true
			last = e
			continue
		}
		if shouldCompare && c(e, last) >= 0 || !shouldCompare {
			last = e
		}
	}
	if !init {
		return defaultValue
	}
	return last
}
