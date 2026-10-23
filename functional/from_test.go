package functional_test

import (
	"cmp"
	"maps"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFromSlice(t *testing.T) {
	v := functional.FromSlice([]int{1, 2, 3}).ToSlice()
	expected := []int{1, 2, 3}
	if !slices.Equal(v, expected) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFromSeq(t *testing.T) {
	v := functional.FromSeq(slices.Values([]int{1, 2, 3})).ToSlice()
	expected := []int{1, 2, 3}
	if !slices.Equal(v, expected) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFromSeq2(t *testing.T) {
	v := functional.FromSeq2(maps.All(map[int]int{
		1: 1,
		2: 2,
		3: 3,
	})).ToSlice()
	expected := []functional.Pair[int, int]{
		{1, 1},
		{2, 2},
		{3, 3},
	}
	// Maps don't have an order, so we shouldn't expect the two to agree on order
	if !equalIgnoringOrder(expected, v, func(a, b functional.Pair[int, int]) int {
		if v := cmp.Compare(a.First, b.First); v != 0 {
			return v
		}
		return cmp.Compare(a.Second, b.Second)
	}) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFromMap(t *testing.T) {
	v := functional.FromMap(map[int]int{
		1: 1,
		2: 2,
		3: 3,
	}).ToSlice()
	expected := []functional.Pair[int, int]{
		{1, 1},
		{2, 2},
		{3, 3},
	}
	if !equalIgnoringOrder(expected, v, func(a, b functional.Pair[int, int]) int {
		if v := cmp.Compare(a.First, b.First); v != 0 {
			return v
		}
		return cmp.Compare(a.Second, b.Second)
	}) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}
