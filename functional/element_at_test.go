package functional_test

import (
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_ElementAt(t *testing.T) {
	t.Run("in-range", func(t *testing.T) {
		v, exists := functional.
			FromSlice([]int{0, 1, 2, 3, 4}).
			ElementAt(2)
		if v != 2 || !exists {
			t.Fatalf("expected (%d, %t), got (%d, %t)\n", 2, true, v, exists)
		}
	})

	t.Run("out-of-range", func(t *testing.T) {
		query := functional.FromSlice([]int{1, 2, 3})

		for _, index := range []int{-1, 3, 100} {
			v, exists := query.ElementAt(index)
			if exists || v != 0 {
				t.Fatalf("ElementAt(%d)=(%d,%t), expected (0,false)", index, v, exists)
			}
		}
	})
}
