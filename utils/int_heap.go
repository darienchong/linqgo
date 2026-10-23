package utils

import "container/heap"

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	e := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return e
}

func (h *IntHeap) Peek() int {
	return (*h)[0]
}

var _ heap.Interface = (*IntHeap)(nil)
