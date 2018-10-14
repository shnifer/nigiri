package vistautils

import (
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
	"github.com/shnifer/nigiri/samples/ew/vista"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

type ViewSectorDrawer struct {
	*nigiri.TriDrawer
	Color color.Color
	Point vec2.V2
	Target vista.ObjectData
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
		dir:= d.Target.Area.Period.MedPart(float64(i)/2)
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