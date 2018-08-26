package main

import (
	"testing"
	"math/rand"
	"sort"
	"log"
	"unsafe"
	"github.com/hajimehoshi/ebiten"
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

type doer interface {
	Do()
}
type doerUpdater interface {
	Updater
	doer
}

type x struct{}
func (x) Update(dt float64){}
func (x) Do(){}

type y struct{}
func (y) Do(){}

func Benchmark_EmptyUpdater(b *testing.B){
	var V doerUpdater = x{}
	for i:=0;i<b.N;i++{
		V.Update(0.1)
	}
}

func Benchmark_CheckUpdate(b *testing.B){
	var V doer = x{}
	for i:=0;i<b.N;i++{
		if u,ok:=V.(Updater);ok {
			u.Update(0.1)
		}
	}
}

func Benchmark_CheckNoUpdate(b *testing.B){
	var V doer = y{}
	o:=ebiten.DrawImageOptions{}
	for i:=0;i<b.N;i++{
		if u,ok:=V.(Updater);ok {
			u.Update(0.1)
		}
	}
}
