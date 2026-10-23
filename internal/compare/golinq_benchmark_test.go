package compare_test

import (
	"cmp"
	"math/rand"
	"slices"
	"testing"

	"example.com/playground/functional"
	"example.com/playground/internal"
	golinq "github.com/ahmetb/go-linq/v5"
)

const (
	takeCount   = 10
	skipCount   = internal.N / 2
	groupCount  = 50
	distinctMod = 1_000
)

func orderedInput() []int64 {
	input := make([]int64, internal.N)
	for i := range input {
		input[i] = int64(i)
	}
	return input
}

type positionCase struct {
	name   string
	target int64
}

func representativePositions() []positionCase {
	return []positionCase{
		{name: "first", target: 0},
		{name: "middle", target: internal.N / 2},
		{name: "last", target: internal.N - 1},
		{name: "missing", target: -1},
	}
}

type prefixCase struct {
	name  string
	count int
}

func representativePrefixes() []prefixCase {
	return []prefixCase{
		{name: "empty", count: 0},
		{name: "small", count: takeCount},
		{name: "half", count: internal.N / 2},
		{name: "all", count: internal.N},
	}
}

func BenchmarkFromSlice(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					ToSlice()

			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.FromSlice(input).
					ToSlice()
			},
		},
	)
}

func BenchmarkFromSeq(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSeq(slices.Values(input)).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSeq(slices.Values(input)).
					ToSlice()
			},
		},
	)
}

func BenchmarkFromMap(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				m := make(map[int64]int64, len(input))
				for i, e := range input {
					m[int64(i)] = e
				}
				return functional.
					FromMap(m).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				m := make(map[int64]int64, len(input))
				for i, e := range input {
					m[int64(i)] = e
				}
				return golinq.
					FromMap(m).
					ToSlice()
			},
		},
	)
}

func BenchmarkSelect(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					Select(func(i int64) int64 { return i * i }).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					Select(func(i int64) int64 { return i * i }).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "functional-multiple",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					Select(func(i int64) int64 { return i * 2 }).
					Select(func(i int64) int64 { return i * 3 }).
					Select(func(i int64) int64 { return i * 5 }).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq-multiple",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					Select(func(i int64) int64 { return i * 2 }).
					Select(func(i int64) int64 { return i * 3 }).
					Select(func(i int64) int64 { return i * 5 }).
					ToSlice()
			},
		},
	)
}

func BenchmarkSelectMany(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					SelectMany(func(i int64) []int64 { return []int64{i, i * i} }).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					SelectMany(func(i int64) golinq.Query[int64] {
						return golinq.FromSlice([]int64{i, i * i})
					}).
					ToSlice()
			},
		},
	)
}

func BenchmarkWhere(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					Where(func(i int64) bool { return i%2 == 0 }).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					Where(func(i int64) bool { return i%2 == 0 }).
					ToSlice()
			},
		},
	)
}

func BenchmarkOrderBy(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					OrderBy(cmp.Compare[int64]).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					Sort(cmp.Less[int64]).
					ToSlice()
			},
		},
	)
}

func BenchmarkMultipleOrderBy(b *testing.B) {
	randCmp := func(i, j int64) int {
		if rand.Float64() > 0.5 {
			return -1
		} else {
			return 1
		}
	}
	randLess := func(i, j int64) bool {
		return rand.Float64() > 0.5
	}

	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					OrderBy(randCmp).
					OrderBy(randCmp).
					OrderBy(randCmp).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					Sort(randLess).
					Sort(randLess).
					Sort(randLess).
					ToSlice()
			},
		},
	)
}

func BenchmarkOrderByKey(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					OrderByKey(func(i int64) int64 { return i }).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					OrderBy(func(i int64) int64 { return i }).
					ToSlice()
			},
		},
	)
}

func BenchmarkOrderByDescending(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					OrderByDescending(cmp.Compare[int64]).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					OrderByDescending(func(i int64) int64 { return i }).
					ToSlice()
			},
		},
	)
}

func BenchmarkOrderByKeyDescending(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					OrderByKeyDescending(func(i int64) int64 { return i }).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					OrderByDescending(func(i int64) int64 { return i }).
					ToSlice()
			},
		},
	)
}

func BenchmarkOrderByThenBy(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					OrderByKey(func(i int64) int64 { return i % groupCount }).
					ThenByKey(func(i int64) int64 { return i }).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					OrderBy(func(i int64) int64 { return i % groupCount }).
					ThenBy(func(i int64) int64 { return i }).
					ToSlice()
			},
		},
	)
}

