package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
	"golang.org/x/image/colornames"
	"image"
)

var C MyCam
var Q *nigiri.Queue
var ClipRect nigiri.Sprite

var Updaters []nigiri.Updater
var Draws []nigiri.DrawRequester

func AddObjects(objs ...interface{}) {
	for _, obj := range objs {
		if updater, ok := obj.(nigiri.Updater); ok {
			Updaters = append(Updaters, updater)
		}
		if drawR, ok := obj.(nigiri.DrawRequester); ok {
			Draws = append(Draws, drawR)
		}
	}
}

func mainLoop(win *ebiten.Image, dt float64) error {
	for _, v := range Updaters {
		v.Update(dt)
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	Q.Clear()
	for _, v := range Draws {
		v.DrawReqs(Q)
	}
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v\nDraws: %v", ebiten.CurrentFPS(), Q.Len()))
	return nil
}

func main() {
	nigiri.StartProfile("gos")
	defer nigiri.StopProfile("gos")

	Q = nigiri.NewQueue()
	nigiri.SetTexLoader(nigiri.FileTexLoader("samples"))
	nigiri.SetFaceLoader(nigiri.FileFaceLoader("samples"))

	C = MyCam{nigiri.NewCamera()}
	C.SetCenter(vec2.V2{X: 500, Y: 500})
	C.SetScale(5)
	C.SetClipRect(image.Rect(200, 100, 800, 600))

	tt, _ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	tt.Fill(colornames.Darkkhaki)
	ClipRect = nigiri.NewSprite(nigiri.NewTex(tt), -1, C.Local())
	ClipRect.Position = vec2.V(500, 350)
	ClipRect.Pivot = vec2.Center
	ClipRect.Scaler = nigiri.NewFixedScaler(600, 500)

	LoadData()
	for _, dat := range Data {
		vo := NewVisualObject(dat, C)
		AddObjects(vo)
	}

	AddObjects(ClipRect, C)

	ebiten.SetVsyncEnabled(true)
	nigiri.Run(mainLoop, 1000, 700, 1, "TEST")
}