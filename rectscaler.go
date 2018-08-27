package nigiri

import (
	"math"
	"fmt"
	"github.com/Shnifer/nigiri/v2"
)

type Scaler struct {
	UseFixed bool
	ScaleFactor float64
	FitProportion bool
	FixedW, FixedH float64
}


func (r Scaler) RectScale(inW,inH int) (outW,outH float64) {
	if !r.UseFixed {
		return float64(inW)*r.ScaleFactor, float64(inH)*r.ScaleFactor
	}
	if !r.FitProportion{
		return r.FixedW, r.FixedH
	}
	if inW<=0 || inH<=0{
		return 0,0
	}
	w,h:=float64(inW), float64(inH)
	scale:=math.Min(r.FixedH/h, r.FixedW/w)
	return w*scale, h*scale
}

type rect struct{
	//position of pivot point in world coordinates
	pos v2.V2
	//position of center point in world coordinates
	center v2.V2
	width  float64
	height float64
	//relative position of pivotRel-point
	pivotRel v2.V2
	//Rotation of rect in Degrees, counter clockwise
	ang float64
}

//do not actually move rect, only change rotation point so pos is also changes
func (r *rect) setPivotRel(v v2.V2)  {
	if r.pivotRel == v{
		return
	}

	r.pivotRel = v
	r.calcPosFromCenter()
}

func (r *rect) calcPosFromCenter() {
	r.pos=r.center.Add(v2.V2{X:(r.pivotRel.X-0.5)*r.width, Y:(r.pivotRel.Y-0.5)*r.height}.Rotate(r.ang))
}
func (r *rect) calcCenterFromPos() {
	r.center = r.pos.Sub(v2.V2{X:(r.pivotRel.X-0.5)*r.width, Y:(r.pivotRel.Y-0.5)*r.height}.Rotate(r.ang))
}

func (r *rect) addAng(dAng float64) {
	if dAng==0{
		return
	}
	r.ang += dAng
	r.calcCenterFromPos()
}

func (r *rect) setAng(ang float64) {
	r.addAng(ang-r.ang)
}

func (r *rect) setSize(w,h float64){
	if w==r.width && h == r.height{
		return
	}
	r.width, r.height = w,h
	r.calcCenterFromPos()
}

func (r *rect) String() string{
	return fmt.Sprintf("rect:[center:%v, Size:%v, %v]", r.center, r.width, r.height)
}

func (r *rect) corners(cam *Camera, o GOpts) (res [4]v2.V2){
	ang:=r.ang
	if !o.DenyRotation {
		ang+=cam.ang
	}
	scale:=1.0
	if !o.DenyScale {
		scale = cam.scale
	}
	res[0] = r.center.Add(v2.V2{X:r.width/2,Y:r.height/2}.Mul(scale).Rotate(ang))
	res[1] = r.center.Add(v2.V2{X:-r.width/2,Y:r.height/2}.Mul(scale).Rotate(ang))
	res[2] = r.center.Add(v2.V2{X:-r.width/2,Y:-r.height/2}.Mul(scale).Rotate(ang))
	res[3] = r.center.Add(v2.V2{X:r.width/2,Y:-r.height/2}.Mul(scale).Rotate(ang))
	return res
}