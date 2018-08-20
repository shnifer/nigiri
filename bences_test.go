package main

import (
	"testing"
	"math/rand"
	"sort"
)

func Benchmark_SortInt(b *testing.B) {
	x:=make([]int, 1000)
	b.ResetTimer()
	for n:=0;n<b.N;n++ {
		b.StopTimer()
		for i := 0; i < 1000; i++ {
			x[i] = rand.Intn(300)
		}
		b.StartTimer()
		sort.Ints(x)
	}
}

func Benchmark_Sortfloat(b *testing.B) {
	x:=make([]float64, 1000)
	b.ResetTimer()
	for n:=0;n<b.N;n++ {
		b.StopTimer()
		for i := 0; i < 1000; i++ {
			x[i] = rand.Float64()*300
		}
		b.StartTimer()
		sort.Float64s(x)
	}
}