package nigiri

//
//import (
//	"github.com/Shnifer/nigiri/v2"
//	"github.com/hajimehoshi/ebiten"
//	"image/color"
//)
//
//type SpriteOpts struct {
//	Src           TexSrcer
//	CamTransform  Transformer
//	Pivot         v2.V2
//	ChangeableTex bool
//	Smooth        bool
//	CompositeMode ebiten.CompositeMode
//	Layer         Layer
//}
//
//func (so SpriteOpts) New() *SpriteTrans {
//	return NewSprite(so)
//}
//
//type SpriteTrans struct {
//	imgD         *Drawer
//	rect         Rect
//	CamTransform Transformer
//	Position     v2.V2
//	Angle        float64
//	Scaler
//}
//
//func (s *SpriteTrans) DrawReqs(Q *Queue) {
//	Q.Add(s.imgD)
//}
//
//func NewSprite(opts SpriteOpts) *SpriteTrans {
//	res := &SpriteTrans{
//		CamTransform: opts.CamTransform,
//		Scaler:       NewScaler(1),
//	}
//	res.imgD = NewDrawer(opts.Src, res, opts.Pivot)
//	res.imgD.CompositeMode = opts.CompositeMode
//	res.imgD.ChangeableTex = opts.ChangeableTex
//	res.imgD.Layer = opts.Layer
//	res.imgD.SetSmooth(opts.Smooth)
//
//	return res
//}
//
//func (s *SpriteTrans) GetRect() Rect {
//	srcRect, _ := s.imgD.Src.GetSrcRect()
//	w, h := float64(srcRect.Dx()), float64(srcRect.Dy())
//	if w <= 0 || h <= 0 {
//		return ZR
//	}
//	return s.TransformRect(NewRect(w, h, pivot))
//}
//
//func (s *SpriteTrans) TransformRect(rect Rect) Rect {
//	s.rect = s.Scaler.TransformRect(rect)
//	s.rect.Position = s.Position
//	s.rect.Angle = s.Angle
//	if s.CamTransform != nil {
//		s.rect = s.CamTransform.TransformRect(s.rect)
//	}
//	return s.rect
//}
//
//func (s *SpriteTrans) SetColor(clr color.Color) {
//	s.imgD.SetColor(clr)
//}
//func (s *SpriteTrans) SetAlpha(v float64) {
//	s.imgD.SetAlpha(v)
//}
//
//func (s *SpriteTrans) ColorAlpha() (color.Color, float64) {
//	return s.imgD.color, s.imgD.alpha
//}
