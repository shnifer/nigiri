package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

var Particle *ebiten.Image
var Particle2 *ebiten.Image
var ScaleFactor float64

func mainLoop(win *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		ScaleFactor*=1.05
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		ScaleFactor/=1.05
	}
	do:=&ebiten.DrawImageOptions{
		Filter: ebiten.FilterLinear,//Filter this is critical
	}
	do.GeoM.Scale(ScaleFactor, ScaleFactor)//Scale factor is critical
	do.ColorM.Scale(1,1,0,1)//Draw with other color is critical
	win.DrawImage(Particle, do)
	do.GeoM.Translate(400,0)
	do.ColorM.Reset()
	do.ColorM.Scale(0,1,1,1)
	win.DrawImage(Particle, do)//and second draw for sure
	return nil
}



func main() {
	ScaleFactor = 1
	Particle,_ = ebiten.NewImage(400,400, ebiten.FilterDefault)
	Particle.Fill(color.White)
	Particle2,_ = ebiten.NewImage(400,400, ebiten.FilterDefault)
	Particle2.Fill(color.White)

	ebiten.Run(mainLoop, 1000, 1000, 1, "TEST")
}
