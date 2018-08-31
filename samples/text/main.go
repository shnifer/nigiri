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
	Q.Add(TD)
	for i := 0; i < 10; i++ {
		S.ScaleFactor = v2.V(1, 1).Mul(float64((i % 2) + 1))
		S.Angle = float64(10 * i)
		S.Position = v2.V(100, 100).Mul(float64(i))
		Q.Add(S)
	}
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v\nDraws: %v", ebiten.CurrentFPS(), Q.Len()))
	return nil
}

func main() {
	nigiri.StartProfile("cam")
	defer nigiri.StopProfile("cam")

	Q = nigiri.NewQueue()

	C = nigiri.NewCamera()
	C.SetCenter(v2.V2{X: 500, Y: 500})
	C.SetScale(5)
	C.SetClipRect(image.Rect(300, 300, 700, 700))

	nigiri.SetFaceLoader(nigiri.FileFaceLoader("samples"))

	TS = nigiri.NewTextSrc(1.2, 1)
	face, err := nigiri.GetFace("furore.ttf", 20)
	if err != nil {
		panic(err)
	}

	TD = nigiri.NewTextDrawer(face, 2)
	TD.Position = v2.V(100, 100)
	TD.Color = colornames.Brown
	TD.Text = "just simple textdrawer"

	TS.SetText("text source sample\n multi line text", face, nigiri.AlignLeft, colornames.Green)

	S = nigiri.SpriteOpts{
		Src:          TS,
		Pivot:        v2.Center,
		Smooth:       false,
		CamTransform: C.Phys(),
	}.New()

	ebiten.SetVsyncEnabled(true)
	ebiten.Run(mainLoop, 1000, 1000, 1, "TEST")
}
