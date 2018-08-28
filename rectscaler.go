package nigiri

import (
	"math"
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
