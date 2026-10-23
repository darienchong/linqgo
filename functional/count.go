package functional

func (f *Functioner[E]) Count() int {
	if !f.shouldFilter() && f.hints.slice != nil {
		return len(f.hints.slice)
	}
	count := 0
	// If we need to filter, then the slice hint is invalidated i.e. we can only use the sequencer
	f.sequencer(sequencerOptions{IsOrdered: false})(func(e E) bool {
		count++
		return true
	})
	return count
}
