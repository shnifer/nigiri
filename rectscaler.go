package nigiri

import (
	"math"
)

type Scaler struct {
	UseFixed       bool
	Sx, Sy         float64
	FitProportion  bool
	FixedW, FixedH float64
}

func (s Scaler) TransformRect(rect Rect) Rect {
	if rect.Width <= 0 || rect.Height <= 0 {
		return ZR
	}

	if !s.UseFixed {
		rect.Width *= s.Sx
		rect.Height *= s.Sy
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
