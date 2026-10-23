package functional_test

import (
	"cmp"
	"math/rand"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_WhereSelect_ShouldCommute(t *testing.T) {
	t.Run("select-then-where", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Select(func(i int) int {
				return i * i
			}).
			Where(func(i int) bool {
				return i%2 == 0
			}).
			ToSlice()
		expected := []int{4, 16}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("where-then-select", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool {
				return i%2 == 0
			}).
			Select(func(i int) int {
				return i * i
			}).
			ToSlice()
		expected := []int{4, 16}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_OrderSelectThenLast_ShouldRespectOrder(t *testing.T) {
	v := functional.FromSlice([]string{"one", "fine", "day", "in", "wonderland"}).
		OrderByDescending(func(s1, s2 string) int { return cmp.Compare[int](len(s1), len(s2)) }).
		Select(func(s string) int { return len(s) }).
		Last(-1)
	expected := len("in")
	if expected != v {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_Stability(t *testing.T) {
	base := []int{1, 2, 3, 4, 5, 6, 7, 8}
	shuffle(base)

	t.Run("noop-sort-to-slice", func(t *testing.T) {
		baseCopy := make([]int, len(base))
		copy(baseCopy, base)
		v := functional.FromSlice(baseCopy).
			OrderBy(func(a, b int) int { return 0 }).
			ToSlice()
		shuffle(baseCopy)
		expected := base

		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, but got %v\n", expected, v)
		}
	})

	t.Run("noop-sort-take", func(t *testing.T) {
		baseCopy := make([]int, len(base))
		copy(baseCopy, base)
		v := functional.FromSlice(baseCopy).
			OrderBy(func(a, b int) int { return 0 }).
			Take(4).
			ToSlice()
		shuffle(baseCopy)
		expected := base[:4]

		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, but got %v\n", expected, v)
		}
	})

	t.Run("sort-to-slice", func(t *testing.T) {
		// Track original index
		input := []functional.Pair[int, int]{
			{0, 3},
			{1, 2},
			{2, 2},
			{3, 2},
			{4, 1},
		}
		v := functional.FromSlice(input).
			OrderBy(func(p1, p2 functional.Pair[int, int]) int { return cmp.Compare(p1.Second, p2.Second) }).
			ToSlice()
		expected := []functional.Pair[int, int]{
			{4, 1},
			{1, 2},
			{2, 2},
			{3, 2},
			{0, 3},
		}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, but got %v\n", expected, v)
		}
	})

	t.Run("sort-take", func(t *testing.T) {
		// Track original index
		input := []functional.Pair[int, int]{
			{0, 3},
			{1, 2},
			{2, 2},
			{3, 2},
			{4, 1},
		}
		v := functional.FromSlice(input).
			OrderBy(func(p1, p2 functional.Pair[int, int]) int { return cmp.Compare(p1.Second, p2.Second) }).
			Take(4).
			ToSlice()
		expected := []functional.Pair[int, int]{
			{4, 1},
			{1, 2},
			{2, 2},
			{3, 2},
		}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, but got %v\n", expected, v)
		}
	})
}

func TestFunctioner_NoopSort_ShouldBeStable(t *testing.T) {

}

func TestFunctioner_OrderByWhere_ShouldPreserveOrder(t *testing.T) {
	base := []int{1, 2, 3, 4, 5, 6, 7, 8}
	shuffle(base)
	v := functional.FromSlice(base).
		OrderBy(cmp.Compare[int]).
		Where(func(i int) bool { return i >= 5 }).
		ToSlice()
	shuffle(base)
	expected := []int{5, 6, 7, 8}
	if !slices.Equal(expected, v) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_OrderByTake_ShouldTakeRespectingOrder(t *testing.T) {
	base := []int{1, 2, 3, 4, 5, 6, 7, 8}
	shuffle(base)
	v := functional.FromSlice(base).
		OrderBy(cmp.Compare[int]).
		Take(4).
		ToSlice()
	shuffle(base)
	expected := []int{1, 2, 3, 4}
	if !slices.Equal(expected, v) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_ReiterationAfterEarlyExit(t *testing.T) {
	query := functional.
		FromSlice([]int{3, 1, 2, 2}).
		Where(func(value int) bool { return value > 1 }).
		OrderBy(cmp.Compare[int]).
		DistinctBy(func(value int) int { return value })

	for range query.ToSequence() {
		break
	}

	v := query.ToSlice()
	expected := []int{2, 3}
	if !slices.Equal(v, expected) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

// === test utils === //
func equalIgnoringOrder[E any](s1, s2 []E, cmp func(E, E) int) bool {
	if len(s1) == 0 {
		return len(s2) == 0
	}
	if len(s2) == 0 {
		return len(s1) == 0
	}
	if len(s1) != len(s2) {
		return false
	}
	for _, e := range s1 {
		found := false
		for _, e2 := range s2 {
			if cmp(e, e2) == 0 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func shuffle[E any](slice []E) []E {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
	return slice
}
