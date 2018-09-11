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

var C *nigiri.Camera
var Q *nigiri.Queue
var FR nigiri.Sprite
var Partice nigiri.Sprite

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

	Q.Clear()
	for x := -100; x <= 100; x += 10 {
		for y := -100; y <= 100; y += 10 {
			Partice.Position = vec2.V(float64(x), float64(y))
			Q.Add(Partice)
		}
	}
	Q.Add(FR)
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v\nDraws: %v", ebiten.CurrentFPS(), Q.Len()))
	return nil
}

func main() {
	nigiri.StartProfile("cam")
	defer nigiri.StopProfile("cam")

	Q = nigiri.NewQueue()
	nigiri.SetTexLoader(nigiri.FileTexLoader("samples"))
	tex, err := nigiri.GetTex("particle.png")
	if err != nil {
		panic(err)
	}
	C = nigiri.NewCamera()
	C.SetCenter(vec2.V2{X: 500, Y: 500})
	C.SetScale(5)
	C.SetClipRect(image.Rect(300, 300, 700, 700))

	tt, _ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	tt.Fill(colornames.Darkkhaki)
	if err != nil {
		panic(err)
	}
	FR = nigiri.NewSprite(nigiri.NewTex(tt), 0, C.Local())
	FR.Position = vec2.V(500, 500)
	FR.Pivot = vec2.Center
	FR.Scaler = nigiri.NewFixedScaler(400, 400)

	Partice = nigiri.NewSprite(tex, 1, C.Phys())
	Partice.SetSmooth(true)
	Partice.Pivot = vec2.Center
	Partice.ScaleFactor = vec2.V(0.2, 0.2)

	ebiten.SetVsyncEnabled(false)
	nigiri.Run(mainLoop, 1000, 1000, 1, "TEST")
}
