package nigiri

import (
	"math"
)

type Scaler struct {
	UseFixed bool
	Sx, Sy float64
	FitProportion bool
	FixedW, FixedH float64
}


func (r Scaler) RectScale(inW,inH int) (outW,outH float64) {
	if !r.UseFixed {
		return r.Sx, r.Sy
	}
	if inW<=0 || inH<=0{
		return 1,1
	}
	w,h:=float64(inW), float64(inH)
	if !r.FitProportion{
		return r.FixedW/w, r.FixedH/h
	}
	scale:=math.Min(r.FixedH/h, r.FixedW/w)
	return scale, scale
}
