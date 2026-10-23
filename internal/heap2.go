package internal

import (
	"iter"
)

// Heap2 implements a min-heap over arbitrary types using custom comparator
type Heap2[E any] struct {
	elements []E
	cmp      func(E, E) int
}

type Heap2Opts struct {
}

func NewHeap2[E any](cmp func(E, E) int, cap int) *Heap2[E] {
	h := &Heap2[E]{
		cmp: cmp,
	}
	if cap > 0 {
		h.elements = make([]E, 0, cap)
	}
	return h
}

func (h *Heap2[E]) Init() {
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

func (h *Heap2[E]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func (h *Heap2[E]) down(i0, n int) bool {
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

func (h *Heap2[E]) LoadSeq(seq iter.Seq[E]) *Heap2[E] {
	for e := range seq {
		h.elements = append(h.elements, e)
	}
	h.Init()
	return h
}

func (h *Heap2[E]) LoadSeqWithPredicate(seq iter.Seq[E], predicate func(E) bool) *Heap2[E] {
	for e := range seq {
		if !predicate(e) {
			continue
		}
		h.elements = append(h.elements, e)
	}
	h.Init()
	return h
}

func (h *Heap2[E]) LoadSlice(slice []E) *Heap2[E] {
	if cap(h.elements)-len(h.elements) < len(slice) {
		h.elements = append(make([]E, 0, len(h.elements)+len(slice)), h.elements...)
	}
	for _, e := range slice {
		h.elements = append(h.elements, e)
	}
	h.Init()
	return h
}

func (h *Heap2[E]) Len() int {
	return len(h.elements)
}

func (h *Heap2[E]) Less(i, j int) bool {
	return h.cmp(h.elements[i], h.elements[j]) < 0
}

func (h *Heap2[E]) Swap(i, j int) {
	t := h.elements[i]
	h.elements[i] = h.elements[j]
	h.elements[j] = t
}

func (h *Heap2[E]) push(e E) {
	h.elements = append(h.elements, e)
}

func (h *Heap2[E]) pop() E {
	e := h.elements[len(h.elements)-1]
	h.elements = h.elements[:len(h.elements)-1]
	return e
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap2[E]) Push(x E) {
	h.push(x)
	h.up(h.Len() - 1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to [Remove](h, 0).
func (h *Heap2[E]) Pop() E {
	n := h.Len() - 1
	h.Swap(0, n)
	h.down(0, n)
	return h.pop()
}

func (h *Heap2[E]) Peek() E {
	return h.elements[0]
}
