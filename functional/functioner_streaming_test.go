package functional_test

import (
	"iter"
	"math/rand/v2"
	"testing"
	"time"

	"example.com/playground/functional"
	"example.com/playground/internal"
)

// Tests in this file confirm that terminal operations on infinite streams should still terminate
// where semantically appropriate.

func infiniteSeq() iter.Seq[int64] {
	return func(yield func(int64) bool) {
		for {
			if !yield(rand.Int64()) {
				return
			}
		}
	}
}

func TestFunctioner_Streaming_WhereAny(t *testing.T) {
	internal.TryUntil(t, 5*time.Second,
		func() any {
			return functional.FromSeq(infiniteSeq()).
				Where(func(i int64) bool { return true }).
				Any(func(i int64) bool { return true })
		},
	)
}

func TestFunctioner_Streaming_ConcatAny(t *testing.T) {
	t.Run("inf-inf", func(t *testing.T) {
		internal.TryUntil(t, 5*time.Second,
			func() any {
				return functional.FromSeq(infiniteSeq()).
					Concat(functional.FromSeq(infiniteSeq())).
					Any(func(i int64) bool { return true })
			},
		)
	})

	t.Run("finite-inf", func(t *testing.T) {
		internal.TryUntil(t, 5*time.Second,
			func() any {
				return functional.FromSlice([]int64{1, 2, 3, 4, 5}).
					Concat(functional.FromSeq(infiniteSeq())).
					Any(func(i int64) bool { return i > 6 })
			},
		)
	})

}

func TestFunctioner_Streaming_UnionAny(t *testing.T) {
	internal.TryUntil(t, 5*time.Second,
		func() any {
			return functional.FromSeq(infiniteSeq()).
				UnionBy(functional.FromSeq(infiniteSeq()), func(i int64) int64 { return i }).
				Any(func(i int64) bool { return true })
		},
	)
}

func TestFunctioner_Streaming_ZipAny(t *testing.T) {
	internal.TryUntil(t, 5*time.Second,
		func() any {
			return functional.FromSeq(infiniteSeq()).
				Zip(functional.FromSeq(infiniteSeq()), func(i1, i2 int64) int64 { return i1 + i2 }).
				Any(func(i int64) bool { return true })
		},
	)
}
