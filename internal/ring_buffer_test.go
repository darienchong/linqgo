package internal_test

import (
	"slices"
	"testing"

	"example.com/playground/internal"
)

func TestRingBuffer(t *testing.T) {
	t.Run("insertion sub-capacity", func(t *testing.T) {
		b := internal.NewRingBuffer[int](3)
		b.Push(0)
		b.Push(1)
		v := b.Collect()
		expected := []int{0, 1}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})

	t.Run("insertion over-capacity", func(t *testing.T) {
		b := internal.NewRingBuffer[int](2)
		b.Push(0)
		b.Push(1)
		b.Push(2)
		v := b.Collect()
		expected := []int{1, 2}
		if !slices.Equal(expected, v) {
			t.Fatalf("expected %v, got %v\n", expected, v)
		}
	})
}
