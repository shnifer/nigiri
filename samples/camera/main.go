package camera

import (
	"fmt"
	"github.com/Shnifer/nigiri"
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/colornames"
	"image"
)

var T *ebiten.Image
var C *nigiri.Camera
var S *nigiri.Sprite
var Q *nigiri.Queue
var FR *nigiri.Sprite

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
	for x := -100; x <= 100; x += 10 {
		for y := -100; y <= 100; y += 10 {
			S.Position = v2.V2{X: float64(x), Y: float64(y)}
			Q.Add(S)
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
	C.SetCenter(v2.V2{X: 500, Y: 500})
	C.SetScale(5)
	C.SetClipRect(image.Rect(300, 300, 700, 700))

	tt, _ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	tt.Fill(colornames.Darkkhaki)

	FR = nigiri.SpriteOpts{
		Src:          nigiri.NewStatic(tt, nil, "particle"),
		Pivot:        v2.Center,
		Layer:        -1,
		CamTransform: C.Local(),
	}.New()
	FR.Position = v2.V2{500, 500}
	FR.UseFixed = true
	FR.FixedH, FR.FixedW = 400, 400

	Opts := nigiri.SpriteOpts{
		Src:          nigiri.NewStatic(tex, nil, "particle"),
		Pivot:        v2.Center,
		Smooth:       true,
		CamTransform: C.Phys(),
	}
	S = nigiri.NewSprite(Opts)
	S.ScaleFactor = v2.V2{X: 0.2, Y: 0.2}
	ebiten.SetVsyncEnabled(true)
	ebiten.Run(mainLoop, 1000, 1000, 1, "TEST")
}
