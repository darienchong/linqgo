package functional

import (
	"iter"
	"slices"

	"example.com/playground/internal"
)

type sequencerOptions struct {
	ShouldIterateInReverseOrder bool
	IsOrdered                   bool
	MaxSize                     int
}

func (f *Functioner[E]) sequencer(options sequencerOptions) iter.Seq[E] {
	if options.MaxSize == 1 {
		return singleSequencer(f, options.ShouldIterateInReverseOrder)
	}
	if options.MaxSize > 0 {
		return boundedSequencer(f, options.MaxSize, options.ShouldIterateInReverseOrder)
	}
	return fullSequencer(f, options.IsOrdered, options.ShouldIterateInReverseOrder)
}

func singleSequencer[E any](f *Functioner[E], shouldReverseComparer bool) iter.Seq[E] {
	shouldCompare := f.shouldCompare()
	shouldFilter := f.shouldFilter()
	p := f.filter()

	if shouldCompare {
		cd := comparerForward
		if shouldReverseComparer {
			cd = comparerReverse
		}
		// For small values of maxSize, boundedSequencer(maxSize) can offer better space complexity
		// than loading the entire sequence at once into the heap
		return func(yield func(E) bool) {
			c := f.comparer(cd)
			init := false
			var singleElt E
			processElt := func(e E) {
				if shouldFilter && !p(e) {
					return
				}
				if !init {
					singleElt = e
					init = true
					return
				}
				if (shouldReverseComparer && c(e, singleElt) <= 0) || (!shouldReverseComparer && c(e, singleElt) < 0) {
					singleElt = e
				}
			}
			if f.hints.slice != nil {
				for i := range f.hints.slice {
					processElt(f.hints.slice[i])
				}
			} else {
				for e := range f.sequence {
					processElt(e)
				}
			}
			if init {
				if !yield(singleElt) {
					return
				}
			}
		}
	}

	if /* !shouldCompare && */ shouldFilter && shouldReverseComparer {
		var last E
		init := false
		return func(yield func(E) bool) {
			if f.hints.slice != nil {
				for i := len(f.hints.slice) - 1; i >= 0; i-- {
					e := f.hints.slice[i]
					if !p(e) {
						continue
					}
					if !init {
						init = true
					}
					last = e
					break
				}
			} else {
				for e := range f.sequence {
					if !p(e) {
						continue
					}
					if !init {
						init = true
					}
					last = e
				}
			}

			if init {
				yield(last)
			}
		}
	}

	if shouldReverseComparer {
		if f.hints.slice != nil {
			return func(yield func(E) bool) {
				yield(f.hints.slice[len(f.hints.slice)-1])
			}
		}

		return func(yield func(E) bool) {
			var last E
			init := false
			for e := range f.sequence {
				last = e
				init = true
			}
			if init {
				yield(last)
			}
		}
	}

	return func(yield func(E) bool) {
		processElt := func(e E) bool {
			if shouldFilter && !p(e) {
				return true
			}
			yield(e)
			return false
		}
		if f.hints.slice != nil {
			for i := range f.hints.slice {
				if !processElt(f.hints.slice[i]) {
					return
				}
			}
		} else {
			for e := range f.sequence {
				if !processElt(e) {
					return
				}
			}
		}
	}
}

