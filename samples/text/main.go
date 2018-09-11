package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image"
	"strconv"
)

var TD *nigiri.TextDrawer

var C *nigiri.Camera
var MultiText nigiri.TextSprite

var Q *nigiri.Queue
var Face font.Face

var UsedText nigiri.TextSprite

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
	Q.Add(TD)
	for i := 0; i < 3; i++ {
		MultiText.Position = vec2.V(0, 150).Mul(float64(i))
		Q.Add(MultiText)
	}
	for i := 0; i < 5; i++ {
		UsedText.SetText(strconv.Itoa(i))
		UsedText.Position = vec2.V(100, 150).AddMul(vec2.V(0, 100), float64(i))
		Q.Add(UsedText)
	}
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v\nDraws: %v", ebiten.CurrentFPS(), Q.Len()))
	return nil
}

func main() {
	nigiri.StartProfile("text")
	defer nigiri.StopProfile("text")

	Q = nigiri.NewQueue()

	C = nigiri.NewCamera()
	C.SetCenter(vec2.V2{X: 500, Y: 500})
	C.SetClipRect(image.Rect(0, 0, 1000, 1000))

	nigiri.SetFaceLoader(nigiri.FileFaceLoader("samples"))

	face, err := nigiri.GetFace("furore.ttf", 20)
	bigFace, err := nigiri.GetFace("furore.ttf", 30)
	if err != nil {
		panic(err)
	}
	Face = bigFace

	TD = nigiri.NewTextDrawer(face, 2)
	TD.Position = vec2.V(100, 100)
	TD.Color = colornames.Brown
	TD.Text = "just simple textdrawer\nsecond line"

	MultiText = nigiri.NewTextSprite(1.2, true, 0, C.Phys())
	MultiText.SetSmooth(true)
	MultiText.Pivot = vec2.Center

	MultiText.AddTextExt("text source sample\nmulti-line", face, 0, colornames.White)
	MultiText.AddTextExt("colored and sized", bigFace, 0, colornames.Greenyellow)
	MultiText.AddTextExt("center or", face, 1, colornames.White)
	MultiText.AddTextExt("right aligned", face, 2, colornames.White)

	UsedText = nigiri.NewTextSprite(1.2, false, 1, C.Phys())
	UsedText.Color = colornames.Red
	UsedText.Face = Face
	UsedText.SetText("text")
	UsedText.ChangeableTex = true

	ebiten.SetVsyncEnabled(true)
	nigiri.Run(mainLoop, 1000, 1000, 1, "TEST")
}
