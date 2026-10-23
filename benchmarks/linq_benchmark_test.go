package benchmarks

import (
	"cmp"
	"iter"
	"math/rand"
	"slices"
	"testing"

	"example.com/playground/functional"
	"example.com/playground/internal"
)

const (
	benchmarkCount    = 100
	benchmarkJoin     = 100
	benchmarkSkip     = 1_000
	benchmarkSeed     = 42
	orderedWindowSize = 100
	groupCount        = 50
)

type fatValue struct {
	value0 int
	value1 int64
	value2 int64
	value3 int64
	value4 int64
	value5 int64
	value6 int64
	value7 int64
}

type sortKey struct {
	first  int
	second int
	third  int
}

func newFatValue(value int) fatValue {
	return fatValue{
		value0: value,
		value1: int64(value),
		value2: int64(value),
		value3: int64(value),
		value4: int64(value),
		value5: int64(value),
		value6: int64(value),
		value7: int64(value),
	}
}

func (value fatValue) isEven() bool {
	return value.value0&1 == 0
}

func addFat(left, right fatValue) fatValue {
	return newFatValue(left.value0 + right.value0)
}

func multiplyFat(value fatValue, multiplier int) fatValue {
	return newFatValue(value.value0 * multiplier)
}

func randomInts(count int) []int {
	random := rand.New(rand.NewSource(benchmarkSeed))
	values := make([]int, count)
	for index := range values {
		values[index] = random.Intn(count)
	}
	return values
}

func randomFatValues(count int) []fatValue {
	values := randomInts(count)
	result := make([]fatValue, len(values))
	for index, value := range values {
		result[index] = newFatValue(value)
	}
	return result
}

func randomSortKeys(count int) []sortKey {
	random := rand.New(rand.NewSource(benchmarkSeed))
	values := make([]sortKey, count)
	for index := range values {
		values[index] = sortKey{
			first:  random.Intn(count),
			second: random.Intn(count),
			third:  random.Intn(count),
		}
	}
	return values
}

func compareSortKey(left, right sortKey) int {
	if comparison := cmp.Compare(left.first, right.first); comparison != 0 {
		return comparison
	}
	if comparison := cmp.Compare(left.second, right.second); comparison != 0 {
		return comparison
	}
	return cmp.Compare(left.third, right.third)
}

func repeatInts(duplicates int) []int {
	values := make([]int, 0, duplicates*benchmarkCount)
	for range duplicates {
		for value := range benchmarkCount {
			values = append(values, value)
		}
	}
	return values
}

func repeatFatValues(duplicates int) []fatValue {
	values := repeatInts(duplicates)
	result := make([]fatValue, len(values))
	for index, value := range values {
		result[index] = newFatValue(value)
	}
	return result
}

func seq[E any](values []E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, value := range values {
			if !yield(value) {
				return
			}
		}
	}
}

func sumSeq(values iter.Seq[int]) int {
	sum := 0
	for value := range values {
		sum += value
	}
	return sum
}

func BenchmarkSliceInt32Contains(b *testing.B) {
	values := randomInts(benchmarkCount)
	target := int(^uint(0) >> 1)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				for _, value := range values {
					if value == target {
						return true
					}
				}
				return false
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).Contains(target, equalInt)
			},
		},
	)
}

