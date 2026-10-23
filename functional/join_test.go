package functional_test

import (
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_Join(t *testing.T) {
	identity := func(i int) int { return i }
	sameKey := func(i int) int { return 0 }

	t.Run("group-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		f2 := functional.FromSlice([]int{1, 1, 2, 2, 3, 3, 4, 4})
		v := f1.GroupJoin(f2, identity, identity, func(i1 int, i2s []int) []int {
			return append([]int{i1}, i2s...)
		}).ToSlice()

		expected := [][]int{
			{1, 1, 1},
			{2, 2, 2},
			{3, 3, 3},
			{4, 4, 4},
			{5},
		}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("left-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		f2 := functional.FromSlice([]int{3, 3, 4, 4})
		v := f1.LeftJoin(f2, identity, identity, func(i1 int, i2 *int) []int {
			if i2 != nil {
				return []int{i1, *i2}
			} else {
				return []int{i1}
			}
		}).ToSlice()

		expected := [][]int{
			{1},
			{2},
			{3, 3},
			{3, 3},
			{4, 4},
			{4, 4},
			{5},
		}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("left-join-with-empty", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		f2 := functional.FromSlice([]int{})
		v := f1.LeftJoin(f2, identity, identity, func(i1 int, i2 *int) []int {
			if i2 != nil {
				return []int{i1, *i2}
			} else {
				return []int{i1}
			}
		}).ToSlice()

		expected := [][]int{
			{1},
			{2},
			{3},
			{4},
			{5},
		}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("right-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{3, 4})
		f2 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		v := f1.RightJoin(f2, identity, identity, func(i1 *int, i2 int) []int {
			if i1 != nil {
				return []int{*i1, i2}
			} else {
				return []int{i2}
			}
		}).ToSlice()

		expected := [][]int{
			{1},
			{2},
			{3, 3},
			{4, 4},
			{5},
		}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("right-join-empty", func(t *testing.T) {
		f1 := functional.FromSlice([]int{})
		f2 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		v := f1.RightJoin(f2, identity, identity, func(i1 *int, i2 int) []int {
			if i1 != nil {
				return []int{*i1, i2}
			} else {
				return []int{i2}
			}
		}).ToSlice()

		expected := [][]int{
			{1},
			{2},
			{3},
			{4},
			{5},
		}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("left-right-inner", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		f2 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		v := f1.Join(f2, identity, identity, func(i1, i2 int) []int {
			return []int{i1, i2}
		}).ToSlice()

		expected := [][]int{
			{1, 1},
			{2, 2},
			{3, 3},
			{4, 4},
			{5, 5},
		}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("left-right-cross-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3})
		f2 := functional.FromSlice([]int{4, 5})
		v := f1.Join(f2, sameKey, sameKey, func(i1, i2 int) []int {
			return []int{i1, i2}
		}).ToSlice()

		expected := [][]int{
			{1, 4},
			{1, 5},
			{2, 4},
			{2, 5},
			{3, 4},
			{3, 5},
		}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("empty-right", func(t *testing.T) {
		f1 := functional.FromSlice([]int{})
		f2 := functional.FromSlice([]int{4, 5})
		v := f1.Join(f2, identity, identity, func(i1, i2 int) []int {
			return []int{i1, i2}
		}).ToSlice()

		expected := [][]int{}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("left-empty", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3})
		f2 := functional.FromSlice([]int{})
		v := f1.Join(f2, identity, identity, func(i1, i2 int) []int {
			return []int{i1, i2}
		}).ToSlice()

		expected := [][]int{}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("empty-empty", func(t *testing.T) {
		f1 := functional.FromSlice([]int{})
		f2 := functional.FromSlice([]int{})
		v := f1.Join(f2, identity, identity, func(i1, i2 int) []int {
			return []int{i1, i2}
		}).ToSlice()

		expected := [][]int{}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_TakeJoin(t *testing.T) {
	identity := func(i int) int { return i }
	t.Run("take-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		f2 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		v := f1.Take(1).Join(f2.Take(1), identity, identity, func(i1, i2 int) []int {
			return []int{i1, i2}
		}).ToSlice()

		expected := [][]int{
			{1, 1},
		}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
	t.Run("skip-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		f2 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		v := f1.Skip(1).Join(f2.Skip(4), identity, identity, func(i1, i2 int) []int {
			return []int{i1, i2}
		}).ToSlice()

		expected := [][]int{
			{5, 5},
		}
		if !equalIgnoringOrder(v, expected, func(s1, s2 []int) int {
			if slices.Equal(s1, s2) {
				return 0
			}
			return -1
		}) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_Join_EarlyTermination(t *testing.T) {
	identity := func(i int) int { return i }
	t.Run("join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		f2 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		v := f1.Join(f2, identity, identity, func(i1, i2 int) []int {
			return []int{i1, i2}
		}).Any(func(s []int) bool { return s[0] == 1 && s[1] == 1 })
		if !v {
			t.Fatalf("expected %v, got %v\n", true, v)
		}
	})

	t.Run("join-by", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		f2 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		v := f1.JoinBy(
			f2,
			func(i1, i2 int) bool {
				return i1 == i2
			},
			func(i1, i2 int) []int {
				return []int{i1, i2}
			}).
			Any(func(s []int) bool {
				return s[0] == 1 && s[1] == 1
			})
		if !v {
			t.Fatalf("expected %v, got %v\n", true, v)
		}
	})

	t.Run("left-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		f2 := functional.FromSlice([]int{4, 5})
		v := f1.LeftJoin(
			f2,
			identity,
			identity,
			func(i1 int, i2 *int) []int {
				if i2 == nil {
					return []int{i1}
				}
				return []int{i1, *i2}
			}).
			Any(func(s []int) bool {
				return s[0] == 1
			})
		if !v {
			t.Fatalf("expected %v, got %v\n", true, v)
		}
	})

	t.Run("right-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{4, 5})
		f2 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		v := f1.RightJoin(
			f2,
			identity,
			identity,
			func(i1 *int, i2 int) []int {
				if i1 == nil {
					return []int{i2}
				}
				return []int{*i1, i2}
			}).
			Any(func(s []int) bool {
				return s[0] == 1
			})
		if !v {
			t.Fatalf("expected %v, got %v\n", true, v)
		}
	})

	t.Run("group-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5})
		f2 := functional.FromSlice([]int{4, 5})
		v := f1.GroupJoin(
			f2,
			identity,
			identity,
			func(i1 int, i2s []int) []int {
				return append([]int{i1}, i2s...)
			}).
			Any(func(s []int) bool {
				return s[0] == 1
			})
		if !v {
			t.Fatalf("expected %v, got %v\n", true, v)
		}
	})
}
