package functional_test

import (
	"cmp"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_OrderByDescending(t *testing.T) {
	v := functional.
		FromSlice(shuffle([]int{1, 2, 3, 4})).
		OrderByDescending(cmp.Compare[int]).
		ToSlice()
	expected := []int{4, 3, 2, 1}
	if !slices.Equal(expected, v) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_OrderByKeyDescending(t *testing.T) {
	v := functional.
		FromSlice(shuffle([]int{1, 2, 4, 3})).
		OrderByKeyDescending(func(a int) int { return a }).
		ToSlice()
	expected := []int{4, 3, 2, 1}
	if !slices.Equal(expected, v) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_MultipleOrderBy(t *testing.T) {
	v := functional.
		FromSlice(shuffle([]functional.Pair[int, int]{
			{1, 2},
			{2, 2},
			{2, 1},
			{1, 1},
		})).
		OrderByKey(func(p functional.Pair[int, int]) int { return p.Second }).
		OrderByKey(func(p functional.Pair[int, int]) int { return p.First }).
		ToSlice()
	expected := []functional.Pair[int, int]{
		{2, 2},
		{2, 1},
		{1, 2},
		{1, 1},
	}
	slices.SortStableFunc(expected, func(a, b functional.Pair[int, int]) int { return cmp.Compare(a.Second, b.Second) })
	slices.SortStableFunc(expected, func(a, b functional.Pair[int, int]) int { return cmp.Compare(a.First, b.First) })
	if !slices.Equal(expected, v) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_OrderByThenBy(t *testing.T) {
	v := functional.
		FromSlice(shuffle([]functional.Pair[int, int]{
			{2, 2},
			{2, 1},
			{1, 2},
			{1, 1},
		})).
		OrderByKey(func(p functional.Pair[int, int]) int { return p.First }).
		ThenByKey(func(p functional.Pair[int, int]) int { return p.Second }).
		ToSlice()
	expected := []functional.Pair[int, int]{
		{2, 2},
		{2, 1},
		{1, 2},
		{1, 1},
	}
	slices.SortStableFunc(expected, func(a, b functional.Pair[int, int]) int { return cmp.Compare(a.Second, b.Second) })
	slices.SortStableFunc(expected, func(a, b functional.Pair[int, int]) int { return cmp.Compare(a.First, b.First) })
	if !slices.Equal(expected, v) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_MultipleThenBy(t *testing.T) {
	type tuple3 struct {
		A, B, C int
	}
	v := functional.
		FromSlice(shuffle([]tuple3{
			{0, 0, 0},
			{1, 1, 1},
			{0, 1, 1},
			{1, 1, 0},
			{0, 1, 0},
			{1, 0, 1},
			{0, 0, 1},
			{1, 0, 0},
		})).
		OrderByKey(func(t tuple3) int { return t.A }).
		ThenByKey(func(t tuple3) int { return t.B }).
		ThenByKey(func(t tuple3) int { return t.C }).
		ToSlice()
	expected := []tuple3{
		{0, 0, 0},
		{1, 1, 1},
		{0, 1, 1},
		{1, 1, 0},
		{0, 1, 0},
		{1, 0, 1},
		{0, 0, 1},
		{1, 0, 0},
	}
	slices.SortStableFunc(expected, func(a, b tuple3) int { return cmp.Compare(a.C, b.C) })
	slices.SortStableFunc(expected, func(a, b tuple3) int { return cmp.Compare(a.B, b.B) })
	slices.SortStableFunc(expected, func(a, b tuple3) int { return cmp.Compare(a.A, b.A) })
	if !slices.Equal(expected, v) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_Unordered(t *testing.T) {
	type tuple3 struct {
		A, B, C int
	}
	v := functional.
		FromSlice([]tuple3{
			{0, 0, 0},
			{1, 1, 1},
			{0, 1, 1},
			{1, 1, 0},
			{0, 1, 0},
			{1, 0, 1},
			{0, 0, 1},
			{1, 0, 0},
		}).
		OrderByKey(func(t tuple3) int { return t.A }).
		ThenByKey(func(t tuple3) int { return t.B }).
		ThenByKey(func(t tuple3) int { return t.C }).
		Unordered().
		ToSlice()
	expected := []tuple3{
		{0, 0, 0},
		{1, 1, 1},
		{0, 1, 1},
		{1, 1, 0},
		{0, 1, 0},
		{1, 0, 1},
		{0, 0, 1},
		{1, 0, 0},
	}
	if !slices.Equal(expected, v) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_BranchingOrderByShouldNotClobber(t *testing.T) {
	// Note we must do 3 OrderBys since this will result in an internal comparers slice of len=3 and cap=4
	f := functional.FromSlice(shuffle([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})).
		OrderBy(cmp.Compare[int]).
		OrderBy(cmp.Compare[int]).
		OrderBy(cmp.Compare[int])
	// Test if the ordering will be clobbered
	w1 := f.OrderBy(cmp.Compare[int])
	w2 := f.OrderByDescending(cmp.Compare[int])
	v1 := w1.ToSlice()
	expected1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if !slices.Equal(expected1, v1) {
		t.Fatalf("expected %v, but got %v\n", expected1, v1)
	}
	v2 := w2.ToSlice()
	expected2 := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	if !slices.Equal(expected2, v2) {
		t.Fatalf("expected %v, but got %v\n", expected2, v2)
	}
}
