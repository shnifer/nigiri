package nigiri

import (
	"github.com/shnifer/nigiri/vec2"
)

var ZR Rect

type Rect struct {
	//position of pivot point in world coordinates
	Position vec2.V2
	//relative position of pivot-point
	pivot vec2.V2
	//size of rect, negative numbers are okay as that is a flipped axis
	Width  float64
	Height float64
	//Rotation of rect in Degrees, counter clockwise
	Angle float64
}

func (r Rect) TransformRect(rect Rect) Rect {
	return r
}

func NewRect(w, h float64, pivotRel vec2.V2) Rect {
	return Rect{
		pivot:  pivotRel,
		Width:  w,
		Height: h,
	}
}

func (r Rect) Empty() bool {
	return r.Width == 0 || r.Height == 0
}

func (r Rect) Corners() (res [4]vec2.V2) {
	rF := vec2.RotateF(r.Angle)
	p := vec2.V2{X: r.Width * r.pivot.X, Y: r.Height * r.pivot.Y}
	res[0] = r.Position.Add(rF(vec2.V2{X: 0, Y: 0}.Sub(p)))
	res[1] = r.Position.Add(rF(vec2.V2{X: r.Width, Y: 0}.Sub(p)))
	res[2] = r.Position.Add(rF(vec2.V2{X: r.Width, Y: r.Height}.Sub(p)))
	res[3] = r.Position.Add(rF(vec2.V2{X: 0, Y: r.Height}.Sub(p)))
	return res
}
