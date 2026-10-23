package functional

import (
	"cmp"
	"slices"
	"testing"
)

func TestFunctioner_boundedSequencer(t *testing.T) {
	t.Run("empty,ordered,reversed=false", func(t *testing.T) {
		v := slices.Collect(boundedSequencer(FromSlice([]int{}).OrderBy(cmp.Compare[int]), 2, false))
		expected := []int{}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("empty,unordered,reversed=false", func(t *testing.T) {
		v := slices.Collect(boundedSequencer(FromSlice([]int{}), 2, false))
		expected := []int{}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("unordered,reversed=false", func(t *testing.T) {
		v := slices.Collect(boundedSequencer(FromSlice([]int{2, 4, 3, 1}), 2, false))
		expected := []int{2, 4}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("unordered,reversed=true", func(t *testing.T) {
		v := slices.Collect(boundedSequencer(FromSlice([]int{2, 4, 3, 1}), 2, true))
		expected := []int{1, 3}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("ordered,reversed=false", func(t *testing.T) {
		v := slices.Collect(boundedSequencer(FromSlice([]int{2, 4, 3, 1}).
			OrderBy(cmp.Compare[int]), 2, false))
		expected := []int{1, 2}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("ordered,reversed=true", func(t *testing.T) {
		v := slices.Collect(boundedSequencer(FromSlice([]int{2, 4, 3, 1}).
			OrderBy(cmp.Compare[int]), 2, true))
		expected := []int{4, 3}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_singleSequencer(t *testing.T) {
	t.Run("empty,ordered,reversed=false", func(t *testing.T) {
		v := slices.Collect(singleSequencer(FromSlice([]int{}).OrderBy(cmp.Compare[int]), false))
		expected := []int{}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("empty,unordered,reversed=false", func(t *testing.T) {
		v := slices.Collect(singleSequencer(FromSlice([]int{}), false))
		expected := []int{}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("unordered,reversed=false", func(t *testing.T) {
		v := slices.Collect(singleSequencer(FromSlice([]int{2, 4, 3, 1}), false))
		expected := []int{2}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("unordered,reversed=true", func(t *testing.T) {
		v := slices.Collect(singleSequencer(FromSlice([]int{2, 4, 3, 1}), true))
		expected := []int{1}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("ordered,reversed=false", func(t *testing.T) {
		v := slices.Collect(singleSequencer(FromSlice([]int{2, 4, 3, 1}).
			OrderBy(cmp.Compare[int]), false))
		expected := []int{1}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("ordered,reversed=true", func(t *testing.T) {
		v := slices.Collect(singleSequencer(FromSlice([]int{2, 4, 3, 1}).
			OrderBy(cmp.Compare[int]), true))
		expected := []int{4}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}
