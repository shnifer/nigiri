package main

import (
	"fmt"
	"github.com/Shnifer/nigiri"
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var T *ebiten.Image
var C *nigiri.Camera
var S *nigiri.Sprite
var Q *nigiri.Queue

func mainLoop(win *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		C.Translate(v2.InDir(C.Angle()).Rotate90().Mul(1))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		C.Translate(v2.InDir(C.Angle()).Rotate90().Mul(-1))
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		C.Translate(v2.InDir(C.Angle()).Mul(1))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		C.Translate(v2.InDir(C.Angle()).Mul(-1))
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
	for x := -100; x <= 100; x += 10 {
		for y := -100; y <= 100; y += 10 {
			S.Position = v2.V2{X: float64(x), Y: float64(y)}
			Q.Add(S)
		}
	}
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprint(ebiten.CurrentFPS()))
	return nil
}

func main() {
	Q = nigiri.NewQueue()
	nigiri.SetTexLoader(nigiri.FileTexLoader("samples"))
	tex, err := nigiri.GetTex("particle.png")
	if err != nil {
		panic(err)
	}
	C = nigiri.NewCamera()
	C.SetCenter(v2.V2{X: 600, Y: 600})

	Opts := nigiri.SpriteOpts{
		Src:          nigiri.NewStatic(tex, nil, "particle"),
		Pivot:        v2.Center,
		Smooth:       true,
		CamTransform: C.Phys(),
	}
	S = nigiri.NewSprite(Opts)
	S.ScaleFactor = v2.V2{X: 0.2, Y: 0.2}
	ebiten.SetVsyncEnabled(true)
	ebiten.Run(mainLoop, 1200, 1200, 1, "TEST")
}
