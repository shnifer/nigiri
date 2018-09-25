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

var Q *nigiri.Queue
var C *nigiri.Camera
var Sector *nigiri.Sector
var ClipRect nigiri.Sprite

func mainLoop(win *ebiten.Image, dt float64) error {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		C.Translate(vec2.InDir(C.Angle()).Rotate90().Mul(1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		C.Translate(vec2.InDir(C.Angle()).Rotate90().Mul(-1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		C.Translate(vec2.InDir(C.Angle()).Mul(1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		C.Translate(vec2.InDir(C.Angle()).Mul(-1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		C.Rotate(1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		C.Rotate(-1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		C.MulScale(1.05)
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		C.MulScale(1 / 1.05)
	}

	if ebiten.IsKeyPressed(ebiten.Key1) {
		Sector.StartAng += 1
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		Sector.StartAng -= 1
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		Sector.EndAng += 1
	}
	if ebiten.IsKeyPressed(ebiten.Key4) {
		Sector.EndAng -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		Sector.Radius *= 1.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		Sector.Radius /= 1.05
	}

	Q.Clear()
	Q.Add(ClipRect)
	Q.Add(Sector)
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v\nDraws: %v", ebiten.CurrentFPS(), Q.Len()))
	return nil
}

func main() {
	nigiri.StartProfile("sector")
	defer nigiri.StopProfile("sector")

	Q = nigiri.NewQueue()
	C = nigiri.NewCamera()
	C.SetCenter(vec2.V2{X: 300, Y: 300})
	C.SetClipRect(image.Rect(200, 200, 400, 400))

	tt, _ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	tt.Fill(colornames.Darkkhaki)
	ClipRect = nigiri.NewSprite(nigiri.NewTex(tt), 0, nil)
	ClipRect.Position = vec2.V(300, 300)
	ClipRect.Pivot = vec2.Center
	ClipRect.Scaler = nigiri.NewFixedScaler(200, 200)

	Sector = nigiri.SectorParams{
		Radius:   100,
		StartAng: 30,
		EndAng:   90,
	}.New(1, C)

	nigiri.Run(mainLoop, 600, 600, 1, "TEST")
}
