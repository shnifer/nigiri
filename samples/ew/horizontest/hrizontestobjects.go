package main

import (
	"github.com/shnifer/nigiri"
	"image/color"
	"github.com/shnifer/nigiri/vec2"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
	"github.com/shnifer/nigiri/samples/ew"
)

type SolidObject struct {
	ew.DiffShadowBody
	*nigiri.Sprite
}

func NewSolidObject(circle ew.Circle) *SolidObject {
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
	ew.Circle
	Density float64
	*nigiri.Sprite
}

func (c *Cloud) HorizonCircle() ew.Circle {
	return c.Circle
}

func (c *Cloud) ShadowDensity(t ew.EmiType) (density float64) {
	return c.Density
}

func (c *Cloud) ShadowBlock() bool {
	return false
}

func NewCloud(circle ew.Circle, density float64) *Cloud{
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
	*ew.Horizon
	Color color.Color
}
func NewLight() *Light{
	k:=[ew.EmiDirCount]float64{}
	for i:=range k{
		k[i] = 1
	}
	res := &Light{
		LightEmitter: ew.NewLightEmitter(1, k, ""),
		Horizon: ew.NewHorizon(),
	}
	return res
}

func (l *Light) SetPosition(pos vec2.V2){
	l.LightEmitter.Center = pos
	l.Horizon.SetPointZoneDist(pos, vec2.FullAnglePeriod, 0)
}


type ViewSectorDrawer struct {
	*nigiri.TriDrawer
	Color color.Color
	Point vec2.V2
	Target ew.HorizonObjectPart
}

func NewViewSectorDrawer(layer nigiri.Layer, vTransformer nigiri.VTransformer) *ViewSectorDrawer{
	res:=&ViewSectorDrawer{}
	res.TriDrawer = nigiri.NewTriDrawer(res, layer, vTransformer)
	res.TriDrawer.ChangeableSrc = true
	return res
}

func (d *ViewSectorDrawer) GetVerticesIndices() ([]ebiten.Vertex, []uint16) {
	v:=make([]ebiten.Vertex,0)
	i:=make([]uint16,0)
	v = append(v, nigiri.VertexColor(d.Point, alpha(color.White, 0.7)))
	for i:=0;i<3;i++{
		dir:= d.Target.Angles.MedPart(float64(i)/2)
		pt:= d.Point.Add(vec2.InDir(dir).Mul(d.Target.Dist))
		var clr color.Color
		clr = color.White
		if i==1 && d.Color!=nil{
			clr = d.Color
		}
		clr = alpha(clr, 0.4)
		v = append(v, nigiri.VertexColor(pt, clr))
	}
	for m := 1; m+1 < len(v); m++ {
		i = append(i, 0, uint16(m), uint16(m+1))
	}
	return v, i
}

func alpha(clr color.Color, alpha float64) color.Color{
	r,g,b,a:=clr.RGBA()
	return color.RGBA64{
		R:uint16(float64(r)*alpha),
		G:uint16(float64(g)*alpha),
		B:uint16(float64(b)*alpha),
		A:uint16(float64(a)*alpha),
	}
}