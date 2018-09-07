package nigiri

import (
	"github.com/shnifer/nigiri/vec2"
)

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
	Position vec2.V2
	Pivot    vec2.V2
	Angle    float64
	Scaler
}

func (s Sprite) TransformRect(rect Rect) Rect {
	rect = s.Scaler.TransformRect(rect)
	rect.Position = s.Position
	rect.Angle = s.Angle
	return rect
}
