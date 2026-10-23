package internal_test

import (
	"cmp"
	"slices"
	"testing"

	"example.com/playground/internal"
)

func TestHeap2(t *testing.T) {
	t.Run("natural order", func(t *testing.T) {
		h := internal.NewHeap2[int](cmp.Compare[int], -1)
		h.Push(2)
		h.Push(4)
		h.Push(1)

		var actual []int
		for h.Len() > 0 {
			actual = append(actual, h.Pop())
		}

		expected := []int{1, 2, 4}
		if !slices.Equal(expected, actual) {
			t.Fatalf("expected %v but got %v\n", expected, actual)
		}
	})

	t.Run("reverse order", func(t *testing.T) {
		h := internal.NewHeap2[int](func(a, b int) int { return -cmp.Compare[int](a, b) }, -1)
		h.Push(2)
		h.Push(4)
		h.Push(1)

		var actual []int
		for h.Len() > 0 {
			actual = append(actual, h.Pop())
		}

		expected := []int{4, 2, 1}
		if !slices.Equal(expected, actual) {
			t.Fatalf("expected %v but got %v\n", expected, actual)
		}
	})
}

func TestHeap2_Stability(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8}
	h := internal.NewStableHeap2[int](func(a, b int) int { return 0 }, -1, internal.DefaultStableHeap2Opts())
	for _, e := range s {
		h.Push(e)
	}
	s2 := make([]int, 0, len(s))
	for h.Len() > 0 {
		s2 = append(s2, h.Pop())
	}
	if !slices.Equal(s, s2) {
		t.Fatalf("expected %v, but got %v", s, s2)
	}
}

func TestHeap3_Stability(t *testing.T) {
	t.Run("stable", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5, 6, 7, 8}
		h := internal.NewStableHeap3[int](func(a, b int) int { return 0 }, -1, internal.StableHeap3Opts{})
		for _, e := range s {
			h.Push(e)
		}
		s2 := make([]int, 0, len(s))
		for h.Len() > 0 {
			s2 = append(s2, h.Pop())
		}
		if !slices.Equal(s, s2) {
			t.Fatalf("expected %v, but got %v", s, s2)
		}
	})

	t.Run("reverse stable", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5, 6, 7, 8}
		h := internal.NewStableHeap3[int](func(a, b int) int { return 0 }, -1, internal.StableHeap3Opts{
			IsUsingReverseInsertionOrderAsTiebreaker: true,
		})
		for _, e := range s {
			h.Push(e)
		}
		s2 := make([]int, 0, len(s))
		for h.Len() > 0 {
			s2 = append(s2, h.Pop())
		}
		slices.Reverse(s)
		if !slices.Equal(s, s2) {
			t.Fatalf("expected %v, but got %v", s, s2)
		}
	})
}
