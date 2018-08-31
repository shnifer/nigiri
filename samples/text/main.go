package main

import (
	"fmt"
	"github.com/Shnifer/nigiri"
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/colornames"
	"image"
)

var TD *nigiri.TextDrawer
var TS *nigiri.TextSrc
var C *nigiri.Camera
var S *nigiri.Sprite
var Q *nigiri.Queue

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
	Q.Add(TD)
	Q.Add(S)
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v\nDraws: %v", ebiten.CurrentFPS(), Q.Len()))
	return nil
}

func main() {
	nigiri.StartProfile("text")
	defer nigiri.StopProfile("text")

	Q = nigiri.NewQueue()

	C = nigiri.NewCamera()
	C.SetCenter(v2.V2{X: 500, Y: 500})
	C.SetClipRect(image.Rect(0, 0, 1000, 1000))

	nigiri.SetFaceLoader(nigiri.FileFaceLoader("samples"))

	face, err := nigiri.GetFace("furore.ttf", 20)
	bigFace, err := nigiri.GetFace("furore.ttf", 30)
	if err != nil {
		panic(err)
	}

	TD = nigiri.NewTextDrawer(face, 2)
	TD.Position = v2.V(100, 100)
	TD.Color = colornames.Brown
	TD.Text = "just simple textdrawer\nsecond line"

	TS = nigiri.NewTextSrc(1.2, 1)
	TS.AddText("text source sample", face, 0, colornames.White)
	TS.AddText("multi-line", face, 0, colornames.White)
	TS.AddText("colored and sized", bigFace, 0, colornames.Greenyellow)
	TS.AddText("center or", face, 1, colornames.White)
	TS.AddText("right aligned", face, 2, colornames.White)

	S = nigiri.SpriteOpts{
		Src:          TS,
		Pivot:        v2.Center,
		Smooth:       true,
		CamTransform: C.Phys(),
	}.New()

	ebiten.SetVsyncEnabled(true)
	ebiten.Run(mainLoop, 1000, 1000, 1, "TEST")
}
