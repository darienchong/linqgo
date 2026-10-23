package functional

func (f *Functioner[E]) ElementAt(i int) (E, bool) {
	if !f.shouldFilter() && !f.shouldCompare() && 0 <= i && i < len(f.hints.slice) {
		return f.hints.slice[i], true
	}

	idx := 0
	for e := range f.sequencer(sequencerOptions{MaxSize: i + 1}) {
		if idx == i {
			return e, true
		}
		idx++
	}
	return *new(E), false
}
