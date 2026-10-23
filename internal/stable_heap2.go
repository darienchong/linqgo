package internal

import (
	"cmp"
	"iter"
)

type heapElt[E any] struct {
	insertionOrder int
	e              E
}

func newHeapElt[E any](insertionOrder int, e E) heapElt[E] {
	return heapElt[E]{
		insertionOrder: insertionOrder,
		e:              e,
	}
}

// StableHeap2 implements a stable min-heap over arbitrary types using custom comparator
type StableHeap2[E any] struct {
	Heap2[heapElt[E]]
	currInsertionOrder int
}

func stableComparer[E any](comparer func(E, E) int, opts StableHeap2Opts) func(heapElt[E], heapElt[E]) int {
	return func(a, b heapElt[E]) int {
		if v := comparer(a.e, b.e); v != 0 {
			return v
		}

		if opts.IsUsingReverseInsertionOrderAsTiebreaker {
			return -cmp.Compare(a.insertionOrder, b.insertionOrder)

		}

		return cmp.Compare(a.insertionOrder, b.insertionOrder)
	}
}

type StableHeap2Opts struct {
	// IsUsingReverseInsertionOrderAsTiebreaker controls the relative emission order of equal elements
	// when IsStable is true. If true (and IsStable is true), equal elements inserted later will be emitted
	// before equal elements inserted earlier.
	IsUsingReverseInsertionOrderAsTiebreaker bool
}

var (
	defaultStableHeapOpts StableHeap2Opts
)

func DefaultStableHeap2Opts() StableHeap2Opts {
	return defaultStableHeapOpts
}

func NewStableHeap2[E any](cmp func(E, E) int, cap int, opts StableHeap2Opts) *StableHeap2[E] {
	return &StableHeap2[E]{
		Heap2: *NewHeap2[heapElt[E]](stableComparer[E](cmp, opts), cap),
	}
}

func (h *StableHeap2[E]) LoadSeq(seq iter.Seq[E]) *StableHeap2[E] {
	for e := range seq {
		h.Heap2.elements = append(h.elements, newHeapElt[E](h.currInsertionOrder, e))
		h.currInsertionOrder++
	}
	h.currInsertionOrder++
	h.Heap2.Init()
	return h
}

func (h *StableHeap2[E]) LoadSeqWithPredicate(seq iter.Seq[E], predicate func(E) bool) *StableHeap2[E] {
	for e := range seq {
		if !predicate(e) {
			continue
		}
		h.Heap2.elements = append(h.elements, newHeapElt[E](h.currInsertionOrder, e))
		h.currInsertionOrder++
	}
	h.currInsertionOrder++
	h.Heap2.Init()
	return h
}

func (h *StableHeap2[E]) LoadSlice(slice []E) *StableHeap2[E] {
	if cap(h.Heap2.elements)-len(h.Heap2.elements) < len(slice) {
		h.Heap2.elements = append(make([]heapElt[E], 0, len(h.Heap2.elements)+len(slice)), h.Heap2.elements...)
	}
	for _, e := range slice {
		h.Heap2.elements = append(h.Heap2.elements, newHeapElt[E](h.currInsertionOrder, e))
		h.currInsertionOrder++
	}
	h.currInsertionOrder++
	h.Heap2.Init()
	return h
}

func (h *StableHeap2[E]) LoadSliceWithPredicate(slice []E, predicate func(E) bool) *StableHeap2[E] {
	if cap(h.Heap2.elements)-len(h.Heap2.elements) < len(slice) {
		h.Heap2.elements = append(make([]heapElt[E], 0, len(h.Heap2.elements)+len(slice)), h.Heap2.elements...)
	}
	for _, e := range slice {
		if !predicate(e) {
			continue
		}
		h.Heap2.elements = append(h.Heap2.elements, newHeapElt[E](h.currInsertionOrder, e))
		h.currInsertionOrder++
	}
	h.currInsertionOrder++
	h.Heap2.Init()
	return h
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *StableHeap2[E]) Push(x E) {
	h.Heap2.Push(newHeapElt(h.currInsertionOrder, x))
	h.currInsertionOrder++
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to [Remove](h, 0).
func (h *StableHeap2[E]) Pop() E {
	return h.Heap2.Pop().e
}

func (h *StableHeap2[E]) Peek() E {
	return h.Heap2.Peek().e
}
