package nigiri

import "github.com/shnifer/nigiri/vec2"

type SpriteTrans struct {
	Position vec2.V2
	Pivot    vec2.V2
	Angle    float64
	Scaler
}

func (s SpriteTrans) TransformRect(rect Rect) Rect {
	rect = s.Scaler.TransformRect(rect)
	rect.Position = s.Position
	rect.Angle = s.Angle
	rect.pivot = s.Pivot
	return rect
}

type Sprite struct {
	*SpriteTrans
	*Drawer
}

func NewSprite(src TexSrcer, layer Layer, transforms ...Transformer) Sprite {
	SpriteT := &SpriteTrans{
		Scaler: NewScaler(1),
	}
	res := Sprite{
		SpriteTrans: SpriteT,
		Drawer:      NewDrawer(src, layer, append(Transforms{SpriteT}, transforms...)),
	}
	return res
}
