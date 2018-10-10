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
)

var C *nigiri.Camera
var Q *nigiri.Queue
var ClipRect *nigiri.Sprite
var Particle *nigiri.Sprite

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
		Particle.ScaleFactor=Particle.ScaleFactor.Mul(1.05)
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		Particle.ScaleFactor=Particle.ScaleFactor.Mul(1/1.05)
	}
	log.Println("===")
	Q.Clear()
	const gridBounds= 5
	const gridStep = 10
	for x := -gridBounds; x <= gridBounds; x += gridStep {
		for y := -gridBounds; y <= gridBounds; y += gridStep {
			Particle.Position = vec2.V(float64(x), float64(y))
			switch ((y+gridBounds)/gridStep)%2 {
			case 0: Particle.SetColor(colornames.Yellow)
			case 1: Particle.SetColor(colornames.Blue)
			}
			Q.Add(Particle)
		}
	}
	Q.Add(ClipRect)
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v\nDraws: %v\nScale: %v",
		ebiten.CurrentFPS(), Q.Len(),C.Scale()))
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
	ClipRect = nigiri.NewSprite(nigiri.NewTex(tt), 0, C.Local())
	ClipRect.Position = vec2.V(500, 500)
	ClipRect.Pivot = vec2.Center
	ClipRect.Scaler = nigiri.NewFixedScaler(400, 400)

	Particle = nigiri.NewSprite(tex, 1, C.Phys())
	Particle.SetSmooth(true)
	Particle.Pivot = vec2.Center
	Particle.ScaleFactor = vec2.V(0.2, 0.2)

	ebiten.SetVsyncEnabled(true)
	nigiri.Run(mainLoop, 1000, 1000, 1, "TEST")
}
