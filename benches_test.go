package nigiri

import (
	"testing"
	"math/rand"
	"sort"
	"github.com/hajimehoshi/ebiten"
	"github.com/Shnifer/nigiri/v2"
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
	for i:=0;i<b.N;i++{
		if u,ok:=V.(Updater);ok {
			u.Update(0.1)
		}
	}
}

func Benchmark_Geom_Scale(b *testing.B){
	G:=ebiten.GeoM{}
	for i:=0;i<b.N;i++{
		G.Scale(1.1, 1.2)
	}
}

func Benchmark_Geom_Rotate(b *testing.B){
	G:=ebiten.GeoM{}
	for i:=0;i<b.N;i++{
		G.Rotate(0.1)
	}
}

func Benchmark_Geom_Translate(b *testing.B){
	G:=ebiten.GeoM{}
	for i:=0;i<b.N;i++{
		G.Translate(0.1, 0.2)
	}
}

func Benchmark_Geom_Concat(b *testing.B){
	G:=ebiten.GeoM{}
	H:=ebiten.GeoM{}
	H.Translate(0.1, 0.2)
	H.Rotate(0.1)
	H.Scale(1.1, 1.2)
	for i:=0;i<b.N;i++{
		G.Concat(H)
	}
}

func Benchmark_Rect_Corners (b *testing.B) {
	r:=NewRect(100,200, v2.V2{0.3, 0.5})
	r.Pos = v2.V2{10,20}
	r.Ang = 40
	for i:=0;i<b.N;i++{
		_ = r.Corners()
	}
}
