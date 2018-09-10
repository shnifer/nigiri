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

var T *ebiten.Image
var C *nigiri.Camera
var I *nigiri.Drawer
var S *nigiri.Sprite
var Q *nigiri.Queue

func mainLoop(win *ebiten.Image, dt float64) error {
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
	Q.Add(I)
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("Position:%v\nAngle:%v\nScale:%v", S.Position, S.Angle, S.ScaleFactor))
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

	S = &nigiri.Sprite{
		Scaler: nigiri.NewScaler(1),
		Pivot:  vec2.Center,
	}
	I = nigiri.NewDrawer(tex, nigiri.Transforms{S, C.Phys()})
	I.SetSmooth(true)

	nigiri.Run(mainLoop, 600, 600, 1, "TEST")
}
