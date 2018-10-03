package nigiri

import (
	"github.com/shnifer/nigiri/vec2"
	"image"
	"image/color"
)

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
	srcUpdate Updater
}

func (s Sprite) Update(dt float64) {
	if s.srcUpdate != nil {
		s.srcUpdate.Update(dt)
	}
}

func NewSprite(src TexSrcer, layer Layer, transforms ...RTransformer) *Sprite {
	SpriteT := &SpriteTrans{
		Scaler: NewScaler(1),
	}
	updater, _ := src.(Updater)
	res := &Sprite{
		SpriteTrans: SpriteT,
		Drawer:      NewDrawer(src, layer, append(Transforms{SpriteT}, transforms...)),
		srcUpdate:   updater,
	}
	return res
}

func (s *Sprite) GetSrcPoint(x,y int) (p image.Point, ok bool) {
	srcR,_:=s.Src.GetSrcRectUID()
	if srcR==nil{
		return image.ZP, false
	}
	resR:=s.Drawer.GetRect()
	if srcR.Empty() || resR.Empty() {
		return image.ZP, false
	}
	rel, ok :=resR.AbsToRel(vec2.V(float64(x), float64(y)))
	if !ok{
		return image.ZP, false
	}
	px:=srcR.Min.X+int(float64(srcR.Dx())*rel.X)
	py:=srcR.Min.Y+int(float64(srcR.Dy())*rel.Y)
	pt:=image.Pt(px,py)
	if !pt.In(*srcR) {
		return image.ZP, false
	}
	return pt, true
}

func (s *Sprite) GetSrcColor(x,y int) (clr color.Color, ok bool) {
	pt, ok:= s.GetSrcPoint(x,y)
	if !ok {
		return nil, false
	}
	img:=s.Src.GetSrcImage()
	return img.At(pt.X, pt.Y), true
}

func (s *Sprite) GetSrcPointColor(x,y int) (pt image.Point, clr color.Color, ok bool) {
	pt, ok = s.GetSrcPoint(x,y)
	if !ok {
		return image.ZP, nil, false
	}
	img:=s.Src.GetSrcImage()
	return pt, img.At(pt.X, pt.Y), true
}

type TextSprite struct {
	*TextSrc
	*SpriteTrans
	*Drawer
}

func NewTextSprite(interlineK float64, permanentTex bool, layer Layer, transforms ...RTransformer) *TextSprite {
	src := NewTextSrc(interlineK, permanentTex)
	SpriteT := &SpriteTrans{
		Scaler: NewScaler(1),
	}
	res := &TextSprite{
		TextSrc:     src,
		SpriteTrans: SpriteT,
		Drawer:      NewDrawer(src, layer, append(Transforms{SpriteT}, transforms...)),
	}
	return res
}
