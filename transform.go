package nigiri

import "github.com/Shnifer/nigiri/v2"

type Transforms []Transformer

func (t Transforms) TransformRect(rect Rect) Rect {
	for _, v := range t {
		if v == nil {
			continue
		}
		rect = v.TransformRect(rect)
		if rect.Empty() {
			return ZR
		}
	}
	return rect
}

type Sprite struct {
	Position v2.V2
	Pivot    v2.V2
	Angle    float64
	Scaler
}

func (s Sprite) TransformRect(rect Rect) Rect {
	rect = s.Scaler.TransformRect(rect)
	rect.Position = s.Position
	rect.Angle = s.Angle
	return rect
}
