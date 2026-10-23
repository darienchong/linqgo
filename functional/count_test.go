package functional_test

import (
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_Count(t *testing.T) {
	t.Run("empty-from-slice", func(t *testing.T) {
		v := functional.FromSlice([]int{}).Count()
		expected := 0
		if expected != v {
			t.Fatalf("expected %v but got %v\n", expected, v)
		}
	})

	t.Run("from-slice", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		v := functional.FromSlice(input).Count()
		expected := len(input)
		if expected != v {
			t.Fatalf("expected %v but got %v\n", expected, v)
		}
	})

	t.Run("empty-from-seq", func(t *testing.T) {
		input := []int{}
		v := functional.FromSeq(slices.Values(input)).Count()
		expected := len(input)
		if expected != v {
			t.Fatalf("expected %v but got %v\n", expected, v)
		}
	})

	t.Run("from-seq", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		v := functional.FromSeq(slices.Values(input)).Count()
		expected := len(input)
		if expected != v {
			t.Fatalf("expected %v but got %v\n", expected, v)
		}
	})

	t.Run("where", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		v := functional.
			FromSlice(input).
			Where(func(i int) bool { return i%2 == 0 }).
			Count()
		expected := 2
		if expected != v {
			t.Fatalf("expected %v but got %v\n", expected, v)
		}
	})
}
