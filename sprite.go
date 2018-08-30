package nigiri

import (
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
)

type SpriteOpts struct{
	Src TexSrcer
	CamTransform Transformer
	Pivot v2.V2
	CompositeMode ebiten.CompositeMode
	Layer Layer
}

type Sprite struct{
	imgD *ImageDrawer
	CamTransform Transformer
	Position v2.V2
	Angle float64
	Scaler
}

func (s *Sprite) DrawReqs(Q *Queue) {
	Q.Add(s.imgD)
}

func NewSprite(opts SpriteOpts) *Sprite{
	res:=&Sprite{
		CamTransform:    opts.CamTransform,
		Scaler: NewScaler(1),
	}
	res.imgD = NewImageDrawer(opts.Src, res, opts.Pivot)
	res.imgD.CompositeMode = opts.CompositeMode
	res.imgD.Layer = opts.Layer

	return res
}

func (s *Sprite) TransformRect(rect Rect) Rect{
	rect = s.Scaler.TransformRect(rect)
	rect.Position = s.Position
	rect.Ang = s.Angle
	if s.CamTransform!=nil {
		rect = s.CamTransform.TransformRect(rect)
	}
	return rect
}