func BenchmarkSliceInt32Distinct(b *testing.B) {
	values := repeatInts(4)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				seen := make(map[int]struct{}, benchmarkCount)
				sum := 0
				for _, value := range values {
					if _, exists := seen[value]; exists {
						continue
					}
					seen[value] = struct{}{}
					sum += value
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).DistinctBy(identityInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSliceInt32Select(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				for _, value := range values {
					sum += value * 3
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).Select(tripleInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSliceInt32SkipTakeSelect(b *testing.B) {
	values := randomInts(benchmarkSkip + benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				for index := benchmarkSkip; index < benchmarkSkip+benchmarkCount; index++ {
					sum += values[index] * 3
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).Skip(benchmarkSkip).Take(benchmarkCount).Select(tripleInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSliceInt32SkipTakeWhere(b *testing.B) {
	values := randomInts(benchmarkSkip + benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				for index := benchmarkSkip; index < benchmarkSkip+benchmarkCount; index++ {
					if isEvenInt(values[index]) {
						sum += values[index]
					}
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).Skip(benchmarkSkip).Take(benchmarkCount).Where(isEvenInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSliceInt32Sum(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				for _, value := range values {
					sum += value
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSliceInt32Where(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				for _, value := range values {
					if isEvenInt(value) {
						sum += value
					}
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).Where(isEvenInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSliceInt32WhereCount(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				count := 0
				for _, value := range values {
					if isEvenInt(value) {
						count++
					}
				}
				return count
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).Where(isEvenInt).Count()
			},
		},
	)
}

func BenchmarkSliceInt32WhereSelect(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				for _, value := range values {
					if isEvenInt(value) {
						sum += value * 3
					}
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).Where(isEvenInt).Select(tripleInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSliceInt32WhereSelectToSlice(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				list := make([]int, 0)
				for _, value := range values {
					if isEvenInt(value) {
						list = append(list, value*3)
					}
				}
				return slices.Clone(list)
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).Where(isEvenInt).Select(tripleInt).ToSlice()
			},
		},
	)
}

func BenchmarkSeqInt32Contains(b *testing.B) {
	values := randomInts(benchmarkCount)
	target := int(^uint(0) >> 1)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				for value := range seq(values) {
					if value == target {
						return true
					}
				}
				return false
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSeq(seq(values)).Contains(target, equalInt)
			},
		},
	)
}

func BenchmarkSeqInt32Distinct(b *testing.B) {
	values := repeatInts(4)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				seen := make(map[int]struct{}, benchmarkCount)
				sum := 0
				for value := range seq(values) {
					if _, exists := seen[value]; exists {
						continue
					}
					seen[value] = struct{}{}
					sum += value
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSeq(seq(values)).DistinctBy(identityInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSeqInt32Select(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				for value := range seq(values) {
					sum += value * 3
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSeq(seq(values)).Select(tripleInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSeqInt32SkipTakeSelect(b *testing.B) {
	values := randomInts(benchmarkSkip + benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				index := 0
				for value := range seq(values) {
					if index >= benchmarkSkip && index < benchmarkSkip+benchmarkCount {
						sum += value * 3
					}
					index++
					if index == benchmarkSkip+benchmarkCount {
						break
					}
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSeq(seq(values)).Skip(benchmarkSkip).Take(benchmarkCount).Select(tripleInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSeqInt32SkipTakeWhere(b *testing.B) {
	values := randomInts(benchmarkSkip + benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				index := 0
				for value := range seq(values) {
					if index >= benchmarkSkip && index < benchmarkSkip+benchmarkCount && isEvenInt(value) {
						sum += value
					}
					index++
					if index == benchmarkSkip+benchmarkCount {
						break
					}
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSeq(seq(values)).Skip(benchmarkSkip).Take(benchmarkCount).Where(isEvenInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSeqInt32Sum(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				return sumSeq(seq(values))
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSeq(seq(values)).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSeqInt32Where(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				for value := range seq(values) {
					if isEvenInt(value) {
						sum += value
					}
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSeq(seq(values)).Where(isEvenInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSeqInt32WhereCount(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				count := 0
				for value := range seq(values) {
					if isEvenInt(value) {
						count++
					}
				}
				return count
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSeq(seq(values)).Where(isEvenInt).Count()
			},
		},
	)
}

func BenchmarkSeqInt32WhereSelect(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				sum := 0
				for value := range seq(values) {
					if isEvenInt(value) {
						sum += value * 3
					}
				}
				return sum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSeq(seq(values)).Where(isEvenInt).Select(tripleInt).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSeqInt32WhereSelectToSlice(b *testing.B) {
	values := randomInts(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				list := make([]int, 0)
				for value := range seq(values) {
					if isEvenInt(value) {
						list = append(list, value*3)
					}
				}
				return slices.Clone(list)
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSeq(seq(values)).Where(isEvenInt).Select(tripleInt).ToSlice()
			},
		},
	)
}

func BenchmarkSliceValueTypeContains(b *testing.B) {
	values := randomFatValues(benchmarkCount)
	target := newFatValue(int(^uint(0) >> 1))
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[fatValue]{
			Name: "baseline",
			F: func(values []fatValue) any {
				for _, value := range values {
					if value == target {
						return true
					}
				}
				return false
			},
		},
		internal.Benchmarker[fatValue]{
			Name: "functional",
			F: func(values []fatValue) any {
				return functional.FromSlice(values).Contains(target, equalFat)
			},
		},
	)
}

func BenchmarkSliceValueTypeDistinct(b *testing.B) {
	values := repeatFatValues(4)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[fatValue]{
			Name: "baseline",
			F: func(values []fatValue) any {
				seen := make(map[fatValue]struct{}, benchmarkCount)
				sum := fatValue{}
				for _, value := range values {
					if _, exists := seen[value]; exists {
						continue
					}
					seen[value] = struct{}{}
					sum = addFat(sum, value)
				}
				return sum
			},
		},
		internal.Benchmarker[fatValue]{
			Name: "functional",
			F: func(values []fatValue) any {
				return functional.FromSlice(values).DistinctBy(identityFat).Aggregate(fatValue{}, addFat)
			},
		},
	)
}

func BenchmarkSliceValueTypeSelect(b *testing.B) {
	values := randomFatValues(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[fatValue]{
			Name: "baseline",
			F: func(values []fatValue) any {
				sum := fatValue{}
				for _, value := range values {
					sum = addFat(sum, multiplyFat(value, 3))
				}
				return sum
			},
		},
		internal.Benchmarker[fatValue]{
			Name: "functional",
			F: func(values []fatValue) any {
				return functional.FromSlice(values).Select(tripleFat).Aggregate(fatValue{}, addFat)
			},
		},
	)
}

func BenchmarkSliceValueTypeSelectSum(b *testing.B) {
	values := randomFatValues(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[fatValue]{
			Name: "baseline",
			F: func(values []fatValue) any {
				sum := 0
				for _, value := range values {
					sum += value.value0
				}
				return sum
			},
		},
		internal.Benchmarker[fatValue]{
			Name: "functional",
			F: func(values []fatValue) any {
				return functional.FromSlice(values).Select(fatValue0).Aggregate(0, addInt)
			},
		},
	)
}

func BenchmarkSliceValueTypeSkipTakeSelect(b *testing.B) {
	values := randomFatValues(benchmarkSkip + benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[fatValue]{
			Name: "baseline",
			F: func(values []fatValue) any {
				sum := fatValue{}
				for index := benchmarkSkip; index < benchmarkSkip+benchmarkCount; index++ {
					sum = addFat(sum, multiplyFat(values[index], 3))
				}
				return sum
			},
		},
		internal.Benchmarker[fatValue]{
			Name: "functional",
			F: func(values []fatValue) any {
				return functional.FromSlice(values).Skip(benchmarkSkip).Take(benchmarkCount).Select(tripleFat).Aggregate(fatValue{}, addFat)
			},
		},
	)
}

func BenchmarkSliceValueTypeSkipTakeWhere(b *testing.B) {
	values := randomFatValues(benchmarkSkip + benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[fatValue]{
			Name: "baseline",
			F: func(values []fatValue) any {
				sum := fatValue{}
				for index := benchmarkSkip; index < benchmarkSkip+benchmarkCount; index++ {
					if values[index].isEven() {
						sum = addFat(sum, values[index])
					}
				}
				return sum
			},
		},
		internal.Benchmarker[fatValue]{
			Name: "functional",
			F: func(values []fatValue) any {
				return functional.FromSlice(values).Skip(benchmarkSkip).Take(benchmarkCount).Where(isEvenFat).Aggregate(fatValue{}, addFat)
			},
		},
	)
}

func BenchmarkSliceValueTypeWhere(b *testing.B) {
	values := randomFatValues(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[fatValue]{
			Name: "baseline",
			F: func(values []fatValue) any {
				sum := fatValue{}
				for _, value := range values {
					if value.isEven() {
						sum = addFat(sum, value)
					}
				}
				return sum
			},
		},
		internal.Benchmarker[fatValue]{
			Name: "functional",
			F: func(values []fatValue) any {
				return functional.FromSlice(values).Where(isEvenFat).Aggregate(fatValue{}, addFat)
			},
		},
	)
}

func BenchmarkSliceValueTypeWhereCount(b *testing.B) {
	values := randomFatValues(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[fatValue]{
			Name: "baseline",
			F: func(values []fatValue) any {
				count := 0
				for _, value := range values {
					if value.isEven() {
						count++
					}
				}
				return count
			},
		},
		internal.Benchmarker[fatValue]{
			Name: "functional",
			F: func(values []fatValue) any {
				return functional.FromSlice(values).Where(isEvenFat).Count()
			},
		},
	)
}

func BenchmarkSliceValueTypeWhereSelect(b *testing.B) {
	values := randomFatValues(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[fatValue]{
			Name: "baseline",
			F: func(values []fatValue) any {
				sum := fatValue{}
				for _, value := range values {
					if value.isEven() {
						sum = addFat(sum, multiplyFat(value, 3))
					}
				}
				return sum
			},
		},
		internal.Benchmarker[fatValue]{
			Name: "functional",
			F: func(values []fatValue) any {
				return functional.FromSlice(values).Where(isEvenFat).Select(tripleFat).Aggregate(fatValue{}, addFat)
			},
		},
	)
}

func BenchmarkSliceValueTypeWhereSelectToSlice(b *testing.B) {
	values := randomFatValues(benchmarkCount)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[fatValue]{
			Name: "baseline",
			F: func(values []fatValue) any {
				list := make([]fatValue, 0)
				for _, value := range values {
					if value.isEven() {
						list = append(list, multiplyFat(value, 3))
					}
				}
				return slices.Clone(list)
			},
		},
		internal.Benchmarker[fatValue]{
			Name: "functional",
			F: func(values []fatValue) any {
				return functional.FromSlice(values).Where(isEvenFat).Select(tripleFat).ToSlice()
			},
		},
	)
}

func BenchmarkSliceOrderedFirst(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				maximum := values[0]
				for _, value := range values[1:] {
					maximum = max(maximum, value)
				}
				return maximum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					OrderByDescending(cmp.Compare[int]).
					First(-1)
			},
		},
	)
}

func BenchmarkSliceFirstBy(b *testing.B) {
	values := internal.UniqueIntegers(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int64]{
			Name: "functional-first",
			F: func(values []int64) any {
				return functional.
					FromSlice(values).
					FirstBy(func(i int64) bool { return i == values[0] }, -1)
			},
		},
		internal.Benchmarker[int64]{
			Name: "functional-mid",
			F: func(values []int64) any {
				return functional.
					FromSlice(values).
					FirstBy(func(i int64) bool { return i == values[len(values)>>1] }, -1)
			},
		},
		internal.Benchmarker[int64]{
			Name: "functional-last",
			F: func(values []int64) any {
				return functional.
					FromSlice(values).
					FirstBy(func(i int64) bool { return i == values[len(values)-1] }, -1)
			},
		},
	)
}

func BenchmarkSliceOrderedFirstBy(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				maximum := -1
				init := false
				for _, value := range values {
					if !init && value%2 == 0 {
						maximum = value
						init = true
						continue
					}
					if value%2 != 0 {
						continue
					}
					maximum = max(maximum, value)
				}
				return maximum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					OrderByDescending(cmp.Compare[int]).
					FirstBy(func(i int) bool { return i%2 == 0 }, -1)
			},
		},
	)
}

func BenchmarkSliceOrderedLast(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				maximum := values[0]
				for _, value := range values[1:] {
					maximum = max(maximum, value)
				}
				return maximum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					OrderBy(cmp.Compare[int]).
					Last(-1)
			},
		},
	)
}

func BenchmarkSliceOrderedLastBy(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				maximum := values[0]
				for _, value := range values[1:] {
					if value%2 == 1 {
						continue
					}
					maximum = max(maximum, value)
				}
				return maximum
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					OrderBy(cmp.Compare[int]).
					LastBy(func(i int) bool { return i%2 == 0 }, -1)
			},
		},
	)
}

func BenchmarkSliceChainedSelects(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline-2-selects",
			F: func(values []int) any {
				result := slices.Clone(values)
				for index := range result {
					result[index] *= 2
				}
				for index := range result {
					result[index] *= 3
				}
				return result
			},
		},
		internal.Benchmarker[int]{
			Name: "functional-2-selects",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					Select(func(value int) int { return value * 2 }).
					Select(func(value int) int { return value * 3 }).
					ToSlice()
			},
		},
		internal.Benchmarker[int]{
			Name: "baseline-3-selects",
			F: func(values []int) any {
				result := slices.Clone(values)
				for index := range result {
					result[index] *= 2
				}
				for index := range result {
					result[index] *= 3
				}
				for index := range result {
					result[index] *= 5
				}
				return result
			},
		},
		internal.Benchmarker[int]{
			Name: "functional-3-selects",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					Select(func(value int) int { return value * 2 }).
					Select(func(value int) int { return value * 3 }).
					Select(func(value int) int { return value * 5 }).
					ToSlice()
			},
		},
		internal.Benchmarker[int]{
			Name: "baseline-4-selects",
			F: func(values []int) any {
				result := slices.Clone(values)
				for index := range result {
					result[index] *= 2
				}
				for index := range result {
					result[index] *= 3
				}
				for index := range result {
					result[index] *= 5
				}
				for index := range result {
					result[index] *= 7
				}
				return result
			},
		},
		internal.Benchmarker[int]{
			Name: "functional-4-selects",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					Select(func(value int) int { return value * 2 }).
					Select(func(value int) int { return value * 3 }).
					Select(func(value int) int { return value * 5 }).
					Select(func(value int) int { return value * 7 }).
					ToSlice()
			},
		},
	)
}

