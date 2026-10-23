package functional_test

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_Where(t *testing.T) {
	t.Run("to-slice", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool {
				return i%2 == 0
			}).
			ToSlice()
		expected := []int{2, 4}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("to-map", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool {
				return i%2 == 0
			}).
			ToMap(
				func(i int) int { return i },
				func(i int) int { return i },
			)
		expected := map[int]int{
			2: 2,
			4: 4,
		}
		if !maps.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

}

func TestFunctioner_MultipleWhere(t *testing.T) {
	v := functional.
		FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
		Where(func(i int) bool { return i != 1 }).
		Where(func(i int) bool { return i != 2 }).
		Where(func(i int) bool { return i != 3 }).
		ToSlice()
	expected := []int{4, 5, 6, 7, 8, 9, 10}
	if !slices.Equal(v, expected) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_MultipleWhereSelect(t *testing.T) {
	t.Run("multi-where-select", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
			Where(func(i int) bool { return i != 1 }).
			Where(func(i int) bool { return i != 2 }).
			Where(func(i int) bool { return i != 3 }).
			Select(func(i int) int { return i * 2 }).
			ToSlice()
		expected := []int{8, 10, 12, 14, 16, 18, 20}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("multi-where-selectmany", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
			Where(func(i int) bool { return i != 1 }).
			Where(func(i int) bool { return i != 2 }).
			Where(func(i int) bool { return i != 3 }).
			SelectMany(func(i int) []int { return []int{i * 2, i * 3} }).
			ToSlice()
		expected := []int{
			8, 12,
			10, 15,
			12, 18,
			14, 21,
			16, 24,
			18, 27,
			20, 30,
		}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_MultipleWhereConcatSelect(t *testing.T) {
	t.Run("multi-where-select", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool { return i != 1 }).
			Where(func(i int) bool { return i != 2 }).
			Where(func(i int) bool { return i != 3 }).
			Concat(functional.FromSlice([]int{6, 7, 8, 9, 10})).
			Select(func(i int) int { return i * 2 }).
			ToSlice()
		expected := []int{8, 10, 12, 14, 16, 18, 20}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_BranchingWhereShouldNotClobber(t *testing.T) {
	f := functional.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
		Where(func(i int) bool { return i != 1 }).
		Where(func(i int) bool { return i != 2 }).
		Where(func(i int) bool { return i != 3 })
	w1 := f.Where(func(i int) bool { return i != 4 })
	w2 := f.Where(func(i int) bool { return i != 5 })
	v1 := w1.ToSlice()
	expected1 := []int{5, 6, 7, 8, 9, 10}
	if !slices.Equal(expected1, v1) {
		t.Fatalf("expected %v, but got %v\n", expected1, v1)
	}
	v2 := w2.ToSlice()
	expected2 := []int{4, 6, 7, 8, 9, 10}
	if !slices.Equal(expected2, v2) {
		t.Fatalf("expected %v, but got %v\n", expected2, v2)
	}
}

func TestFunctioner_JoinShouldRespectOuterOrdering(t *testing.T) {
	identity := func(i int) int { return i }

	t.Run("inner-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5}).OrderByDescending(cmp.Compare[int])
		f2 := functional.FromSlice([]int{2, 3, 4})
		v := f1.Join(f2, identity, identity, func(i1, i2 int) int { return i1*10 + i2 }).ToSlice()
		expected := []int{44, 33, 22}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("left-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{1, 2, 3, 4, 5}).OrderByDescending(cmp.Compare[int])
		f2 := functional.FromSlice([]int{2, 3, 4})
		v := f1.LeftJoin(f2, identity, identity, func(i1 int, i2 *int) int {
			if i2 == nil {
				return i1 * 10
			}
			return i1*10 + *i2
		}).ToSlice()
		expected := []int{50, 44, 33, 22, 10}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("right-join", func(t *testing.T) {
		f1 := functional.FromSlice([]int{2, 3, 4})
		f2 := functional.FromSlice([]int{1, 2, 3, 4, 5}).OrderByDescending(cmp.Compare[int])
		v := f1.RightJoin(f2, identity, identity, func(i1 *int, i2 int) int {
			if i1 == nil {
				return i2 * 10
			}
			return i2*10 + *i1
		}).ToSlice()
		expected := []int{50, 44, 33, 22, 10}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("group-join", func(t *testing.T) {
		null := func(r rune) int { return 0 }
		f1 := functional.FromSlice([]rune("abc")).OrderByDescending(cmp.Compare[rune])
		f2 := functional.FromSlice([]rune("abcde")).OrderByDescending(cmp.Compare[rune])
		v := f1.GroupJoin(f2, null, null, func(r1 rune, r2 []rune) string {
			return fmt.Sprintf("%s-%s", string(r1), string(r2))
		}).ToSlice()
		expected := []string{
			"c-edcba",
			"b-edcba",
			"a-edcba",
		}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}
