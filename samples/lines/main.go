package main

import (
	"github.com/Shnifer/nigiri"
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
	"image/color"
)

var C *nigiri.Camera
var Q *nigiri.Queue
var L nigiri.Line

func mainLoop(win *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		C.Translate(v2.InDir(C.Angle()).Rotate90().Mul(1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		C.Translate(v2.InDir(C.Angle()).Rotate90().Mul(-1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		C.Translate(v2.InDir(C.Angle()).Mul(1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		C.Translate(v2.InDir(C.Angle()).Mul(-1 / C.Scale()))
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

	Q.Clear()
	Q.Add(L)
	L.From = v2.ZV
	colors := [3]color.Color{colornames.Green, colornames.Palevioletred, colornames.Aquamarine}
	for i := 0; i < 360; i += 10 {
		L.To = v2.InDir(float64(i)).Mul(100)
		L.Color = colors[(i/10)%3]
		L.Width = 1
		Q.Add(L)
	}
	for i := 0; i < 10; i++ {
		v := v2.V(1.1, 0).Mul(float64(i))
		L.From = v2.V(100, 0).Add(v)
		L.To = L.From.Add(v2.V(0, 100))
		L.Color = colors[(i/10)%3]
		L.Width = 1
		Q.Add(L)
	}
	Q.Run(win)
	return nil
}

func main() {
	Q = nigiri.NewQueue()

	C = nigiri.NewCamera()
	C.SetCenter(v2.V2{X: 300, Y: 300})
	C.SetScale(5)

	L = nigiri.NewLine(C.NoRot(), 1)

	ebiten.Run(mainLoop, 600, 600, 1, "TEST")
}