func BenchmarkSliceChainedWheresSelect(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline-2-wheres",
			F: func(values []int) any {
				result := make([]int, 0, len(values))
				for _, value := range values {
					if value >= 0 && value < internal.N {
						result = append(result, value*2)
					}
				}
				return result
			},
		},
		internal.Benchmarker[int]{
			Name: "functional-2-wheres",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					Where(func(value int) bool { return value >= 0 }).
					Where(func(value int) bool { return value < internal.N }).
					Select(func(value int) int { return value * 2 }).
					ToSlice()
			},
		},
		internal.Benchmarker[int]{
			Name: "baseline-3-wheres",
			F: func(values []int) any {
				result := make([]int, 0, len(values))
				for _, value := range values {
					if value >= 0 && value < internal.N && value%2 == 0 {
						result = append(result, value*2)
					}
				}
				return result
			},
		},
		internal.Benchmarker[int]{
			Name: "functional-3-wheres",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					Where(func(value int) bool { return value >= 0 }).
					Where(func(value int) bool { return value < internal.N }).
					Where(isEvenInt).
					Select(func(value int) int { return value * 2 }).
					ToSlice()
			},
		},
	)
}

func BenchmarkSliceOrderBy(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				result := slices.Clone(values)
				slices.SortStableFunc(result, cmp.Compare[int])
				return result
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.FromSlice(values).OrderBy(cmp.Compare[int]).ToSlice()
			},
		},
	)
}

