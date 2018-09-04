package nigiri

import (
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
	_ "image/png"
	"math/rand"
	"sort"
	"testing"
)

func Benchmark_SortInt(b *testing.B) {
	x := make([]int, 1000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		for i := 0; i < 1000; i++ {
			x[i] = rand.Intn(300)
		}
		b.StartTimer()
		sort.Ints(x)
	}
}

func Benchmark_Sortfloat(b *testing.B) {
	x := make([]float64, 1000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		for i := 0; i < 1000; i++ {
			x[i] = rand.Float64() * 300
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

func (x) Update(dt float64) {}
func (x) Do()               {}

type y struct{}

func (y) Do() {}

func Benchmark_EmptyUpdater(b *testing.B) {
	var V doerUpdater = x{}
	for i := 0; i < b.N; i++ {
		V.Update(0.1)
	}
}

func Benchmark_CheckUpdate(b *testing.B) {
	var V doer = x{}
	for i := 0; i < b.N; i++ {
		if u, ok := V.(Updater); ok {
			u.Update(0.1)
		}
	}
}

func Benchmark_CheckNoUpdate(b *testing.B) {
	var V doer = y{}
	for i := 0; i < b.N; i++ {
		if u, ok := V.(Updater); ok {
			u.Update(0.1)
		}
	}
}

func Benchmark_Geom_Scale(b *testing.B) {
	G := ebiten.GeoM{}
	for i := 0; i < b.N; i++ {
		G.Scale(1.1, 1.2)
	}
}

func Benchmark_Geom_Rotate(b *testing.B) {
	G := ebiten.GeoM{}
	for i := 0; i < b.N; i++ {
		G.Rotate(0.1)
	}
}

func Benchmark_Geom_Translate(b *testing.B) {
	G := ebiten.GeoM{}
	for i := 0; i < b.N; i++ {
		G.Translate(0.1, 0.2)
	}
}

func Benchmark_Geom_Concat(b *testing.B) {
	G := ebiten.GeoM{}
	H := ebiten.GeoM{}
	H.Translate(0.1, 0.2)
	H.Rotate(0.1)
	H.Scale(1.1, 1.2)
	for i := 0; i < b.N; i++ {
		G.Concat(H)
	}
}

func Benchmark_Rect_Corners(b *testing.B) {
	r := NewRect(100, 200, v2.V2{0.3, 0.5})
	r.Position = v2.V2{10, 20}
	r.Ang = 40
	for i := 0; i < b.N; i++ {
		_ = r.Corners()
	}
}

func BenchmarkSampleQueue(b *testing.B) {
	Q := NewQueue()
	SetTexLoader(FileTexLoader("samples"))
	tex, err := GetTex("HUD_Ship.png")
	if err != nil {
		panic(err)
	}
	C := NewCamera()
	C.SetCenter(v2.V2{X: 300, Y: 300})

	Opts := SpriteOpts{
		Src:          NewStatic(tex, nil, "ship"),
		Pivot:        v2.TopMid,
		Smooth:       true,
		CamTransform: C.Phys(),
		//		ChangeableTex: true,
	}
	S := NewSprite(Opts)

	SetFaceLoader(FileFaceLoader("samples"))

	face, err := GetFace("furore.ttf", 20)
	if err != nil {
		panic(err)
	}
	TD := NewTextDrawer(face, 2)
	TD.Position = v2.V(100, 100)
	TD.Color = colornames.Brown
	TD.Text = "just simple textdrawer\nsecond line"

	TS := NewTextSrc(1.2, 1, true)
	TS.AddText("text source sample\nmulti-line", face, 0, colornames.White)
	TS.AddText("center or", face, 1, colornames.White)
	TS.AddText("right aligned", face, 2, colornames.White)

	S2 := SpriteOpts{
		Src:          TS,
		Pivot:        v2.Center,
		Smooth:       true,
		CamTransform: C.Phys(),
		//		ChangeableTex: true,
	}.New()

	dest, _ := ebiten.NewImage(1000, 1000, ebiten.FilterDefault)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Q.Clear()
		for j := 0; j < 10; j++ {
			S.Position.DoAddMul(v2.V(1, 1), 1)
			Q.Add(S)
		}
		for j := 0; j < 10; j++ {
			S.Angle += 1
			Q.Add(S2)
		}
		//	Q.Add(TD)
		Q.Run(dest)
	}
}
