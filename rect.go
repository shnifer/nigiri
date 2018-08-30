package nigiri

import (
	"github.com/Shnifer/nigiri/v2"
)

var ZR Rect

type Rect struct {
	//position of pivot point in world coordinates
	Position v2.V2
	//relative position of pivotRel-point
	pivotRel v2.V2
	//size of rect
	Width  float64
	Height float64
	//Rotation of rect in Degrees, counter clockwise
	Ang float64
}

func NewRect(w, h float64, pivotRel v2.V2) Rect {
	return Rect{
		pivotRel: pivotRel,
		Width:    w,
		Height:   h,
	}
}

func (r Rect) Empty() bool {
	return r == ZR
}

func (r Rect) Corners() (res [4]v2.V2) {
	rF := v2.RotateF(r.Ang)
	p := v2.V2{X: r.Width * r.pivotRel.X, Y: r.Height * r.pivotRel.Y}
	res[0] = r.Position.Add(rF(v2.V2{X: 0, Y: 0}.Sub(p)))
	res[1] = r.Position.Add(rF(v2.V2{X: r.Width, Y: 0}.Sub(p)))
	res[2] = r.Position.Add(rF(v2.V2{X: r.Width, Y: r.Height}.Sub(p)))
	res[3] = r.Position.Add(rF(v2.V2{X: 0, Y: r.Height}.Sub(p)))
	return res
}