func BenchmarkSliceMultipleOrderBy(b *testing.B) {
	values := randomSortKeys(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[sortKey]{
			Name: "baseline",
			F: func(values []sortKey) any {
				result := slices.Clone(values)
				slices.SortStableFunc(result, compareSortKey)
				return result
			},
		},
		internal.Benchmarker[sortKey]{
			Name: "functional",
			F: func(values []sortKey) any {
				return functional.
					FromSlice(values).
					OrderBy(func(left, right sortKey) int { return cmp.Compare(left.third, right.third) }).
					OrderBy(func(left, right sortKey) int { return cmp.Compare(left.second, right.second) }).
					OrderBy(func(left, right sortKey) int { return cmp.Compare(left.first, right.first) }).
					ToSlice()
			},
		},
	)
}

func BenchmarkSliceSelectOrderByThenBySkip(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				result := make([]int, len(values))
				for index, value := range values {
					result[index] = value * value
				}
				slices.SortStableFunc(result, func(left, right int) int {
					if comparison := -cmp.Compare(left, right); comparison != 0 {
						return comparison
					}
					return -cmp.Compare(left%2, right%2)
				})
				return result[orderedWindowSize:]
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					Select(func(value int) int { return value * value }).
					OrderByDescending(cmp.Compare[int]).
					ThenByKeyDescending(func(value int) int { return value % 2 }).
					Skip(orderedWindowSize).
					ToSlice()
			},
		},
	)
}

