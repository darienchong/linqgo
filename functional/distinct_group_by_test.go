package functional_test

import (
	"cmp"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_DistinctBy(t *testing.T) {
	v := functional.
		FromSlice([]int{1, 1, 1}).
		DistinctBy(func(i int) int { return i }).
		ToSlice()
	expected := []int{1}
	if !slices.Equal(expected, v) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}

	t.Run("preserves first occurrence", func(t *testing.T) {
		type user struct {
			id   int
			name string
		}

		v := functional.
			FromSlice([]user{{1, "Foo"}, {2, "Bar"}, {3, "Foo"}}).
			OrderByDescending(func(u1, u2 user) int { return cmp.Compare(u1.id, u2.id) }).
			DistinctBy(func(value user) string { return value.name }).
			ToSlice()
		expected := []user{{3, "Foo"}, {2, "Bar"}}

		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_GroupBy(t *testing.T) {
	v := functional.
		FromSlice([]int{1, 2, 3, 4, 5, 6}).
		GroupBy[bool, []int](func(i int) bool { return i%2 == 0 }).
		ToMap(
			func(p functional.Pair[bool, []int]) bool {
				return p.First
			},
			func(p functional.Pair[bool, []int]) []int {
				return p.Second
			},
		)
	expected := map[bool][]int{
		true:  {2, 4, 6},
		false: {1, 3, 5},
	}
	for k, val := range v {
		expectedVal := expected[k]
		if !slices.Equal(expectedVal, val) {
			t.Fatalf("%t: expected %v, got %v\n", k, expectedVal, val)
		}
	}

	t.Run("preserves first key order", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}).
			GroupBy[int, []int](func(value int) int { return value % 2 }).
			ToSlice()
		expected := []functional.Pair[int, []int]{
			{First: 1, Second: []int{1, 3, 5, 7, 9}},
			{First: 0, Second: []int{2, 4, 6, 8}},
		}

		if !slices.EqualFunc(v, expected, func(left, right functional.Pair[int, []int]) bool {
			return left.First == right.First && slices.Equal(left.Second, right.Second)
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_OrderGroupBy_ShouldGroupPreservingOrder(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	shuffle(input)

	inputCopy := make([]int, len(input))
	copy(inputCopy, input)
	v := functional.
		FromSlice(inputCopy).
		OrderBy(cmp.Compare[int]).
		GroupBy[int, []int](func(i int) int { return i % 2 }).
		ToMap(
			func(p functional.Pair[int, []int]) int {
				return p.First
			},
			func(p functional.Pair[int, []int]) []int {
				return p.Second
			},
		)
	expected := map[int][]int{
		0: {2, 4, 6},
		1: {1, 3, 5},
	}
	for k, val := range v {
		expectedVal := expected[k]
		if !slices.Equal(expectedVal, val) {
			t.Fatalf("%d: expected %v, got %v\n", k, expectedVal, val)
		}
	}
}
