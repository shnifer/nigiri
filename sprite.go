package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"github.com/Shnifer/nigiri/v2"
)

type Sprite struct{
	Src           SpriteSrcer
	CompositeMode ebiten.CompositeMode
	Layer float32
	GOpts GOpts
	RectScaler
	filter ebiten.Filter
	color color.Color
	alpha float64
	colorM ebiten.ColorM
	rect
	tagSuffix string
}

type GOpts struct{
	denyRotation bool
	denyScale bool
}

func NewSprite(src SpriteSrcer, rectScaler RectScaler) *Sprite{
	res := &Sprite{
		Src:    src,
		RectScaler: rectScaler,
		color: color.White,
		alpha: 1,
		rect: rect{
			pivotRel:v2.V2{X:0.5, Y:0.5},
		},
	}
	res.calcTagSuffix()
	return res
}

func (s *Sprite) SetSmooth(smooth bool) {
	if smooth{
		s.filter = ebiten.FilterLinear
	} else {
		s.filter = ebiten.FilterDefault
	}
	s.calcTagSuffix()
}

func (s *Sprite) calcColorM() {
	const MaxColor = 0xffff
	s.colorM.Reset()
	r, g, b, a := s.color.RGBA()
	s.colorM.Scale(s.alpha*float64(r)/MaxColor,
		s.alpha*float64(g)/MaxColor,
			s.alpha*float64(b)/MaxColor,
				s.alpha*float64(a)/MaxColor)
}

func (s *Sprite) ColorAlpha() (color.Color, float64) {
	return s.color, s.alpha
}

func (s *Sprite) SetColor (color color.Color)  {
	if s.color==color{
		return
	}
	s.color = color
	s.calcColorM()
	s.calcTagSuffix()
}

func (s *Sprite) SetAlpha (alpha float64)  {
	if s.alpha == alpha{
		return
	}
	s.alpha = alpha
	s.calcColorM()
	s.calcTagSuffix()
}

func (s *Sprite) calcTagSuffix() {
	r,g,b,a:=s.color.RGBA()
	const k = 0xff
	br,bg,bb,ba:=byte(r/k), byte(g/k), byte(b/k), byte(a/k)
	baa:=byte(s.alpha*k)

	s.tagSuffix = string([]byte{byte(s.filter), br,bg,bb,ba,baa})
}

func (s *Sprite) Update(dt float64){
	if s.Src!=nil {
		s.Src.Update(dt)
	}
}

func (s *Sprite) DrawReqs(Q *Queue){
	if s.Src==nil{
		return
	}
	tex, srcRect, tag:= s.Src.GetSpriteSrc()

	w,h := srcRect.Max.X, srcRect.Max.Y
	if w==0 || h==0{
		return
	}

	if s.RectScaler==nil{
		s.rect.setSize(float64(w), float64(h))
	} else {
		s.rect.setSize(s.RectScale(w, h))
	}
	if Q.cam.isClipped(s.rect, s.GOpts){
		return
	}

	order:=reqOrder{
		layer: s.Layer,
		groupTag: tag+s.tagSuffix,
	}

	do:=&ebiten.DrawImageOptions{
		CompositeMode: s.CompositeMode,
		Filter: s.filter,
		ColorM: s.colorM,
		SourceRect: srcRect,
		GeoM: s.geom(float64(w),float64(h), Q.cam, s.GOpts),
	}
	Q.add(drawReq{
		f: ebiDrawF(tex, do),
		reqOrder: order,
	})
}

func (s *Sprite) geom(srcW,srcH float64,cam *Camera, gOpt GOpts) ebiten.GeoM {
	pos:=s.rect.pos
	rot:=s.rect.ang
	scale := 1.0
	if cam!=nil {
		pos = cam.applyV2(s.rect.pos)
		if !gOpt.denyRotation {
			rot += cam.ang
		}
		if !gOpt.denyScale {
			scale = cam.scale
		}
	}
	g:=ebiten.GeoM{}
	g.Translate((0.5-s.rect.pivotRel.X)*srcW, (0.5-s.rect.pivotRel.Y)*srcH)
	g.Scale(s.rect.width/srcW*scale, s.rect.height/srcH*scale)
	g.Rotate(rot*v2.Deg2Rad)
	g.Translate(pos.X, pos.Y)
	return g
}

func ebiDrawF(img *ebiten.Image, do *ebiten.DrawImageOptions) drawF{
	return func(dest *ebiten.Image) {
		dest.DrawImage(img, do)
	}
}