func BenchmarkSliceOrderedTake(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				result := slices.Clone(values)
				slices.SortStableFunc(result, cmp.Compare[int])
				return result[:orderedWindowSize]
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					OrderBy(cmp.Compare[int]).
					Take(orderedWindowSize).
					ToSlice()
			},
		},
	)
}

func BenchmarkSliceOrderedSkip(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				result := slices.Clone(values)
				slices.SortStableFunc(result, cmp.Compare[int])
				return result[orderedWindowSize:]
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					OrderBy(cmp.Compare[int]).
					Skip(orderedWindowSize).
					ToSlice()
			},
		},
	)
}

func BenchmarkSliceWhereAny(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				for _, value := range values {
					if value > 500_000 && value > 600_000 {
						return true
					}
				}
				return false
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					Where(func(value int) bool { return value > 500_000 }).
					Any(func(value int) bool { return value > 600_000 })
			},
		},
	)
}

func BenchmarkSliceOrderedWhereAny(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				result := slices.Clone(values)
				slices.SortStableFunc(result, cmp.Compare[int])
				for _, value := range result {
					if value > 500_000 && value > 600_000 {
						return true
					}
				}
				return false
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					OrderBy(cmp.Compare[int]).
					Where(func(value int) bool { return value > 500_000 }).
					Any(func(value int) bool { return value > 600_000 })
			},
		},
	)
}

