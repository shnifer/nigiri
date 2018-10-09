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
var Circle nigiri.Sprite
var Sector nigiri.Sector

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
	Q.Add(Circle)
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
	C.SetClipRect(image.Rect(0, 0, 600, 600))

	Circle = nigiri.NewSprite(nigiri.CircleTex(), 0, C.Phys())
	Circle.Scaler = nigiri.NewFixedScaler(400, 400)
	Circle.Pivot = vec2.Center
	Circle.SetColor(colornames.Burlywood)

	Sector = nigiri.NewSector(1, C.Phys())
	Sector.Radius = 150
	Sector.SetColor(colornames.Blue)
	Sector.EndAng = 45

	nigiri.Run(mainLoop, 600, 600, 1, "TEST")
}
