package functional_test

import (
	"cmp"
	"maps"
	"slices"
	"testing"

	"example.com/playground/functional"
)

func TestFunctioner_All(t *testing.T) {
	v := functional.
		FromSlice([]int{1, 3, 5}).
		All(func(i int) bool { return i%2 == 1 })
	if !v {
		t.Fatalf("expected %t, got %t\n", !v, v)
	}

	t.Run("empty", func(t *testing.T) {
		v := functional.FromSlice([]int{}).All(func(int) bool { return false })
		if !v {
			t.Fatalf("expected %t, got %t\n", true, v)
		}
	})
}

func TestFunctioner_Any(t *testing.T) {
	v := functional.
		FromSlice([]int{1, 3, 5}).
		Any(func(i int) bool { return i%2 == 0 })
	if v {
		t.Fatalf("expected %t, got %t\n", !v, v)
	}

	t.Run("empty", func(t *testing.T) {
		v := functional.FromSlice([]int{}).Any(func(int) bool { return true })
		if v {
			t.Fatalf("expected %t, got %t\n", false, v)
		}
	})
}

func TestFunctioner_Aggregate(t *testing.T) {
	v := functional.
		FromSlice([]int{1, 2, 3, 4, 5}).
		Aggregate(0, func(v, next int) int {
			return v + next
		})
	expected := 15
	if v != expected {
		t.Fatalf("expected %d, got %d\n", expected, v)
	}

	t.Run("empty returns seed", func(t *testing.T) {
		v := functional.
			FromSlice([]int{}).
			Aggregate(42, func(accumulator, next int) int {
				return accumulator + next
			})
		if v != 42 {
			t.Fatalf("expected %d, got %d\n", 42, v)
		}
	})

	t.Run("changes accumulator type", func(t *testing.T) {
		v := functional.
			FromSlice([]string{"apple", "mango", "orange"}).
			Aggregate(0, func(total int, value string) int {
				return total + len(value)
			})
		if v != 16 {
			t.Fatalf("expected %d, got %d\n", 16, v)
		}
	})
}

func TestFunctioner_ToSequence(t *testing.T) {
	expected := []int{1, 2, 3, 4}
	v := slices.Collect(functional.FromSeq(slices.Values(expected)).ToSequence())
	if !slices.Equal(v, expected) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_ToSlice(t *testing.T) {
	expected := []int{1, 2, 3, 4}
	v := functional.FromSeq(slices.Values(expected)).ToSlice()
	if !slices.Equal(v, expected) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}

	t.Run("returns independent copy", func(t *testing.T) {
		input := []int{1, 2, 3}
		v := functional.FromSlice(input).ToSlice()
		v[0] = 99

		if input[0] != 1 {
			t.Fatalf("result aliases input: input=%v result=%v", input, v)
		}
	})
}

func TestFunctioner_ToMap(t *testing.T) {
	v := functional.
		FromSlice([]int{1, 2, 3}).
		ToMap(func(i int) int {
			return i
		}, func(i int) int {
			return i * i
		})
	expected := map[int]int{
		1: 1,
		2: 4,
		3: 9,
	}
	if !maps.Equal(v, expected) {
		t.Fatalf("expected %v, got %v\n", expected, v)
	}
}

func TestFunctioner_Chunk(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := functional.FromSlice([]int{}).Chunk[[]int](2).ToSlice()
		expected := [][]int{}
		if len(v) != len(expected) {
			t.Fatalf("expected %v, but got %v\n", expected, v)
		}
	})

	t.Run("uneven-chunk", func(t *testing.T) {
		const sz = 5
		input := make([]int, 0, sz)
		for i := range sz {
			input = append(input, i)
		}
		v := functional.FromSlice(input).Chunk[[]int](2).ToSlice()
		expected := [][]int{
			{0, 1},
			{2, 3},
			{4},
		}
		if len(v) != len(expected) {
			t.Fatalf("expected %v, but got %v\n", expected, v)
		}
		for i, elt := range expected {
			if !slices.Equal(elt, v[i]) {
				t.Fatalf("expected %v, but got %v\n", expected, v)
			}
		}
	})

	t.Run("even-chunk", func(t *testing.T) {
		const sz = 6
		input := make([]int, 0, sz)
		for i := range sz {
			input = append(input, i)
		}
		v := functional.FromSlice(input).Chunk[[]int](2).ToSlice()
		expected := [][]int{
			{0, 1},
			{2, 3},
			{4, 5},
		}
		if len(v) != len(expected) {
			t.Fatalf("expected %v, but got %v\n", expected, v)
		}
		for i, elt := range expected {
			if !slices.Equal(elt, v[i]) {
				t.Fatalf("expected %v, but got %v\n", expected, v)
			}
		}
	})
}