func BenchmarkSliceGroupBy(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				groups := make(map[int][]int, groupCount)
				for _, value := range values {
					key := value % groupCount
					groups[key] = append(groups[key], value)
				}
				return groups
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					GroupBy[int, []int](func(value int) int { return value % groupCount }).
					ToMap(
						func(pair functional.Pair[int, []int]) int { return pair.First },
						func(pair functional.Pair[int, []int]) []int { return pair.Second },
					)
			},
		},
	)
}

func BenchmarkSliceTakeSkip(b *testing.B) {
	values := randomInts(internal.N)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int]{
			Name: "baseline",
			F: func(values []int) any {
				return values[internal.N/3 : 2*internal.N/3]
			},
		},
		internal.Benchmarker[int]{
			Name: "functional",
			F: func(values []int) any {
				return functional.
					FromSlice(values).
					Skip(internal.N / 3).
					Take(internal.N / 3).
					ToSlice()
			},
		},
	)
}

func BenchmarkSliceJoin(b *testing.B) {
	values := internal.UniqueIntegers(benchmarkJoin)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int64]{
			Name: "baseline",
			F: func(values []int64) any {
				evens, odds := make([]int64, len(values)>>1), make([]int64, len(values)>>1)
				for _, e := range values {
					if e%2 == 0 {
						evens = append(evens, e)
					} else {
						odds = append(odds, e)
					}
				}

				res := make([]int64, len(evens)*len(evens)+len(odds)*len(odds))
				for i := 0; i < len(evens); i++ {
					for j := 0; j < len(evens); j++ {
						res = append(res, evens[i]+evens[j])
					}
				}
				for i := 0; i < len(odds); i++ {
					for j := 0; j < len(odds); j++ {
						res = append(res, odds[i]+odds[j])
					}
				}
				return res
			},
		},
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(values []int64) any {
				f1 := functional.FromSlice(values)
				f2 := functional.FromSlice(values)
				return f1.Join(f2,
					func(i int64) int64 { return i % 2 },
					func(i int64) int64 { return i % 2 },
					func(i1, i2 int64) int64 { return i1 + i2 },
				).ToSlice()
			},
		},
	)
}

func BenchmarkSliceGroupJoin(b *testing.B) {
	values := internal.UniqueIntegers(benchmarkJoin)
	internal.CompareUsing(
		b,
		values,
		internal.Benchmarker[int64]{
			Name: "functional",
			F: func(values []int64) any {
				f1 := functional.FromSlice(values)
				f2 := functional.FromSlice(values)
				return f1.GroupJoin(f2,
					func(i int64) int64 { return i % 2 },
					func(i int64) int64 { return i % 2 },
					func(i1 int64, i2 []int64) int64 {
						sum := i1
						for _, i := range i2 {
							sum += i
						}
						return sum
					},
				).ToSlice()
			},
		},
	)
}

func addInt(left, right int) int {
	return left + right
}

func tripleInt(value int) int {
	return value * 3
}

func identityInt(value int) int {
	return value
}

func equalInt(left, right int) bool {
	return left == right
}

func isEvenInt(value int) bool {
	return value&1 == 0
}

func tripleFat(value fatValue) fatValue {
	return multiplyFat(value, 3)
}

func identityFat(value fatValue) fatValue {
	return value
}

func equalFat(left, right fatValue) bool {
	return left == right
}

func isEvenFat(value fatValue) bool {
	return value.isEven()
}

func fatValue0(value fatValue) int {
	return value.value0
}
