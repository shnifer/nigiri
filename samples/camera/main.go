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

var T *ebiten.Image
var C *nigiri.Camera
var S *nigiri.SpriteTrans
var SI *nigiri.Drawer
var Q *nigiri.Queue
var FR *nigiri.Drawer

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
			S.Position = vec2.V(float64(x), float64(y))
			Q.Add(SI)
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
	frSprite := nigiri.SpriteTrans{
		Position: vec2.V(500, 500),
		Pivot:    vec2.Center,
		Scaler:   nigiri.NewFixedScaler(400, 400),
	}
	FR = nigiri.NewDrawer(nigiri.NewTex(tt), 0, nigiri.Transforms{frSprite, C.Local()})

	S = &nigiri.SpriteTrans{
		Pivot: vec2.Center,
	}
	S.ScaleFactor = vec2.V2{X: 0.2, Y: 0.2}
	SI = nigiri.NewDrawer(tex, 1, nigiri.Transforms{S, C.Phys()})
	SI.SetSmooth(true)
	SI.Layer = 1
	ebiten.SetVsyncEnabled(false)
	nigiri.Run(mainLoop, 1000, 1000, 1, "TEST")
}
