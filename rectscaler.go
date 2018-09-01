package nigiri

import (
	"github.com/Shnifer/nigiri/v2"
	"math"
)

type Scaler struct {
	UseFixed       bool
	FitProportion  bool
	ScaleFactor    v2.V2
	FixedW, FixedH float64
}

func NewScaler(scale float64) Scaler {
	return Scaler{ScaleFactor: v2.V2{X: scale, Y: scale}}
}

func NewFixedScaler(w, h float64) Scaler {
	return Scaler{UseFixed: true, FixedH: h, FixedW: w}
}

func NewFitScaler(w, h float64) Scaler {
	return Scaler{UseFixed: true, FitProportion: true, FixedH: h, FixedW: w}
}

func (s Scaler) TransformRect(rect Rect) Rect {
	if rect.Width == 0 || rect.Height == 0 {
		return ZR
	}

	if !s.UseFixed {
		rect.Width *= s.ScaleFactor.X
		rect.Height *= s.ScaleFactor.Y
		return rect
	}

	if !s.FitProportion {
		rect.Width = s.FixedW
		rect.Height = s.FixedH
		return rect
	}

	scale := math.Min(s.FixedH/rect.Height, s.FixedW/rect.Width)
	rect.Width *= scale
	rect.Height *= scale
	return rect
}
