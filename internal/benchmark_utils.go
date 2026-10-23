package internal

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
)

const (
	N = 100_000
)

//goland:noinspection GoUnusedGlobalVariable
var (
	global any // Used to prevent the compiler from optimising away calls
)

func UniqueIntegers(n int64) []int64 {
	res := make([]int64, 0, n)
	for i := int64(1); i <= n; i++ {
		res = append(res, i)
	}
	rand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})
	return res
}

func Integers(n int64) []int64 {
	res := make([]int64, 0, n)
	for i := int64(1); i <= n; i++ {
		res = append(res, rand.Int63n(n))
	}
	return res
}

func Shuffled[E any](n int, generator func() E) []E {
	res := make([]E, 0, n)
	for i := 1; i <= n; i++ {
		res = append(res, generator())
	}
	return res
}

type Benchmarker[E any] struct {
	Name string
	F    func([]E) any
}

func Compare(b *testing.B, benchmarkers ...Benchmarker[int64]) {
	b.Helper()
	input := Integers(N)
	inputs := make([][]int64, len(benchmarkers))
	for i := range inputs {
		inputs[i] = make([]int64, len(input))
		copy(inputs[i], input)
	}
	b.ResetTimer()

	for i, benchmarker := range benchmarkers {
		b.Run(fmt.Sprintf("%s @ %d", benchmarker.Name, N), func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				var result any
				for pb.Next() {
					result = benchmarker.F(inputs[i])
				}
				runtime.KeepAlive(result)

			})
		})
	}
}

func CompareUsing[E any](b *testing.B, input []E, benchmarkers ...Benchmarker[E]) {
	b.Helper()
	inputs := make([][]E, len(benchmarkers))
	for i := range inputs {
		inputs[i] = make([]E, len(input))
		copy(inputs[i], input)
	}
	b.ResetTimer()

	for i, benchmarker := range benchmarkers {
		b.Run(fmt.Sprintf("%s @ %d", benchmarker.Name, len(input)), func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				var result any
				for pb.Next() {
					result = benchmarker.F(inputs[i])
				}
				runtime.KeepAlive(result)
			})
		})
	}
}