func boundedSequencer[E any](f *Functioner[E], maxSize int, shouldReverseComparer bool) iter.Seq[E] {
	shouldCompare := f.shouldCompare()
	shouldFilter := f.shouldFilter()
	p := f.filter()

	var resolvedComparerDirection, reverseResolvedComparerDirection comparerDirection
	if shouldReverseComparer {
		resolvedComparerDirection = comparerReverse
		reverseResolvedComparerDirection = comparerForward
	} else {
		resolvedComparerDirection = comparerForward
		reverseResolvedComparerDirection = comparerReverse
	}

	if shouldCompare {
		// For small values of maxSize, boundedSequencer(maxSize) can offer better space complexity
		// than loading the entire sequence at once into the heap
		return func(yield func(E) bool) {
			c := f.comparer(resolvedComparerDirection)
			h := internal.NewStableHeap3[E](
				f.comparer(reverseResolvedComparerDirection),
				maxSize,
				internal.StableHeap3Opts{
					IsUsingReverseInsertionOrderAsTiebreaker: true,
				},
			)
			processElt := func(e E) {
				if shouldFilter && !p(e) {
					return
				}
				if h.Len() < maxSize {
					h.Push(e)
					return
				}
				if c(e, h.Peek()) < 0 {
					h.Pop()
					h.Push(e)
				}
			}
			if f.hints.slice != nil {
				for i := range f.hints.slice {
					processElt(f.hints.slice[i])
				}
			} else {
				for e := range f.sequence {
					processElt(e)
				}
			}
			// This is in reverse order, so we need to collect it and traverse in reverse order
			s := make([]E, 0, maxSize)
			for h.Len() > 0 {
				s = append(s, h.Pop())
			}
			//goland:noinspection GoFixSlicesBackward
			for i := len(s) - 1; i >= 0; i-- {
				if !yield(s[i]) {
					return
				}
			}
		}
	}

	if shouldReverseComparer {
		return func(yield func(E) bool) {
			b := internal.NewRingBuffer[E](maxSize)
			if f.hints.slice != nil {
				for i := range f.hints.slice {
					b.Push(f.hints.slice[i])
				}
			} else {
				for e := range f.sequence {
					b.Push(e)
				}
			}
			b.ReverseIterate2(func(e E) bool {
				if shouldFilter && !p(e) {
					return true
				}
				if !yield(e) {
					return false
				}
				return true
			})
		}
	}

	return func(yield func(E) bool) {
		counter := maxSize
		if counter <= 0 {
			return
		}
		processElt := func(e E) bool {
			if shouldFilter && !p(e) {
				return true
			}
			if !yield(e) {
				return false
			}
			counter--
			if counter <= 0 {
				return false
			}
			return true
		}
		if f.hints.slice != nil {
			for i := range f.hints.slice {
				if !processElt(f.hints.slice[i]) {
					return
				}
			}
		} else {
			for e := range f.sequence {
				if !processElt(e) {
					return
				}
			}
		}
	}
}

func fullSequencer[E any](f *Functioner[E], isOrdered, shouldIterateInReverseOrder bool) iter.Seq[E] {
	shouldCompare := isOrdered && f.shouldCompare()
	shouldFilter := f.shouldFilter()

	if shouldCompare {
		return func(yield func(E) bool) {
			var h *internal.StableHeap3[E]
			if shouldIterateInReverseOrder {
				h = internal.NewStableHeap3[E](
					f.comparer(comparerReverse),
					f.hints.lengthEstimate,
					internal.StableHeap3Opts{
						IsUsingReverseInsertionOrderAsTiebreaker: true,
					})
			} else {
				h = internal.NewStableHeap3[E](
					f.comparer(comparerForward),
					f.hints.lengthEstimate,
					internal.DefaultStableHeap3Opts(),
				)
			}
			if f.hints.slice != nil {
				if shouldFilter {
					h.LoadSliceWithPredicate(f.hints.slice, f.filter())
				} else {
					h.LoadSlice(f.hints.slice)
				}
			} else {
				if shouldFilter {
					h.LoadSeqWithPredicate(f.sequence, f.filter())
				} else {
					h.LoadSeq(f.sequence)
				}
			}
			for h.Len() > 0 {
				if !yield(h.Pop()) {
					return
				}
			}
		}
	} else if /* !shouldCompare && */ shouldFilter {
		p := f.filter()
		return func(yield func(E) bool) {
			processElt := func(e E) bool {
				if !p(e) {
					return true
				}
				if !yield(e) {
					return false
				}
				return true
			}
			if f.hints.slice != nil {
				if shouldIterateInReverseOrder {
					for i := len(f.hints.slice) - 1; i >= 0; i-- {
						if !processElt(f.hints.slice[i]) {
							return
						}
					}
				} else {
					for i := range f.hints.slice {
						if !processElt(f.hints.slice[i]) {
							return
						}
					}
				}

			} else {
				if shouldIterateInReverseOrder {
					s := slices.Collect(f.sequence)
					for i := len(s) - 1; i >= 0; i-- {
						if !processElt(s[i]) {
							return
						}
					}
				} else {
					for e := range f.sequence {
						if !processElt(e) {
							return
						}
					}
				}

			}
		}
	}

	// !shouldCompare && !shouldFilter
	if shouldIterateInReverseOrder {
		return func(yield func(E) bool) {
			s := slices.Collect(f.sequence)
			for i := len(s) - 1; i >= 0; i-- {
				if !yield(s[i]) {
					return
				}
			}
		}
	}

	return f.sequence
}
