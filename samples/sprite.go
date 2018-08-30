package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
)

var T *ebiten.Image

func mainLoop(win *ebiten.Image) error {
	G := ebiten.GeoM{}
	G.Rotate(0.1)
	src := image.Rect(100, 100, 350, 350)
	do := &ebiten.DrawImageOptions{
		GeoM:       G,
		Filter:     ebiten.FilterLinear,
		SourceRect: &src,
	}
	win.DrawImage(T, do)
	return nil
}

func main() {
	tex, _, err := ebitenutil.NewImageFromFile("samples/HUD_Ship.png", ebiten.FilterDefault)
	T = tex
	if err != nil {
		panic(err)
	}
	ebiten.Run(mainLoop, 600, 600, 1, "TEST")
}
