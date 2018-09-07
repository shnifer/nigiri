package main

import (
	"fmt"
	"github.com/Shnifer/nigiri"
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var C *nigiri.Camera
var Q *nigiri.Queue
var L nigiri.Line

var from vec2.V2
var dir float64

func mainLoop(win *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		from.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		from.X += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		from.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		from.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		dir += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		dir -= 1
	}

	Q.Clear()
	L.From = from
	L.To = from.AddMul(vec2.InDir(dir), 100)
	Q.Add(L)
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("Dir: %v", dir))
	return nil
}

func main() {
	from = vec2.V(200, 200)
	Q = nigiri.NewQueue()

	C = nigiri.NewCamera()
	C.SetCenter(vec2.V2{X: 300, Y: 300})

	L = nigiri.NewLine(nil, 1)

	ebiten.Run(mainLoop, 600, 600, 1, "TEST")
}