func TestFunctioner_Single(t *testing.T) {
	emptySeq := func(yield func(int) bool) {}
	oneSeq := func(yield func(int) bool) {
		yield(1)
	}
	twoSeq := func(yield func(int) bool) {
		for _, e := range []int{1, 2} {
			if !yield(e) {
				return
			}
		}
	}
	t.Run("empty", func(t *testing.T) {
		_, err := functional.FromSeq(emptySeq).Single()
		if err == nil {
			t.Fatalf("expected an error, but got no error\n")
		}
	})

	t.Run("one", func(t *testing.T) {
		v, err := functional.FromSeq(oneSeq).Single()
		if err != nil {
			t.Fatalf("expected no error, but got %v\n", err)
		}
		if v != 1 {
			t.Fatalf("expected 1, but got %d\n", v)
		}
	})

	t.Run("two", func(t *testing.T) {
		_, err := functional.FromSeq(twoSeq).Single()
		if err == nil {
			t.Fatalf("expected an error, but got no error\n")
		}
	})

	t.Run("filtered-empty", func(t *testing.T) {
		_, err := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool { return i > 10 }).
			Single()
		if err == nil {
			t.Fatalf("expected an error, but got no error\n")
		}
	})

	t.Run("filtered-one", func(t *testing.T) {
		v, err := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool { return i > 4 }).
			Single()
		if err != nil {
			t.Fatalf("expected no error, but got %v\n", err)
		}
		if v != 5 {
			t.Fatalf("expected 5, but got %d\n", v)
		}
	})

	t.Run("filtered-two", func(t *testing.T) {
		_, err := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool { return i > 3 }).
			Single()
		if err == nil {
			t.Fatalf("expected an error, but got no error\n")
		}
	})

	t.Run("ordered-empty", func(t *testing.T) {
		_, err := functional.
			FromSlice(shuffle([]int{})).
			OrderBy(cmp.Compare[int]).
			Single()
		if err == nil {
			t.Fatalf("expected an error, but got no error\n")
		}
	})

	t.Run("ordered-one", func(t *testing.T) {
		v, err := functional.
			FromSeq(oneSeq).
			OrderBy(cmp.Compare[int]).
			Single()
		if err != nil {
			t.Fatalf("expected no error, but got %v\n", err)
		}
		if v != 1 {
			t.Fatalf("expected 1, but got %d\n", v)
		}
	})

	t.Run("ordered-two", func(t *testing.T) {
		_, err := functional.
			FromSlice(shuffle([]int{1, 2})).
			OrderBy(cmp.Compare[int]).
			Single()
		if err == nil {
			t.Fatalf("expected an error, but got no error\n")
		}
	})

	t.Run("filtered-ordered-empty", func(t *testing.T) {
		_, err := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool { return i > 10 }).
			OrderBy(cmp.Compare[int]).
			Single()
		if err == nil {
			t.Fatalf("expected an error, but got no error\n")
		}
	})

	t.Run("filtered-ordered-one", func(t *testing.T) {
		v, err := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool { return i > 4 }).
			OrderBy(cmp.Compare[int]).
			Single()
		if err != nil {
			t.Fatalf("expected no error, but got %v\n", err)
		}
		if v != 5 {
			t.Fatalf("expected 5, but got %d\n", v)
		}
	})

	t.Run("filtered-ordered-two", func(t *testing.T) {
		_, err := functional.
			FromSlice([]int{1, 2, 3, 4, 5}).
			Where(func(i int) bool { return i > 3 }).
			OrderBy(cmp.Compare[int]).
			Single()
		if err == nil {
			t.Fatalf("expected an error, but got no error\n")
		}
	})
}
