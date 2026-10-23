package functional

func (f *Functioner[E]) Reverse() *Functioner[E] {
	hints := f.hints.clone()
	hints.slice = nil

	return newF(
		f.sequencer(sequencerOptions{
			IsOrdered:                   true,
			ShouldIterateInReverseOrder: true,
		}),
		nil,
		nil,
		hints,
	)
}
