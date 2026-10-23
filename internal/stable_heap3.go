package internal

import "iter"

type StableHeap3Opts struct {
	// IsUsingReverseInsertionOrderAsTiebreaker controls the relative emission order of equal elements
	// when IsStable is true. If true (and IsStable is true), equal elements inserted later will be emitted
	// before equal elements inserted earlier.
	IsUsingReverseInsertionOrderAsTiebreaker bool
}

var (
	defaultStableHeap3Opts StableHeap3Opts
)

func DefaultStableHeap3Opts() StableHeap3Opts {
	return defaultStableHeap3Opts
}

type StableHeap3[E any] struct {
	cmp                func(E, E) int
	opts               StableHeap3Opts
	elements           []E
	indices            []int
	currInsertionOrder int
}

func NewStableHeap3[E any](cmp func(E, E) int, cap int, opts StableHeap3Opts) *StableHeap3[E] {
	h := &StableHeap3[E]{
		cmp:  cmp,
		opts: opts,
	}

	if cap > 0 {
		h.elements = make([]E, 0, cap)
		h.indices = make([]int, 0, cap)
	}

	return h
}

func (h *StableHeap3[E]) Len() int {
	return len(h.elements)
}

func (h *StableHeap3[E]) Less(i, j int) bool {
	if v := h.cmp(h.elements[i], h.elements[j]); v != 0 {
		if v < 0 {
			return true
		}
		return false
	}

	if h.opts.IsUsingReverseInsertionOrderAsTiebreaker {
		return h.indices[j] < h.indices[i]
	}
	return h.indices[i] < h.indices[j]
}

func (h *StableHeap3[E]) Swap(i, j int) {
	h.elements[i], h.elements[j] = h.elements[j], h.elements[i]
	h.indices[i], h.indices[j] = h.indices[j], h.indices[i]
}

func (h *StableHeap3[E]) LoadSeq(seq iter.Seq[E]) *StableHeap3[E] {
	for e := range seq {
		h.elements = append(h.elements, e)
		h.indices = append(h.indices, h.currInsertionOrder)
		h.currInsertionOrder++
	}
	h.currInsertionOrder++
	h.Init()
	return h
}

func (h *StableHeap3[E]) LoadSeqWithPredicate(seq iter.Seq[E], predicate func(E) bool) *StableHeap3[E] {
	for e := range seq {
		if predicate != nil && !predicate(e) {
			continue
		}
		h.elements = append(h.elements, e)
		h.indices = append(h.indices, h.currInsertionOrder)
		h.currInsertionOrder++
	}
	h.currInsertionOrder++
	h.Init()
	return h
}

func (h *StableHeap3[E]) LoadSlice(slice []E) *StableHeap3[E] {
	if cap(h.elements)-len(h.elements) < len(slice) {
		h.elements = append(make([]E, 0, len(h.elements)+len(slice)), h.elements...)
	}
	for _, e := range slice {
		h.elements = append(h.elements, e)
		h.indices = append(h.indices, h.currInsertionOrder)
		h.currInsertionOrder++
	}
	h.currInsertionOrder++
	h.Init()
	return h
}

func (h *StableHeap3[E]) LoadSliceWithPredicate(slice []E, predicate func(E) bool) *StableHeap3[E] {
	if cap(h.elements)-len(h.elements) < len(slice) {
		h.elements = append(make([]E, 0, len(h.elements)+len(slice)), h.elements...)
	}
	for _, e := range slice {
		if predicate != nil && !predicate(e) {
			continue
		}
		h.elements = append(h.elements, e)
		h.indices = append(h.indices, h.currInsertionOrder)
		h.currInsertionOrder++
	}
	h.currInsertionOrder++
	h.Init()
	return h
}

func (h *StableHeap3[E]) Init() {
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

func (h *StableHeap3[E]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func (h *StableHeap3[E]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *StableHeap3[E]) Push(x E) {
	h.elements = append(h.elements, x)
	h.indices = append(h.indices, h.currInsertionOrder)
	h.currInsertionOrder++
	h.up(h.Len() - 1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to [Remove](h, 0).
func (h *StableHeap3[E]) Pop() E {
	n := h.Len() - 1
	h.Swap(0, n)
	h.down(0, n)
	e := h.elements[len(h.elements)-1]
	h.elements = h.elements[:len(h.elements)-1]
	h.indices = h.indices[:len(h.indices)-1]
	return e
}

func (h *StableHeap3[E]) Peek() E {
	return h.elements[0]
}
