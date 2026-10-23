package utils

import (
	"container/heap"
)

var _ heap.Interface = (*Heap[int])(nil)

type Heap[T any] struct {
	backingSlice []T
	lessFn       func(a, b T) bool
}

// New returns a new heap instance using the given function to order the elements.
func New[T any](lessFn func(a, b T) bool) *Heap[T] {
	heapImpl := &Heap[T]{
		lessFn: lessFn,
	}
	heap.Init(heapImpl)
	return heapImpl
}

// Len returns the number of elements inside the heap.
func (h *Heap[T]) Len() int {
	return len(h.backingSlice)
}

// Deprecated: Less should not be used, except by container/heap.
func (h *Heap[T]) Less(i, j int) bool {
	return h.lessFn(h.backingSlice[i], h.backingSlice[j])
}

// Deprecated: Swap should not be used, except by container/heap.
func (h *Heap[T]) Swap(i, j int) {
	h.backingSlice[i], h.backingSlice[j] = h.backingSlice[j], h.backingSlice[i]
}

// Deprecated: Push should not be used, except by container/heap. Users looking to push elements into the heap should instead use
// container/heap.Push.
func (h *Heap[T]) Push(x any) {
	t, ok := x.(T)
	if !ok {
		return
	}

	h.backingSlice = append(h.backingSlice, t)
}

// Deprecated: Pop should not be used, except by container/heap. Users looking to pop elements from the heap should instead use
// container/heap.Pop.
func (h *Heap[T]) Pop() any {
	last := h.backingSlice[len(h.backingSlice)-1]
	h.backingSlice = h.backingSlice[:h.Len()-1]
	return last
}

// Peek returns the top of the heap. Returns nil if the heap is empty.
func (h *Heap[T]) Peek() *T {
	if h.Len() == 0 {
		return nil
	}

	return &h.backingSlice[0]
}
