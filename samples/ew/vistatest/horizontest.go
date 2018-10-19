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
	"log"
	"github.com/shnifer/nigiri/samples/ew/vista"
	"github.com/shnifer/nigiri/samples/ew/vista/vistautils"
)

var Q *nigiri.Queue
var C MyCam
var L nigiri.Line
var ViewSector *vistautils.ViewSectorDrawer

var Lights []*Light
var SolidObjects []*SolidObject
var Clouds []*Cloud
var HorizonObjects []vista.Object

func mainLoopUpdate(dt float64){
	C.Update(dt)
	for _, v := range SolidObjects {
		v.Update(dt)
	}
	speed:=0.0
	if ebiten.IsKeyPressed(ebiten.Key1) {
		speed+= 30
	} else if ebiten.IsKeyPressed(ebiten.Key2){
		speed-= 30
	}
	for _, l:=range Lights{
		l.SetPosition(l.Center.Rotate(speed*dt))
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
	for _, v := range Clouds {
		Q.Add(v)
	}
	for _,light:=range Lights {
		light.Calculate(HorizonObjects)
		hRes := light.Result()
		for _, rec := range hRes {
			ViewSector.Point = light.Center
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

	L = nigiri.NewLine(cam.Phys(), 1)

	C = MyCam{cam}
	C.SetScale(0.4)


	for i:=0; i<20; i++{
		circle:=vista.Circle{Center: vec2.RandomInCircle(800), Radius: rand.Float64()*50+10}
		SolidObjects = append(SolidObjects, NewSolidObject(circle))
	}

	for i:=0; i<0; i++{
		circle:=vista.Circle{Center: vec2.RandomInCircle(800), Radius: rand.Float64()*50+10}
		Clouds = append(Clouds, NewCloud(circle, 1))
	}

	HorizonObjects = make([]vista.Object, 0)
	for i := 0; i < len(SolidObjects); i++ {
		HorizonObjects = append(HorizonObjects, SolidObjects[i])
	}
	for i := 0; i < len(Clouds); i++ {
		HorizonObjects = append(HorizonObjects, Clouds[i])
	}

	colors:=[...]color.Color{colornames.Red, colornames.Orange, colornames.Yellow, colornames.Green,
	colornames.Cyan, colornames.Blue, colornames.Purple}
	ViewSector = vistautils.NewViewSectorDrawer(-1,C)
	lightCount := len(colors)
	lightCount = 7
	for i:=0;i<lightCount;i++ {
		light := NewLight()
		light.Color = colors[i]
		light.SetPosition(vec2.InDir(float64(i)*360/float64(lightCount)).Mul(900))
		Lights = append(Lights, light)
	}
	ebiten.SetVsyncEnabled(false)
	err:=nigiri.Run(mainLoop, 800, 800, 1, "TEST")
	if err!=nil{
		log.Println("ERROR: ",err)
	}
}