func BenchmarkOrderByThenByDescending(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					OrderByKey(func(i int64) int64 { return i % groupCount }).
					ThenByKeyDescending(func(i int64) int64 { return i }).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					OrderBy(func(i int64) int64 { return i % groupCount }).
					ThenByDescending(func(i int64) int64 { return i }).
					ToSlice()
			},
		},
	)
}

func BenchmarkSkip(b *testing.B) {
	for _, tc := range []prefixCase{
		{name: "none", count: 0},
		{name: "small", count: takeCount},
		{name: "half", count: skipCount},
		{name: "all", count: internal.N},
	} {
		b.Run(tc.name, func(b *testing.B) {
			internal.Compare(b,
				internal.Benchmarker[int64]{
					Name: "functional",
					F: func(input []int64) any {
						return functional.
							FromSlice(input).
							Skip(tc.count).
							ToSlice()
					},
				},
				internal.Benchmarker[int64]{
					Name: "go-linq",
					F: func(input []int64) any {
						return golinq.
							FromSlice(input).
							Skip(tc.count).
							ToSlice()
					},
				},
			)
		})
	}
}

func BenchmarkSkipWhile(b *testing.B) {
	input := orderedInput()
	for _, tc := range representativePrefixes() {
		b.Run(tc.name, func(b *testing.B) {
			internal.CompareUsing(b, input,
				internal.Benchmarker[int64]{
					Name: "functional",
					F: func(input []int64) any {
						return functional.
							FromSlice(input).
							SkipWhile(func(i int64) bool { return i < int64(tc.count) }).
							ToSlice()
					},
				},
				internal.Benchmarker[int64]{
					Name: "go-linq",
					F: func(input []int64) any {
						return golinq.
							FromSlice(input).
							SkipWhile(func(i int64) bool { return i < int64(tc.count) }).
							ToSlice()
					},
				},
			)
		})
	}
}

func BenchmarkTake(b *testing.B) {
	for _, tc := range representativePrefixes() {
		b.Run(tc.name, func(b *testing.B) {
			internal.Compare(b,
				internal.Benchmarker[int64]{
					Name: "functional",
					F: func(input []int64) any {
						return functional.
							FromSlice(input).
							Take(tc.count).
							ToSlice()
					},
				},
				internal.Benchmarker[int64]{
					Name: "go-linq",
					F: func(input []int64) any {
						return golinq.
							FromSlice(input).
							Take(tc.count).
							ToSlice()
					},
				},
			)
		})
	}
}

func BenchmarkTakeWhile(b *testing.B) {
	input := orderedInput()
	for _, tc := range representativePrefixes() {
		b.Run(tc.name, func(b *testing.B) {
			internal.CompareUsing(b, input,
				internal.Benchmarker[int64]{
					Name: "functional",
					F: func(input []int64) any {
						return functional.
							FromSlice(input).
							TakeWhile(func(i int64) bool { return i < int64(tc.count) }).
							ToSlice()
					},
				},
				internal.Benchmarker[int64]{
					Name: "go-linq",
					F: func(input []int64) any {
						return golinq.
							FromSlice(input).
							TakeWhile(func(i int64) bool { return i < int64(tc.count) }).
							ToSlice()
					},
				},
			)
		})
	}
}

func BenchmarkFirst(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					First(-1)
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				value, _ := golinq.
					FromSlice(input).
					First()
				return value
			},
		},
	)
}

func BenchmarkFirstBy(b *testing.B) {
	input := orderedInput()
	for _, tc := range representativePositions() {
		b.Run(tc.name, func(b *testing.B) {
			internal.CompareUsing(b, input,
				internal.Benchmarker[int64]{
					Name: "functional",
					F: func(input []int64) any {
						return functional.
							FromSlice(input).
							FirstBy(func(i int64) bool { return i == tc.target }, -1)
					},
				},
				internal.Benchmarker[int64]{
					Name: "go-linq",
					F: func(input []int64) any {
						value, _ := golinq.
							FromSlice(input).
							FirstWith(func(i int64) bool { return i == tc.target })
						return value
					},
				},
			)
		})
	}
}

func BenchmarkLast(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					Last(-1)
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				value, _ := golinq.
					FromSlice(input).
					Last()
				return value
			},
		},
	)
}

