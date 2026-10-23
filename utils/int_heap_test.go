package utils

import (
	"container/heap"
	"testing"
)

func TestIntHeap(t *testing.T) {
	h := &IntHeap{}
	heap.Push(h, 5)
	if h.Peek() != 5 {
		t.Errorf("expect 5 but got %d", h.Peek())
	}

	heap.Push(h, 3)
	if h.Peek() != 3 {
		t.Errorf("expect 4 but got %d", h.Peek())
	}

	heap.Push(h, 4)
	if h.Peek() != 3 {
		t.Errorf("expect 3 but got %d", h.Peek())
	}

	for _, expectedE := range []int{3, 4, 5} {
		e := heap.Pop(h).(int)
		if e != expectedE {
			t.Errorf("expect %d but got %d", expectedE, e)
		}
	}
}
