package main

import (
	"github.com/shnifer/nigiri"
	"image/color"
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

type Cloud struct{
	vista.Circle
	Density float64
	*nigiri.Sprite
}

func (c *Cloud) VistaTypes() (isObstacle, isTarget, isBlocker bool) {
	return true, false, false
}

func (c *Cloud) VistaCircle() vista.Circle {
	return c.Circle
}

func (c *Cloud) ShadowDensity(t ew.EmiType) (density float64) {
	return c.Density
}

func (c *Cloud) ShadowBlock() bool {
	return false
}

func NewCloud(circle vista.Circle, density float64) *Cloud{
	sprite := nigiri.NewSprite(nigiri.CircleTex(), 0, C.Phys())
	sprite.Pivot = vec2.Center
	sprite.SetSmooth(false)
	sprite.Position = circle.Center
	sprite.Scaler = nigiri.NewFixedScaler(circle.Radius * 2, circle.Radius * 2)
	sprite.SetColor(colornames.Gray)

	return &Cloud{
		Circle: circle,
		Density: density,
		Sprite: sprite,
	}
}

type Light struct{
	*ew.LightEmitter
	*vista.Vista
	Color color.Color
}
func NewLight() *Light{
	k:=[ew.EmiDirCount]float64{}
	for i:=range k{
		k[i] = 1
	}
	res := &Light{
		LightEmitter: ew.NewLightEmitter(1, k, ""),
		Vista: vista.NewArea(),
	}
	return res
}

func (l *Light) SetPosition(pos vec2.V2){
	l.LightEmitter.Center = pos
	l.Vista.Position = pos
}


