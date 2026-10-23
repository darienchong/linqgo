package functional_test

import (
	"cmp"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_Concat(t *testing.T) {
	v := functional.
		FromSlice([]int{1, 2, 3}).
		Concat(functional.FromSlice([]int{4, 5, 6})).
		ToSlice()
	expected := []int{1, 2, 3, 4, 5, 6}
	if !slices.Equal(expected, v) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}

	t.Run("early exit does not start second source", func(t *testing.T) {
		firstPulled := 0
		secondStarted := false
		first := functional.FromSeq(func(yield func(int) bool) {
			for value := 1; value <= 3; value++ {
				firstPulled++
				if !yield(value) {
					return
				}
			}
		})
		second := functional.FromSeq(func(yield func(int) bool) {
			secondStarted = true
			yield(4)
		})

		var values []int
		for value := range first.Concat(second).ToSequence() {
			values = append(values, value)
			if len(values) == 2 {
				break
			}
		}

		expected := []int{1, 2}
		if !slices.Equal(values, expected) {
			t.Fatalf("expected %v, got %v\n", expected, values)
		}
		if firstPulled != 2 {
			t.Fatalf("expected 2 pulls, got %d\n", firstPulled)
		}
		if secondStarted {
			t.Fatal("second source started after early exit")
		}
	})
}

func TestFunctioner_IntersectBy(t *testing.T) {
	t.Run("left-empty", func(t *testing.T) {
		v := functional.
			FromSlice([]int{}).
			IntersectBy(functional.FromSlice([]int{1, 2, 3}), func(i int) int { return i }).
			ToSlice()
		var expected []int
		if !equalIgnoringOrder(expected, v, cmp.Compare[int]) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("right-empty", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3}).
			IntersectBy(functional.FromSlice([]int{}), func(i int) int { return i }).
			ToSlice()
		var expected []int
		if !equalIgnoringOrder(expected, v, cmp.Compare[int]) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("left, right non-empty", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3}).
			IntersectBy(functional.FromSlice([]int{2, 3, 4}), func(i int) int { return i }).
			ToSlice()
		expected := []int{2, 3}
		if !equalIgnoringOrder(expected, v, cmp.Compare[int]) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("preserves left representatives, removes duplicates", func(t *testing.T) {
		v := functional.
			FromSlice([]int{5, 7, 8, 10}).
			IntersectBy(
				functional.FromSlice([]int{1, 4, 7, 9, 12, 3}),
				func(value int) int { return value % 2 },
			).
			ToSlice()
		expected := []int{5, 8}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_UnionBy(t *testing.T) {
	t.Run("left-empty", func(t *testing.T) {
		v := functional.
			FromSlice([]int{}).
			UnionBy(functional.FromSlice([]int{1, 2, 3}), func(i int) int { return i }).
			ToSlice()
		expected := []int{1, 2, 3}
		if !equalIgnoringOrder(expected, v, cmp.Compare[int]) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("right-empty", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3}).
			UnionBy(functional.FromSlice([]int{}), func(i int) int { return i }).
			ToSlice()
		expected := []int{1, 2, 3}
		if !equalIgnoringOrder(expected, v, cmp.Compare[int]) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("left, right non-empty", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3}).
			UnionBy(functional.FromSlice([]int{2, 3, 4}), func(i int) int { return i }).
			ToSlice()
		expected := []int{1, 2, 3, 4}
		if !equalIgnoringOrder(expected, v, cmp.Compare[int]) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("preserves order and removes duplicates", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3}).
			UnionBy(
				functional.FromSlice([]int{2, 4, 5, 1}),
				func(value int) int { return value },
			).
			ToSlice()
		expected := []int{1, 2, 3, 4, 5}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_ExceptBy(t *testing.T) {
	t.Run("left-empty", func(t *testing.T) {
		v := functional.
			FromSlice([]int{}).
			ExceptBy(functional.FromSlice([]int{1, 2, 3}), func(i int) int { return i }).
			ToSlice()
		expected := []int{}
		if !equalIgnoringOrder(expected, v, cmp.Compare[int]) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("right-empty", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3}).
			ExceptBy(functional.FromSlice([]int{}), func(i int) int { return i }).
			ToSlice()
		expected := []int{1, 2, 3}
		if !equalIgnoringOrder(expected, v, cmp.Compare[int]) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("left, right non-empty", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3}).
			ExceptBy(functional.FromSlice([]int{2, 3, 4}), func(i int) int { return i }).
			ToSlice()
		expected := []int{1}
		if !equalIgnoringOrder(expected, v, cmp.Compare[int]) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("preserves unmatched duplicates", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 5}).
			ExceptBy(
				functional.FromSlice([]int{1}),
				func(value int) int { return value % 2 },
			).
			ToSlice()
		expected := []int{2, 4, 2}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

// Set operations are one class of operations that iterate over two sequences at once.
// This runs the risk of calling yield after it returns false if both sequences are iterated over in a single new sequence.
// This test guards against those kinds of regressions.
func TestFunctioner_SetOperations_EarlyTermination(t *testing.T) {
	t.Run("concat", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Concat(functional.FromSlice([]int{6, 7, 8, 9, 10})).
			First(-1)
		expected := 1
		if expected != v {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("intersect-by", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
			IntersectBy(functional.FromSlice([]int{6, 7, 8, 9, 10}), func(i int) int { return i }).
			First(-1)
		expected := 6
		if expected != v {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("union-by", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
			UnionBy(functional.FromSlice([]int{6, 7, 8, 9, 10}), func(i int) int { return i }).
			Take(7).
			ToSlice()
		expected := []int{1, 2, 3, 4, 5, 6, 7}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("except-by", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
			ExceptBy(functional.FromSlice([]int{6, 7, 8, 9, 10}), func(i int) int { return i }).
			First(-1)
		expected := 1
		if expected != v {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}