func BenchmarkLastBy(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					LastBy(func(i int64) bool { return i%2 == 0 }, -1)
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				value, _ := golinq.
					FromSlice(input).
					LastWith(func(i int64) bool { return i%2 == 0 })
				return value
			},
		},
	)
}

func BenchmarkElementAt(b *testing.B) {
	input := orderedInput()
	positions := []struct {
		name string
		idx  int
	}{
		{name: "first", idx: 0},
		{name: "middle", idx: internal.N / 2},
		{name: "last", idx: internal.N - 1},
	}
	for _, tc := range positions {
		b.Run(tc.name, func(b *testing.B) {
			internal.CompareUsing(b, input,
				internal.Benchmarker[int64]{
					Name: "functional",
					F: func(input []int64) any {
						v, _ := functional.
							FromSlice(input).
							ElementAt(tc.idx)
						return v
					},
				},
				internal.Benchmarker[int64]{
					Name: "go-linq",
					F: func(input []int64) any {
						value, _ := golinq.
							FromSlice(input).
							Skip(tc.idx).
							First()
						return value
					},
				},
			)
		})
	}
}

func BenchmarkCount(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					Count()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					Count()
			},
		},
	)
}

func BenchmarkContains(b *testing.B) {
	input := orderedInput()
	for _, tc := range representativePositions() {
		b.Run(tc.name, func(b *testing.B) {
			internal.CompareUsing(b, input,
				internal.Benchmarker[int64]{
					Name: "functional",
					F: func(input []int64) any {
						return functional.
							FromSlice(input).
							Contains(tc.target, func(a, b int64) bool { return a == b })
					},
				},
				internal.Benchmarker[int64]{
					Name: "go-linq",
					F: func(input []int64) any {
						return golinq.
							FromSlice(input).
							Contains(tc.target)
					},
				},
			)
		})
	}
}

func BenchmarkMaxBy(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					MaxBy(cmp.Compare[int64])
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				value, _ := golinq.Max(golinq.FromSlice(input))
				return value
			},
		},
	)
}

func BenchmarkMinBy(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					MinBy(cmp.Compare[int64])
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				value, _ := golinq.Min(golinq.FromSlice(input))
				return value
			},
		},
	)
}

func BenchmarkAggregate(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					Aggregate(int64(0), func(acc, next int64) int64 { return acc + next })
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					AggregateWithSeed(int64(0), func(acc, next int64) int64 { return acc + next })
			},
		},
	)
}

func BenchmarkAll(b *testing.B) {
	input := orderedInput()
	for _, tc := range representativePositions() {
		b.Run(tc.name, func(b *testing.B) {
			internal.CompareUsing(b, input,
				internal.Benchmarker[int64]{
					Name: "functional",
					F: func(input []int64) any {
						return functional.
							FromSlice(input).
							All(func(i int64) bool { return i != tc.target })
					},
				},
				internal.Benchmarker[int64]{
					Name: "go-linq",
					F: func(input []int64) any {
						return golinq.
							FromSlice(input).
							All(func(i int64) bool { return i != tc.target })
					},
				},
			)
		})
	}
}

func BenchmarkAny(b *testing.B) {
	input := orderedInput()
	for _, tc := range representativePositions() {
		b.Run(tc.name, func(b *testing.B) {
			internal.CompareUsing(b, input,
				internal.Benchmarker[int64]{
					Name: "functional",
					F: func(input []int64) any {
						return functional.
							FromSlice(input).
							Any(func(i int64) bool { return i == tc.target })
					},
				},
				internal.Benchmarker[int64]{
					Name: "go-linq",
					F: func(input []int64) any {
						return golinq.
							FromSlice(input).
							AnyWith(func(i int64) bool { return i == tc.target })
					},
				},
			)
		})
	}
}

func BenchmarkDistinctBy(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					DistinctBy(func(i int64) int64 { return i % distinctMod }).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					DistinctBy(func(i int64) int64 { return i % distinctMod }).
					ToSlice()
			},
		},
	)
}

func BenchmarkGroupBy(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					GroupBy[int64, []int64](func(i int64) int64 { return i % groupCount }).
					ToMap(
						func(p functional.Pair[int64, []int64]) int64 { return p.First },
						func(p functional.Pair[int64, []int64]) []int64 { return p.Second },
					)
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					GroupBy(
						func(i int64) int64 { return i % groupCount },
						func(i int64) int64 { return i },
					).
					ToSlice()
			},
		},
	)
}

