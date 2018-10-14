package main

import (
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
	"golang.org/x/image/colornames"
	"github.com/shnifer/nigiri/samples/ew"
	"github.com/shnifer/nigiri/samples/ew/vista"
)

type SolidObject struct {
	ew.DiffShadowBody
	*nigiri.Sprite
}


func NewSolidObject(circle vista.Circle) *SolidObject {
	Sprite := nigiri.NewSprite(nigiri.CircleTex(), 0, C.Phys())
	Sprite.Pivot = vec2.Center
	Sprite.SetSmooth(false)
	Sprite.SetColor(colornames.White)
	Sprite.Position = circle.Center
	visualSize := circle.Radius * 2
	Sprite.Scaler = nigiri.NewFixedScaler(visualSize, visualSize)

	return &SolidObject{
		DiffShadowBody: ew.DiffShadowBody{Circle: circle, Albedo: 1},
		Sprite:         Sprite,
	}
}
