package functional_test

import (
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_Zip(t *testing.T) {
	t.Run("zip", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3})
		f2 := functional.FromSlice([]int{4, 5, 6})
		v := f1.Zip(f2, func(a, b int) int { return a*10 + b }).ToSlice()
		expected := []int{14, 25, 36}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}
