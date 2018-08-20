package main

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

type Camera struct{
	posG *ebiten.GeoM
	rot float64
	scale float64
	clipRect image.Rectangle
}

func NewCamera()*Camera{
	res := &Camera{
		posG: new(ebiten.GeoM),
		scale:1,
		clipRect: image.ZR,
	}
	return res
}