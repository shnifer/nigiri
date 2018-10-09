package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"math/rand"
	"time"
)

var Q *nigiri.Queue
var C MyCam
var Ani nigiri.TexSrcer
var L nigiri.Line

type SolidObject struct {
	DiffShadowBody
	*nigiri.Sprite
}

var SolidObjects []*SolidObject
var HorizonObjects []HorizonObject
var horizon *Horizon

func NewSolidObject(circle Circle) *SolidObject {
	Sprite := nigiri.NewSprite(Ani, 0, C.Phys())
	Sprite.Pivot = vec2.Center
	Sprite.SetSmooth(true)
	Sprite.Position = circle.Center
	visualSize := circle.Radius * 2 * 1.3
	Sprite.Scaler = nigiri.NewFixedScaler(visualSize, visualSize)

	return &SolidObject{
		DiffShadowBody: DiffShadowBody{Circle: circle, Albedo: 1},
		Sprite:         Sprite,
	}
}

func mainLoop(win *ebiten.Image, dt float64) error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		point := C.UnApplyPoint(vec2.V(float64(x), float64(y)))
		horizon.SetPointZoneDist(point, vec2.FullAnglePeriod, 0)
	}
	C.Update(dt)
	if ebiten.IsRunningSlowly() {
		return nil
	}
	for _, v := range SolidObjects {
		v.Update(dt)
	}

	hRes := horizon.Calculate(HorizonObjects, nil, HorizonObjects, nil)
	_ = hRes

	Q.Clear()
	for _, v := range SolidObjects {
		Q.Add(v)
	}
	for _, rec := range hRes {
		start, end := rec.Target.Angles.Get()
		L.From = horizon.point
		L.To = L.From.Add(vec2.InDir(start).Mul(rec.Target.Dist))
		L.SetColor(colornames.Red)
		Q.Add(L)
		L.From = horizon.point
		L.To = L.From.Add(vec2.InDir(end).Mul(rec.Target.Dist))
		L.SetColor(colornames.Royalblue)
		Q.Add(L)
	}
	Q.Run(win)
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	Q = nigiri.NewQueue()
	cam := nigiri.NewCamera()
	cam.SetCenter(vec2.V2{X: 400, Y: 400})
	cam.SetClipRect(image.Rect(0, 0, 800, 800))
	nigiri.SetTexLoader(nigiri.FileTexLoader("samples"))
	tex, err := nigiri.GetTex("planet_ani.png")
	if err != nil {
		panic(err)
	}
	Ani, err = nigiri.NewFrameTexSrc(tex, 64, 64, 19,
		nigiri.AnimateFrameCycle(5))
	L = nigiri.NewLine(cam.Phys(), 1)

	C = MyCam{cam}

	SolidObjects = append(SolidObjects, NewSolidObject(Circle{Center: vec2.V(0, 0), Radius: 50}))
	SolidObjects = append(SolidObjects, NewSolidObject(Circle{Center: vec2.V(200, 0), Radius: 50}))

	HorizonObjects = make([]HorizonObject, len(SolidObjects))
	for i := 0; i < len(SolidObjects); i++ {
		HorizonObjects[i] = SolidObjects[i]
	}
	horizon = NewHorizon()
	horizon.SetPointZoneDist(vec2.V(-250, 0), vec2.FullAnglePeriod, 0)
	nigiri.Run(mainLoop, 800, 800, 1, "TEST")
}
