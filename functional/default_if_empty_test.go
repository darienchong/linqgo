package functional_test

import (
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_DefaultIfEmpty(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		v := functional.FromSlice([]int{1, 2, 3}).
			DefaultIfEmpty(-1).
			Last(999)
		expected := 3
		if expected != v {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("empty", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3}).
			Where(func(i int) bool { return i > 10 }).
			DefaultIfEmpty(-1).
			Last(999)
		expected := -1
		if expected != v {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}
