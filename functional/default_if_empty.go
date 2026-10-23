package functional

func (f *Functioner[E]) DefaultIfEmpty(defaultValue E) *Functioner[E] {
	hints := f.hints.clone()
	hints.slice = nil

	return newF(func(yield func(E) bool) {
		yielded := false
		f.sequencer(sequencerOptions{IsOrdered: true})(func(e E) bool {
			yielded = true
			return yield(e)
		})
		if !yielded {
			if !yield(defaultValue) {
				return
			}
		}
	}, nil, nil, hints)
}
