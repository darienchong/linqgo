package functional_test

import (
	"cmp"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_Reverse(t *testing.T) {
	t.Run("unordered,unfiltered", func(t *testing.T) {
		v := functional.FromSlice([]int{1, 2, 3, 4, 5}).
			Reverse().
			ToSlice()
		expected := []int{5, 4, 3, 2, 1}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("unordered,filtered", func(t *testing.T) {
		v := functional.FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool { return i%2 == 0 }).
			Reverse().
			ToSlice()
		expected := []int{4, 2}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("ordered,unfiltered", func(t *testing.T) {
		v := functional.FromSlice([]int{1, 2, 3, 4, 5}).
			OrderBy(func(a, b int) int { return cmp.Compare[int](a%2, b%2) }). // [2 4 1 3 5]
			Reverse().                                                         // [5 3 1 4 2]
			ToSlice()
		expected := []int{5, 3, 1, 4, 2}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("ordered,filtered", func(t *testing.T) {
		v := functional.FromSlice([]int{1, 3, 5, 2, 4}).
			OrderBy(cmp.Compare[int]).
			Where(func(i int) bool { return i%2 == 1 }).
			Reverse().
			ToSlice()
		expected := []int{5, 3, 1}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("ordering-post-reversal-should-be-stable", func(t *testing.T) {
		v := functional.FromSlice([]int{1, 2, 3, 4, 5}).
			Reverse().                                                                   // [5 4 3 2 1]
			OrderByDescending(func(a, b int) int { return cmp.Compare[int](a%2, b%2) }). // [5 3 1 4 2]
			ToSlice()
		expected := []int{5, 3, 1, 4, 2}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}
