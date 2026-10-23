package functional

func (f *Functioner[E]) MaxBy(c comparer[E]) E {
	var init bool
	var best E

	if f.hints.slice != nil && !f.shouldFilter() {
		for _, e := range f.hints.slice {
			if !init {
				best = e
				init = true
				continue
			}
			if c(best, e) < 0 {
				best = e
			}
		}
		return best
	}

	for e := range f.sequencer(sequencerOptions{IsOrdered: false}) {
		if !init {
			best = e
			init = true
			continue
		}
		if c(best, e) < 0 {
			best = e
		}
	}
	return best
}

func (f *Functioner[E]) MinBy(c comparer[E]) E {
	var init bool
	var best E

	if f.hints.slice != nil && !f.shouldFilter() {
		for _, e := range f.hints.slice {
			if !init {
				best = e
				init = true
				continue
			}
			if c(e, best) < 0 {
				best = e
			}
		}
		return best
	}

	for e := range f.sequencer(sequencerOptions{IsOrdered: false}) {
		if !init {
			best = e
			init = true
			continue
		}
		if c(e, best) < 0 {
			best = e
		}
	}
	return best
}
