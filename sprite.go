package nigiri

import (
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

type SpriteOpts struct {
	Src           TexSrcer
	CamTransform  Transformer
	Pivot         v2.V2
	Smooth        bool
	CompositeMode ebiten.CompositeMode
	Layer         Layer
}

func (so SpriteOpts) New() *Sprite {
	return NewSprite(so)
}

type Sprite struct {
	imgD         *ImageDrawer
	rect         Rect
	CamTransform Transformer
	Position     v2.V2
	Angle        float64
	Scaler
}

func (s *Sprite) DrawReqs(Q *Queue) {
	Q.Add(s.imgD)
}

func NewSprite(opts SpriteOpts) *Sprite {
	res := &Sprite{
		CamTransform: opts.CamTransform,
		Scaler:       NewScaler(1),
	}
	res.imgD = NewImageDrawer(opts.Src, res, opts.Pivot)
	res.imgD.CompositeMode = opts.CompositeMode
	res.imgD.Layer = opts.Layer
	res.imgD.SetSmooth(opts.Smooth)

	return res
}

func (s *Sprite) GetRect() Rect {
	return s.rect
}

func (s *Sprite) TransformRect(rect Rect) Rect {
	s.rect = s.Scaler.TransformRect(rect)
	s.rect.Position = s.Position
	s.rect.Ang = s.Angle
	if s.CamTransform != nil {
		s.rect = s.CamTransform.TransformRect(s.rect)
	}
	return s.rect
}

func (s *Sprite) SetColor(clr color.Color) {
	s.imgD.SetColor(clr)
}
func (s *Sprite) SetAlpha(v float64) {
	s.imgD.SetAlpha(v)
}

func (s *Sprite) ColorAlpha() (color.Color, float64) {
	return s.imgD.color, s.imgD.alpha
}
