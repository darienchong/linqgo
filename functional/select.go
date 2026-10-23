package functional

type SelectFunctioner[E, V any] struct {
	*Functioner[V]
	source    *Functioner[E]
	projector func(E) V
}

func (f *SelectFunctioner[E, V]) Select[V2 any](projector func(V) V2) *SelectFunctioner[E, V2] {
	return &SelectFunctioner[E, V2]{
		Functioner: newF(func(yield func(V2) bool) {
			if f.source.shouldCompare() {
				s := f.source.slice(true, true)
				for i := range s {
					if !yield(projector(f.projector(s[i]))) {
						return
					}
				}
			} else {
				f.source.sequencer(sequencerOptions{IsOrdered: false})(func(e E) bool {
					return yield(projector(f.projector(e)))
				})
			}
		}, nil, nil, f.hints.clone2[V2]()),
		source:    f.source,
		projector: func(e E) V2 { return projector(f.projector(e)) },
	}
}

func (f *Functioner[E]) Select[V any](projector func(E) V) *SelectFunctioner[E, V] {
	hints := f.hints.clone2[V]()
	hints.op = opSelect

	return &SelectFunctioner[E, V]{
		Functioner: newF(func(yield func(V) bool) {
			if f.shouldCompare() {
				s := f.slice(true, true)
				for i := range s {
					if !yield(projector(s[i])) {
						return
					}
				}
			} else {
				f.sequencer(sequencerOptions{IsOrdered: true})(func(e E) bool {
					return yield(projector(e))
				})
			}
		}, nil, nil, f.hints.clone2[V]()),
		source:    f,
		projector: projector,
	}
}

func (f *Functioner[E]) SelectMany[V any](projector func(E) []V) *Functioner[V] {
	hints := f.hints.clone2[V]()
	hints.lengthEstimate = -1

	return newF(func(yield func(V) bool) {
		if f.shouldCompare() {
			s := f.slice(true, true)
			for i := range s {
				for _, p := range projector(s[i]) {
					if !yield(p) {
						return
					}
				}
			}
		} else {
			f.sequencer(sequencerOptions{IsOrdered: true})(func(e E) bool {
				projected := projector(e)
				for i := range projected {
					if !yield(projected[i]) {
						return false
					}
				}
				return true
			})
		}
	}, nil, nil, hints)
}
