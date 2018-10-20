package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
	"golang.org/x/image/colornames"
	"image"
	"log"
	"github.com/shnifer/prof"
	"github.com/gobuffalo/packr"
	"github.com/shnifer/nigiri/nigiriutil"
)

var C MyCam
var Q *nigiri.Queue
var ClipRect *nigiri.Sprite

var Updaters []nigiri.Updater
var Draws []nigiri.DrawRequester
var MainMouser nigiri.MouseRect

func AddObjects(objs ...interface{}) {
	for _, obj := range objs {
		if updater, ok := obj.(nigiri.Updater); ok {
			Updaters = append(Updaters, updater)
		}
		if drawR, ok := obj.(nigiri.DrawRequester); ok {
			Draws = append(Draws, drawR)
		}
		if onMouser, ok := obj.(nigiri.OnMouser); ok {
			MainMouser.AddChild(onMouser)
		}
	}
}

func mainLoop(win *ebiten.Image, dt float64) error {
	for _, v := range Updaters {
		v.Update(dt)
	}
	//	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft){
	MainMouser.OnMouse(ebiten.CursorPosition())
	//	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	Q.Clear()
	for _, v := range Draws {
		v.DrawReqs(Q)
	}
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v\nDraw calls: %v\nUse WASD to move camera\n"+
		"Q-E to rotate\nZ-X to scale", ebiten.CurrentFPS(), Q.Len()))
	return nil
}

func main() {
	prof.StartProfile("gos")
	defer prof.StopProfile("gos")

	Q = nigiri.NewQueue()
	if !nigiri.IsJS {
		nigiri.SetTexLoader(nigiri.FileTexLoader("samples"))
		nigiri.SetFaceLoader(nigiri.FileFaceLoader("samples"))
	} else {
		box := packr.NewBox("res/")
		nigiri.SetTexLoader(nigiriutil.PackrTexLoader(box, ""))
		nigiri.SetFaceLoader(nigiriutil.PackrFaceLoader(box, ""))
	}

	C = MyCam{nigiri.NewCamera()}
	C.SetCenter(vec2.V2{X: 500, Y: 350})
	C.SetScale(5)
	C.SetClipRect(image.Rect(200, 100, 800, 600))

	MainMouser = nigiri.NewClickRect(func(x, y int) bool {
		log.Println("main mouser self ", x, ", ", y)
		return true
	})
	MainMouser.CatchRect = C.ClipRect()

	tt, _ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	tt.Fill(colornames.Darkkhaki)
	ClipRect = nigiri.NewSprite(nigiri.NewTex(tt), -1, C.Local())
	ClipRect.Position = vec2.V(500, 350)
	ClipRect.Pivot = vec2.Center
	ClipRect.Scaler = nigiri.NewFixedScaler(600, 500)

	StartData()
	for _, dat := range Data {
		vo := NewVisualObject(dat, C)
		AddObjects(vo)
	}

	AddObjects(ClipRect, C)

	ebiten.SetVsyncEnabled(false)
	nigiri.Run(mainLoop, 1000, 700, 1, "TEST")
}
