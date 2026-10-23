package functional_test

import (
	"cmp"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_Select(t *testing.T) {
	t.Run("square", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Select(func(i int) int {
				return i * i
			}).
			ToSlice()
		expected := []int{1, 4, 9, 16, 25}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("int-to-str", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3}).
			Select(func(i int) string {
				switch i {
				case 1:
					return "one"
				case 2:
					return "two"
				case 3:
					return "three"
				default:
					return "?"
				}
			}).
			ToSlice()
		expected := []string{"one", "two", "three"}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_MultipleSelect(t *testing.T) {
	t.Run("multiple-select", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Select(func(i int) int {
				return i * 2
			}).
			Select(func(i int) int {
				return i * 3
			}).
			Select(func(i int) int {
				return i * 5
			}).
			ToSlice()
		expected := []int{1 * 30, 2 * 30, 3 * 30, 4 * 30, 5 * 30}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("multiple-select-should-respect-order", func(t *testing.T) {
		v := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			OrderByDescending(cmp.Compare[int]).
			Select(func(i int) int {
				return i * 2
			}).
			Select(func(i int) int {
				return i * 3
			}).
			Select(func(i int) int {
				return i * 5
			}).
			ToSlice()
		expected := []int{5 * 30, 4 * 30, 3 * 30, 2 * 30, 1 * 30}
		if !slices.Equal(v, expected) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}

func TestFunctioner_SelectMany(t *testing.T) {
	v := functional.
		FromSlice([]int{1, 2, 3}).
		SelectMany(func(i int) []int {
			return []int{
				i,
				i * i,
				i * i * i,
			}
		}).
		ToSlice()
	expected := []int{
		1, 1, 1,
		2, 4, 8,
		3, 9, 27,
	}
	if !equalIgnoringOrder(expected, v, cmp.Compare[int]) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}
