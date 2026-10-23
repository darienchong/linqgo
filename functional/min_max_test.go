package functional_test

import (
	"cmp"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_MinMaxBy(t *testing.T) {
	t.Run("min-by", func(t *testing.T) {
		v := functional.
			FromSlice([]string{"one", "fine", "day", "in", "wonderland"}).
			MinBy(func(s1, s2 string) int { return cmp.Compare[int](len(s1), len(s2)) })
		expected := "in"
		if expected != v {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("max-by", func(t *testing.T) {
		v := functional.
			FromSlice([]string{"one", "fine", "day", "in", "wonderland"}).
			MaxBy(func(s1, s2 string) int { return cmp.Compare[int](len(s1), len(s2)) })
		expected := "wonderland"
		if expected != v {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}
