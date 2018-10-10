package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
	"image"
	_ "image/png"
	"math/rand"
	"image/color"
	"golang.org/x/image/colornames"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
)

var Q *nigiri.Queue
var C MyCam
var L nigiri.Line
var ViewSector *ViewSectorDrawer

type SolidObject struct {
	DiffShadowBody
	*nigiri.Sprite
}

type Light struct{
	*LightEmitter
	*Horizon
	Color color.Color
}
func NewLight() *Light{
	k:=[EmiDirCount]float64{}
	for i:=range k{
		k[i] = 1
	}
	res := &Light{
		LightEmitter: NewLightEmitter(1, k, ""),
		Horizon: NewHorizon(),
	}
	return res
}

func (l *Light) SetPosition(pos vec2.V2){
	l.LightEmitter.Center = pos
	l.Horizon.point = pos
}

var Lights []*Light
var SolidObjects []*SolidObject
var HorizonObjects []HorizonObject

func NewSolidObject(circle Circle) *SolidObject {
	Sprite := nigiri.NewSprite(nigiri.CircleTex(), 0, C.Phys())
	Sprite.Pivot = vec2.Center
	Sprite.SetSmooth(true)
	Sprite.Position = circle.Center
	visualSize := circle.Radius * 2
	Sprite.Scaler = nigiri.NewFixedScaler(visualSize, visualSize)

	return &SolidObject{
		DiffShadowBody: DiffShadowBody{Circle: circle, Albedo: 1},
		Sprite:         Sprite,
	}
}

type ViewSectorDrawer struct {
	*nigiri.TriDrawer
	Color color.Color
	Point vec2.V2
	Target HorizonObjectPart
}

func alpha(clr color.Color, alpha float64) color.Color{
	r,g,b,a:=clr.RGBA()
	return color.RGBA64{
		R:uint16(float64(r)*alpha),
		G:uint16(float64(g)*alpha),
		B:uint16(float64(b)*alpha),
		A:uint16(float64(a)*alpha),
	}
}

func (d *ViewSectorDrawer) GetVerticesIndices() ([]ebiten.Vertex, []uint16) {
	v:=make([]ebiten.Vertex,0)
	i:=make([]uint16,0)
	v = append(v, nigiri.VertexColor(d.Point, alpha(color.White, 0.7)))
	for i:=0;i<3;i++{
		dir:= d.Target.Angles.MedPart(float64(i)/2)
		pt:= d.Point.Add(vec2.InDir(dir).Mul(d.Target.Dist))
		var clr color.Color
		clr = color.White
		if i==1 && d.Color!=nil{
			clr = d.Color
		}
		clr = alpha(clr, 0.4)
		v = append(v, nigiri.VertexColor(pt, clr))
	}
	for m := 1; m+1 < len(v); m++ {
		i = append(i, 0, uint16(m), uint16(m+1))
	}
	return v, i
}

func NewViewSectorDrawer(layer nigiri.Layer, vTransformer nigiri.VTransformer) *ViewSectorDrawer{
	res:=&ViewSectorDrawer{}
	res.TriDrawer = nigiri.NewTriDrawer(res, layer, vTransformer)
	res.TriDrawer.ChangeableSrc = true
	return res
}

func mainLoopUpdate(dt float64){
	C.Update(dt)
	for _, v := range SolidObjects {
		v.Update(dt)
	}
	for _, l:=range Lights{
		l.SetPosition(l.point.Rotate(10*dt))
	}
}

func mainLoop(win *ebiten.Image, dt float64) error {
	//if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
	//	x, y := ebiten.CursorPosition()
	//	point := C.UnApplyPoint(vec2.V(float64(x), float64(y)))
	//	horizon.SetPointZoneDist(point, vec2.FullAnglePeriod, 0)
	//}
	mainLoopUpdate(dt)

	if ebiten.IsRunningSlowly() {
		return nil
	}
	Q.Clear()
	for _, v := range SolidObjects {
		Q.Add(v)
	}
	for _,light:=range Lights {
		hRes := light.Calculate(HorizonObjects, nil, HorizonObjects, nil)
		for _, rec := range hRes {
			ViewSector.Point = light.point
			ViewSector.Target = rec.Target
			ViewSector.Color = light.Color
			Q.Add(ViewSector)
		}
	}
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v",ebiten.CurrentFPS()))
	return nil
}

func main() {
	nigiri.StartProfile("ew")
	defer nigiri.StopProfile("ew")

	rand.Seed(3)
	Q = nigiri.NewQueue()
	cam := nigiri.NewCamera()
	cam.SetCenter(vec2.V2{X: 400, Y: 400})
	cam.SetClipRect(image.Rect(0, 0, 800, 800))
	nigiri.SetTexLoader(nigiri.FileTexLoader("samples"))
	//tex, err := nigiri.GetTex("planet_ani.png")
	//if err != nil {
	//	panic(err)
	//}
	//Ani, err = nigiri.NewFrameTexSrc(tex, 64, 64, 19,
	//	nigiri.AnimateFrameCycle(5))
	L = nigiri.NewLine(cam.Phys(), 1)

	C = MyCam{cam}
	C.SetScale(0.5)


	for i:=0; i<15; i++{
		circle:=Circle{Center: vec2.RandomInCircle(800), Radius: rand.Float64()*50+10}
		SolidObjects = append(SolidObjects, NewSolidObject(circle))
	}

	HorizonObjects = make([]HorizonObject, len(SolidObjects))
	for i := 0; i < len(SolidObjects); i++ {
		HorizonObjects[i] = SolidObjects[i]
	}

	colors:=[...]color.Color{colornames.Red, colornames.Orange, colornames.Yellow, colornames.Green,
	colornames.Cyan, colornames.Blue, colornames.Purple}
	ViewSector = NewViewSectorDrawer(-1,C)
	lightCount := len(colors)
	for i:=0;i<lightCount;i++ {
		light := NewLight()
		light.Color = colors[i]
		light.SetPosition(vec2.InDir(float64(i)*360/float64(lightCount)).Mul(900))
		Lights = append(Lights, light)
	}
	ebiten.SetVsyncEnabled(false)
	nigiri.Run(mainLoop, 800, 800, 1, "TEST")
}