func BenchmarkConcat(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					Concat(functional.FromSlice(input)).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					Concat(golinq.FromSlice(input)).
					ToSlice()
			},
		},
	)
}

func BenchmarkToSlice(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					ToSlice()
			},
		},
	)
}

func BenchmarkToMap(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				return functional.
					FromSlice(input).
					ToMap(
						func(i int64) int64 { return i },
						func(i int64) int64 { return i * i },
					)
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				return golinq.
					FromSlice(input).
					ToMapBy(
						func(i int64) int64 { return i },
						func(i int64) int64 { return i * i },
					)
			},
		},
	)
}

func BenchmarkToSequence(b *testing.B) {
	internal.Compare(b,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				var sum int64
				for e := range functional.FromSlice(input).ToSequence() {
					sum += e
				}
				return sum
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				var sum int64
				for item := range golinq.FromSlice(input).Iterate {
					sum += item
				}
				return sum
			},
		},
	)
}

func BenchmarkJoin(b *testing.B) {
	const n = 100_000
	mod := int64(n / 10)
	input := internal.UniqueIntegers(n)
	internal.CompareUsing(b, input,
		internal.Benchmarker[int64]{
			Name: "functional/same-size",
			F: func(input []int64) any {
				f1 := functional.FromSlice(input)
				f2 := functional.FromSlice(input)
				return f1.Join(f2,
					func(i int64) int64 { return i % mod },
					func(i int64) int64 { return i % mod },
					func(i1, i2 int64) int64 { return i1 + i2 },
				).ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq/same-size",
			F: func(input []int64) any {
				f1 := golinq.FromSlice(input)
				f2 := golinq.FromSlice(input)

				return f1.Join(f2,
					func(i int64) int64 { return i % mod },
					func(i int64) int64 { return i % mod },
					func(i1, i2 int64) int64 { return i1 + i2 },
				).ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "functional/unbalanced-size-left",
			F: func(input []int64) any {
				f1 := functional.FromSlice(input[:len(input)>>3])
				f2 := functional.FromSlice(input)
				return f1.Join(f2,
					func(i int64) int64 { return i % mod },
					func(i int64) int64 { return i % mod },
					func(i1, i2 int64) int64 { return i1 + i2 },
				).ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq/unbalanced-size-left",
			F: func(input []int64) any {
				f1 := golinq.FromSlice(input[:len(input)>>3])
				f2 := golinq.FromSlice(input)

				return f1.Join(f2,
					func(i int64) int64 { return i % mod },
					func(i int64) int64 { return i % mod },
					func(i1, i2 int64) int64 { return i1 + i2 },
				).ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "functional/unbalanced-size-right",
			F: func(input []int64) any {
				f1 := functional.FromSlice(input)
				f2 := functional.FromSlice(input[:len(input)>>3])
				return f1.Join(f2,
					func(i int64) int64 { return i % mod },
					func(i int64) int64 { return i % mod },
					func(i1, i2 int64) int64 { return i1 + i2 },
				).ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq/unbalanced-size-right",
			F: func(input []int64) any {
				f1 := golinq.FromSlice(input)
				f2 := golinq.FromSlice(input[:len(input)>>3])

				return f1.Join(f2,
					func(i int64) int64 { return i % mod },
					func(i int64) int64 { return i % mod },
					func(i1, i2 int64) int64 { return i1 + i2 },
				).ToSlice()
			},
		},
	)
}

func BenchmarkGroupJoin(b *testing.B) {
	const n = 100_000
	mod := int64(n / 10)
	input := internal.UniqueIntegers(n)
	internal.CompareUsing(b, input,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(input []int64) any {
				f1 := functional.FromSlice(input)
				f2 := functional.FromSlice(input)
				return f1.GroupJoin(f2,
					func(i int64) int64 { return i % mod },
					func(i int64) int64 { return i % mod },
					func(i1 int64, i2 []int64) []int64 { return append(append(make([]int64, 0, 1+len(i2)), i1), i2...) },
				).ToSlice()
			},
		},
		internal.Benchmarker[int64]{
			Name: "go-linq",
			F: func(input []int64) any {
				f1 := golinq.FromSlice(input)
				f2 := golinq.FromSlice(input)

				return f1.GroupJoin(f2,
					func(i int64) int64 { return i % mod },
					func(i int64) int64 { return i % mod },
					func(i1 int64, i2 []int64) []int64 { return append(append(make([]int64, 0, 1+len(i2)), i1), i2...) },
				).ToSlice()
			},
		},
	)
}
