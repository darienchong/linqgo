package functional

import "iter"

func (f *Functioner[E1]) Zip[E2, V any](f2 *Functioner[E2], projector func(E1, E2) V) *Functioner[V] {
	return newF(func(yield func(V) bool) {
		next, stop := iter.Pull(f.sequencer(sequencerOptions{IsOrdered: true}))
		next2, stop2 := iter.Pull(f2.sequencer(sequencerOptions{IsOrdered: true}))
		defer stop()
		defer stop2()
		for {
			e1, ok := next()
			if !ok {
				return
			}
			e2, ok := next2()
			if !ok {
				return
			}
			if !yield(projector(e1, e2)) {
				return
			}
		}
	}, nil, nil, newHints[V]())
}
