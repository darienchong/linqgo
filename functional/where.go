package functional

func (f *Functioner[E]) Where(predicate func(E) bool) *Functioner[E] {
	hints := f.hints.clone()
	hints.slice = nil
	hints.op = opWhere
	if f.hints.lengthEstimate > 0 {
		hints.lengthEstimate = f.hints.lengthEstimate >> 1 // Use 1/2 of the elements as an estimate for allocation
		if hints.lengthEstimate == 0 {
			hints.lengthEstimate = -1
		}
	}

	return newF(f.sequence, copyOnAppend(f.filters, predicate), f.comparers, hints)
}
