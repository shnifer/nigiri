package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
	_ "image/png"
)

var Q *nigiri.Queue
var C *nigiri.Camera
var Sprite nigiri.Sprite

func mainLoop(win *ebiten.Image, dt float64) error {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		Sprite.Position.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		Sprite.Position.X += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		Sprite.Position.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		Sprite.Position.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		Sprite.Angle += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		Sprite.Angle -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		Sprite.ScaleFactor.X *= -1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		Sprite.ScaleFactor.Y *= -1
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		Sprite.ScaleFactor.X *= 1.05
		Sprite.ScaleFactor.Y *= 1.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		Sprite.ScaleFactor.X /= 1.05
		Sprite.ScaleFactor.Y /= 1.05
	}

	Q.Clear()
	Q.Add(Sprite)
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("Position:%v\nAngle:%v\nScale:%v", Sprite.Position, Sprite.Angle, Sprite.ScaleFactor))
	return nil
}

func main() {
	nigiri.StartProfile("sprite")
	defer nigiri.StopProfile("sprite")

	Q = nigiri.NewQueue()
	nigiri.SetTexLoader(nigiri.FileTexLoader("samples"))
	tex, err := nigiri.GetTex("HUD_Ship.png")
	if err != nil {
		panic(err)
	}
	C := nigiri.NewCamera()
	C.SetCenter(vec2.V2{X: 300, Y: 300})

	Sprite = nigiri.NewSprite(tex, 0, C.Phys())
	Sprite.Pivot = vec2.Center
	Sprite.SetSmooth(true)

	nigiri.Run(mainLoop, 600, 600, 1, "TEST")
}
