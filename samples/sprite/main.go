package main

import (
	"fmt"
	"github.com/Shnifer/nigiri"
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	_ "image/png"
)

var T *ebiten.Image
var C *nigiri.Camera
var S *nigiri.Sprite
var Q *nigiri.Queue

func mainLoop(win *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		S.Position.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		S.Position.X += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		S.Position.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		S.Position.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		S.Angle += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		S.Angle -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		S.ScaleFactor.X *= -1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		S.ScaleFactor.Y *= -1
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		S.ScaleFactor.X *= 1.05
		S.ScaleFactor.Y *= 1.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		S.ScaleFactor.X /= 1.05
		S.ScaleFactor.Y /= 1.05
	}

	Q.Clear()
	Q.Add(S)
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("Position:%v\nAngle:%v\nScale:%v", S.Position, S.Angle, S.ScaleFactor))
	return nil
}

func main() {
	Q = nigiri.NewQueue()
	nigiri.SetTexLoader(nigiri.FileTexLoader("samples"))
	tex, err := nigiri.GetTex("HUD_Ship.png")
	if err != nil {
		panic(err)
	}
	C := nigiri.NewCamera()
	C.SetCenter(v2.V2{X: 300, Y: 300})

	Opts := nigiri.SpriteOpts{
		Src:          nigiri.NewStatic(tex, nil, "ship"),
		Pivot:        v2.TopMid,
		Smooth:       true,
		CamTransform: C.Phys(),
	}
	S = nigiri.NewSprite(Opts)

	ebiten.Run(mainLoop, 600, 600, 1, "TEST")
}
