package functional

import (
	"errors"
	"fmt"
	"iter"
)

func (f *Functioner[E]) Cache() *Functioner[E] {
	return FromSlice(f.ToSlice())
}

func (f *Functioner[E]) All(predicate func(E) bool) bool {
	// Don't use slice() or sequencer() here since we specifically want to integrate
	// filtering from Where() and predicate testing with Any() into a single loop
	shouldFilter := f.shouldFilter()
	var filter filter[E]
	if shouldFilter {
		filter = f.filter()
	}
	if f.hints.slice != nil {
		for i := range f.hints.slice {
			e := f.hints.slice[i]
			if shouldFilter {
				if !filter(e) {
					continue
				}

				if !predicate(e) {
					return false
				}
			} else {
				if !predicate(e) {
					return false
				}
			}
		}
	} else {
		for e := range f.sequence {
			if shouldFilter {
				if !filter(e) {
					continue
				}
				if !predicate(e) {
					return false
				}
			} else {
				if !predicate(e) {
					return false
				}
			}
		}
	}

	return true
}

func (f *Functioner[E]) Any(predicate func(E) bool) bool {
	// Don't use slice() or sequencer() here since we specifically want to integrate
	// filtering from Where() and predicate testing with Any() into a single loop
	shouldFilter := f.shouldFilter()
	var filter filter[E]
	if shouldFilter {
		filter = f.filter()
	}

	if f.hints.slice != nil {
		for i := range f.hints.slice {
			e := f.hints.slice[i]
			if shouldFilter {
				if !filter(e) {
					continue
				}
				if predicate(e) {
					return true
				}
			} else {
				if predicate(e) {
					return true
				}
			}
		}
		return false
	}

	for e := range f.sequence {
		if shouldFilter {
			if !filter(e) {
				continue
			}
			if predicate(e) {
				return true
			}
		} else {
			if predicate(e) {
				return true
			}
		}
	}
	return false
}

func (f *Functioner[E]) Aggregate[V any](seed V, aggregator func(v V, next E) V) V {
	if !f.shouldCompare() && f.hints.slice != nil {
		shouldFilter := f.shouldFilter()
		var p filter[E]
		if shouldFilter {
			p = f.filter()
		}
		agg := seed
		for i := range f.hints.slice {
			e := f.hints.slice[i]
			if shouldFilter && !p(e) {
				continue
			}
			agg = aggregator(agg, e)
		}
		return agg
	}

	if f.shouldCompare() {
		agg := seed
		s := f.slice(true, true)
		for i := range s {
			agg = aggregator(agg, s[i])
		}
		return agg
	}

	agg := seed
	f.sequencer(sequencerOptions{IsOrdered: false})(func(e E) bool {
		agg = aggregator(agg, e)
		return true
	})
	return agg
}

func (f *Functioner[E]) ToSequence() iter.Seq[E] {
	return f.sequencer(sequencerOptions{IsOrdered: true})
}

func (f *Functioner[E]) ToSlice() []E {
	return f.slice(true, false)
}

func (f *Functioner[E]) ToMap[K comparable, V any](keyer func(E) K, valuer func(E) V) map[K]V {
	m := make(map[K]V, max(f.hints.lengthEstimate, 0))

	shouldFilter := f.shouldFilter()
	var filter filter[E]
	if shouldFilter {
		filter = f.filter()
	}
	if !f.shouldCompare() && f.hints.slice != nil {
		for i := range f.hints.slice {
			e := f.hints.slice[i]
			if shouldFilter {
				if !filter(e) {
					continue
				}
			}
			m[keyer(e)] = valuer(e)
		}
	} else if f.shouldCompare() {
		s := f.slice(true, true)
		for i := range s {
			m[keyer(s[i])] = valuer(s[i])
		}
	} else {
		f.sequencer(sequencerOptions{IsOrdered: false})(func(e E) bool {
			m[keyer(e)] = valuer(e)
			return true
		})
	}

	return m
}

func (f *Functioner[E]) Chunk[S ~[]E](chunkSize int) *Functioner[S] {
	if chunkSize <= 0 {
		panic(fmt.Sprintf("cannot Chunk() with non-positive chunk size: %d", chunkSize))
	}
	return newF[S](func(yield func(S) bool) {
		var chunk []E
		if f.shouldCompare() {
			s := f.slice(true, true)
			for i := range s {
				if chunk == nil {
					chunk = make([]E, 0, chunkSize)
				}
				chunk = append(chunk, s[i])
				if len(chunk) == chunkSize {
					if !yield(chunk) {
						return
					}
					chunk = nil
				}
			}
		} else {
			for e := range f.sequencer(sequencerOptions{IsOrdered: false}) {
				if chunk == nil {
					chunk = make([]E, 0, chunkSize)
				}
				chunk = append(chunk, e)
				if len(chunk) == chunkSize {
					if !yield(chunk) {
						return
					}
					chunk = nil
				}
			}
		}
		if chunk != nil && !yield(chunk) {
			return
		}
	}, nil, nil, newHints[S]())
}

var (
	SingleError = errors.New("invoked Single() but the sequence did not have exactly 1 element")
)

func (f *Functioner[E]) Single() (E, error) {
	var single E
	var init bool
	for e := range f.sequencer(sequencerOptions{
		IsOrdered: true,
		MaxSize:   2,
	}) {
		if init {
			return single, SingleError
		}
		single = e
		init = true
	}

	if !init {
		return single, SingleError
	}
	return single, nil
}
