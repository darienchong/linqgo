package functional_test

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_Skip(t *testing.T) {
	t.Run("skip < length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			Skip(3).
			ToSlice()
		expected := []int{4}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("skip == length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			Skip(4).
			ToSlice()
		expected := []int{}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("skip > length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			Skip(1e9).
			ToSlice()
		expected := []int{}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_SkipWhile(t *testing.T) {
	t.Run("non-monotonic skip", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			SkipWhile(func(i int) bool { return i%2 == 1 }).
			ToSlice()
		expected := []int{2, 3, 4, 5}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("skip < length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			SkipWhile(func(i int) bool { return i < 3 }).
			ToSlice()
		expected := []int{3, 4}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("skip == length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			SkipWhile(func(i int) bool { return i < 5 }).
			ToSlice()
		expected := []int{}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("skip > length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			SkipWhile(func(i int) bool { return true }).
			ToSlice()
		expected := []int{}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_Take(t *testing.T) {
	t.Run("take from empty", func(t *testing.T) {
		v := functional.FromSlice([]int{}).
			Take(1).
			ToSlice()
		var expected []int
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %d, but got %d\n", expected, v)
		}
	})

	t.Run("take == 0", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			Take(0).
			ToSlice()
		expected := []int{}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("take < length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			Take(3).
			ToSlice()
		expected := []int{1, 2, 3}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("take == length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			Take(4).
			ToSlice()
		expected := []int{1, 2, 3, 4}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("take > length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			Take(1e9).
			ToSlice()
		expected := []int{1, 2, 3, 4}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	for _, tc := range []struct {
		count  int
		want   []int
		pulled int
	}{
		{count: 0, want: nil, pulled: 0},
		{count: 1, want: []int{1}, pulled: 1},
		{count: 3, want: []int{1, 2, 3}, pulled: 3},
		{count: 4, want: []int{1, 2, 3, 4}, pulled: 4},
		{count: 9, want: []int{1, 2, 3, 4}, pulled: 4},
	} {
		t.Run(fmt.Sprintf("pulls exactly %d", tc.count), func(t *testing.T) {
			pulled := 0
			source := func(yield func(int) bool) {
				for _, value := range []int{1, 2, 3, 4} {
					pulled++
					if !yield(value) {
						return
					}
				}
			}

			v := functional.FromSeq(iter.Seq[int](source)).Take(tc.count).ToSlice()
			if !slices.Equal(v, tc.want) {
				t.Fatalf("expected %v, got %v\n", tc.want, v)
			}
			if pulled != tc.pulled {
				t.Fatalf("expected %d pulls, got %d\n", tc.pulled, pulled)
			}
		})
	}
}

func TestFunctioner_TakeWhile(t *testing.T) {
	t.Run("take < length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			TakeWhile(func(i int) bool { return i < 3 }).
			ToSlice()
		expected := []int{1, 2}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("take == length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			TakeWhile(func(i int) bool { return i < 5 }).
			ToSlice()
		expected := []int{1, 2, 3, 4}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("take > length", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4}).
			TakeWhile(func(i int) bool { return true }).
			ToSlice()
		expected := []int{1, 2, 3, 4}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("ordered take", func(t *testing.T) {
		v := functional.
			FromSlice(shuffle([]int{1, 2, 3, 4})).
			OrderBy(cmp.Compare[int]).
			Take(2).
			ToSlice()
		expected := []int{1, 2}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}
