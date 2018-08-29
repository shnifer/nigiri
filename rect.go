package nigiri

import (
	"github.com/Shnifer/nigiri/v2"
)

type rect struct{
	//position of pivot point in world coordinates
	pos v2.V2
	//relative position of pivotRel-point
	pivotRel v2.V2
	//size of rect
	width  float64
	height float64
	//Rotation of rect in Degrees, counter clockwise
	ang float64
}

//do not actually move rect, only change rotation point so pos is also changes
func (r *rect) movePivotRel(v v2.V2)  {
	if r.pivotRel == v{
		return
	}

	r.pivotRel = v
}