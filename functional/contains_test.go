package functional_test

import (
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_Contains(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		v := functional.FromSlice(input).Contains(3, func(a, b int) bool { return a == b })
		const expected = true
		if expected != v {
			t.Fatalf("expected %v, but got %v\n", expected, v)
		}
	})

	t.Run("not-contains", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		v := functional.FromSlice(input).Contains(999, func(a, b int) bool { return a == b })
		const expected = false
		if expected != v {
			t.Fatalf("expected %v, but got %v\n", expected, v)
		}
	})
}
