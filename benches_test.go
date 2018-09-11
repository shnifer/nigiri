package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri/vec2"
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
	r := NewRect(100, 200, vec2.V(0.3, 0.5))
	r.Position = vec2.V(10, 20)
	r.Angle = 40
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
	C.SetCenter(vec2.V2{X: 300, Y: 300})

	S := SpriteTrans{
		Pivot: vec2.TopMid,
	}
	I := NewDrawer(tex, 0, Transforms{S, C.Phys()})
	I.Layer = 1
	I.ChangeableTex = true
	I.SetSmooth(true)

	SetFaceLoader(FileFaceLoader("samples"))

	face, err := GetFace("furore.ttf", 20)
	if err != nil {
		panic(err)
	}
	TD := NewTextDrawer(face, 2)
	TD.Position = vec2.V(100, 100)
	TD.Color = colornames.Brown
	TD.Text = "just simple textdrawer\nsecond line"

	TS := NewTextSrc(1.2, true)
	TS.AddTextExt("text source sample\nmulti-line", face, 0, colornames.White)
	TS.AddTextExt("center or", face, 1, colornames.White)
	TS.AddTextExt("right aligned", face, 2, colornames.White)

	S2 := SpriteTrans{
		Pivot: vec2.Center,
	}
	I2 := NewDrawer(TS, 0, Transforms{S2, C.Phys()})

	dest, _ := ebiten.NewImage(1000, 1000, ebiten.FilterDefault)

	//log.Printf("SpriteTrans object Size = %v\nImageDrawer Size = %v\nTextSrc Size = %v\n",
	//	unsafe.Sizeof(S), unsafe.Sizeof(*I), unsafe.Sizeof(*TS))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Q.Clear()
		for j := 0; j < 10; j++ {
			S.Position.DoAddMul(vec2.V(1, 1), 1)
			Q.Add(I)
		}
		for j := 0; j < 10; j++ {
			S.Angle += 1
			Q.Add(I2)
		}
		//	Q.Add(TD)
		Q.Run(dest)
	}
}

//run me
func BenchmarkImageDrawer_calcDrawTag(b *testing.B) {
	id := NewDrawer(nil, 0)
	id.SetColor(colornames.Aliceblue)
	id.SetAlpha(0.3)
	id.SetSmooth(true)
	id.SetCompositeMode(ebiten.CompositeModeDestinationOut)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id.calcDrawTag()
	}
}

func BenchmarkCircleTex1(b *testing.B) {
	const radius = 512
	d := radius*2 + 1

	p := make([]byte, d*d*4)
	dw := d * 4
	r2 := radius * radius

	var ix, iy int
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for x := -radius; x <= radius; x++ {
			for y := -radius; y <= radius; y++ {
				if x*x+y*y <= r2 {
					ix = x + radius
					iy = y + radius
					for i := 0; i < 4; i++ {
						p[4*ix+dw*iy+i] = 255
					}
				}
			}
		}
	}
}

func BenchmarkCircleTex2(b *testing.B) {
	const radius = 512
	d := radius*2 + 1

	p := make([]byte, d*d*4)
	dw := d * 4
	r2 := radius * radius

	white := []byte{255, 255, 255, 255}

	var ix, iy int
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for x := -radius; x <= radius; x++ {
			for y := -radius; y <= radius; y++ {
				if x*x+y*y <= r2 {
					ix = x + radius
					iy = y + radius
					copy(p[4*ix+dw*iy:], white)
				}
			}
		}
	}
}
