package functional

import (
	"iter"
	"slices"
)

type Functioner[E any] struct {
	sequence  iter.Seq[E]
	filters   []filter[E]
	comparers []comparer[E]
	hints     hints[E]
}

func newF[E any](
	sequence iter.Seq[E],
	filters []filter[E],
	comparers []comparer[E],
	hints hints[E],
) *Functioner[E] {
	return &Functioner[E]{
		sequence:  sequence,
		filters:   filters,
		comparers: comparers,
		hints:     hints,
	}
}

func (f *Functioner[E]) filter() filter[E] {
	if len(f.filters) == 0 {
		return nil
	}
	if len(f.filters) == 1 {
		return f.filters[0]
	}
	return func(e E) bool {
		for _, predicate := range f.filters {
			if !predicate(e) {
				return false
			}
		}
		return true
	}
}

type comparerDirection int

const (
	comparerForward comparerDirection = iota + 1
	comparerReverse
)

func (f *Functioner[E]) comparer(comparerDirection comparerDirection) comparer[E] {
	if len(f.comparers) == 0 {
		return func(E, E) int { return 0 }
	}
	if len(f.comparers) == 1 {
		switch comparerDirection {
		case comparerForward:
			return f.comparers[0]
		case comparerReverse:
			return func(e1, e2 E) int { return -f.comparers[0](e1, e2) }
		default:
			return f.comparers[0]
		}
	}
	return func(a, b E) int {
		//goland:noinspection GoFixSlicesBackward
		for i := len(f.comparers) - 1; i >= 0; i-- {
			if v := f.comparers[i](a, b); v != 0 {
				switch comparerDirection {
				case comparerForward:
					return v
				case comparerReverse:
					return -v
				default:
					return v
				}
			}
		}
		return 0
	}
}

func (f *Functioner[E]) slice(ordered bool, isReadOnly bool) []E {
	var s []E
	// We integrate materialisation and filtering to skip an additional loop for performance
	if len(f.filters) == 0 {
		if f.hints.slice != nil {
			if !(ordered && f.shouldCompare()) && isReadOnly {
				s = f.hints.slice
			} else {
				s = slices.Clone(f.hints.slice)
			}
		} else if f.hints.lengthEstimate > 0 {
			s = slices.AppendSeq(make([]E, 0, f.hints.lengthEstimate), f.sequence)
		} else if f.hints.lengthEstimate <= 0 {
			s = slices.Collect(f.sequence)
		}
	} else if f.shouldFilter() {
		filter := f.filter()
		if f.hints.slice != nil {
			// Don't preallocate the entire source slice's worth of space since we don't know how many elements will be filtered
			s = make([]E, 0, max(1, len(f.hints.slice)>>1))
			for _, e := range f.hints.slice {
				if !filter(e) {
					continue
				}
				s = append(s, e)
			}
		} else {
			s = make([]E, 0, max(1, f.hints.lengthEstimate))
			for e := range f.sequence {
				if !filter(e) {
					continue
				}
				s = append(s, e)
			}
		}
	}

	if ordered && f.shouldCompare() {
		slices.SortStableFunc(s, f.comparer(comparerForward))
	}

	return s
}
