package functional_test

import (
	"cmp"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_FirstBy(t *testing.T) {
	t.Run("empty defaults", func(t *testing.T) {
		empty := functional.FromSlice([]int{})
		if v := empty.First(42); v != 42 {
			t.Fatalf("expected %d, got %d\n", 42, v)
		}
		if v := empty.FirstBy(func(int) bool { return true }, 42); v != 42 {
			t.Fatalf("expected %d, got %d\n", 42, v)
		}
	})

	t.Run("non-default value", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			FirstBy(func(i int) bool { return i%2 == 0 }, -1)
		expected := 2
		if v != expected {
			t.Fatalf("expected %d, got %d\n", expected, v)
		}
	})

	t.Run("default value", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			FirstBy(func(i int) bool { return i%2 == 5 }, -1)
		expected := -1
		if v != expected {
			t.Fatalf("expected %d, got %d\n", expected, v)
		}
	})

	t.Run("order-first", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			OrderByDescending(cmp.Compare[int]).
			First(-1)
		expected := 4
		if v != expected {
			t.Fatalf("expected %d, got %d\n", expected, v)
		}
	})

	t.Run("order-first-stable", func(t *testing.T) {
		v := functional.
			FromSlice([][]int{
				{1, 1},
				{2, 2},
				{3, 3},
				{3, 4},
			}).
			OrderByDescending(func(a, b []int) int { return cmp.Compare[int](a[0], b[0]) }).
			First([]int{-1, -1})
		expected := []int{3, 3}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %d, got %d\n", expected, v)
		}
	})
}

func TestFunctioner_LastBy(t *testing.T) {
	t.Run("empty defaults", func(t *testing.T) {
		empty := functional.FromSlice([]int{})
		if v := empty.Last(42); v != 42 {
			t.Fatalf("expected %d, got %d\n", 42, v)
		}
		if v := empty.LastBy(func(int) bool { return true }, 42); v != 42 {
			t.Fatalf("expected %d, got %d\n", 42, v)
		}
	})

	t.Run("non-default value", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			LastBy(func(i int) bool { return i%2 == 0 }, -1)
		expected := 4
		if v != expected {
			t.Fatalf("expected %d, got %d\n", expected, v)
		}
	})

	t.Run("default value", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			LastBy(func(i int) bool { return i%2 == 5 }, -1)
		expected := -1
		if v != expected {
			t.Fatalf("expected %d, got %d\n", expected, v)
		}
	})

	t.Run("order-last-by", func(t *testing.T) {
		v := functional.
			FromSlice(shuffle([]int{1, 2, 3, 4})).
			OrderBy(cmp.Compare[int]).
			LastBy(func(i int) bool { return true }, -1)
		expected := 4
		if v != expected {
			t.Fatalf("expected %d, got %d\n", expected, v)
		}
	})

	t.Run("order-last-stable", func(t *testing.T) {
		v := functional.
			FromSlice([][]int{
				{1, 1},
				{2, 2},
				{3, 3},
				{3, 4},
			}).
			OrderBy(func(a, b []int) int { return cmp.Compare[int](a[0], b[0]) }).
			Last([]int{-1, -1})
		expected := []int{3, 4}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %d, got %d\n", expected, v)
		}
	})
}
