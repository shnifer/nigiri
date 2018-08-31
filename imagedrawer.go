package nigiri

import (
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

type ImageDrawer struct {
	Src           TexSrcer
	CompositeMode ebiten.CompositeMode
	Layer         Layer
	Transform     Transformer
	pivot         v2.V2

	//color
	color  color.Color
	alpha  float64
	colorM ebiten.ColorM

	filter ebiten.Filter

	//tagSuffix depends on color and filter
	tagSuffix string

	//just a temp to not alloc, rewritten each draw
	r Rect
}

func NewImageDrawer(src TexSrcer, transform Transformer, pivotRel v2.V2) *ImageDrawer {
	res := &ImageDrawer{
		Src:       src,
		Transform: transform,
		color:     color.White,
		alpha:     1,
		pivot:     pivotRel,
	}
	res.calcTagSuffix()
	return res
}

func (id *ImageDrawer) SetSmooth(smooth bool) {
	if smooth {
		id.filter = ebiten.FilterLinear
	} else {
		id.filter = ebiten.FilterDefault
	}
	id.calcTagSuffix()
}

func (id *ImageDrawer) calcColorM() {
	const MaxColor = 0xffff
	id.colorM.Reset()
	r, g, b, a := id.color.RGBA()
	id.colorM.Scale(id.alpha*float64(r)/MaxColor,
		id.alpha*float64(g)/MaxColor,
		id.alpha*float64(b)/MaxColor,
		id.alpha*float64(a)/MaxColor)
}

func (id *ImageDrawer) ColorAlpha() (color.Color, float64) {
	return id.color, id.alpha
}

func (id *ImageDrawer) SetColor(color color.Color) {
	if id.color == color {
		return
	}
	id.color = color
	id.calcColorM()
	id.calcTagSuffix()
}

func (id *ImageDrawer) SetAlpha(alpha float64) {
	if id.alpha == alpha {
		return
	}
	id.alpha = alpha
	id.calcColorM()
	id.calcTagSuffix()
}

func (id *ImageDrawer) calcTagSuffix() {
	r, g, b, a := id.color.RGBA()
	const k = 0xff
	br, bg, bb, ba := byte(r/k), byte(g/k), byte(b/k), byte(a/k)
	baa := byte(id.alpha * k)

	id.tagSuffix = string([]byte{byte(id.filter), br, bg, bb, ba, baa})
}

func (id *ImageDrawer) DrawReqs(Q *Queue) {
	if id.Src == nil {
		return
	}
	tex, srcRect, tag, afterDrawCb := id.Src.GetTexSrc()

	if tex == nil {
		afterDrawCb(tex)
		return
	}
	w, h := float64(srcRect.Dx()), float64(srcRect.Dy())
	if w <= 0 || h <= 0 {
		afterDrawCb(tex)
		return
	}

	id.r = NewRect(w, h, id.pivot)

	if id.Transform != nil {
		id.r = id.Transform.TransformRect(id.r)
	}
	if id.r.Empty() {
		afterDrawCb(tex)
		return
	}

	order := reqOrder{
		layer:    id.Layer,
		groupTag: tag + id.tagSuffix,
	}

	do := &ebiten.DrawImageOptions{
		CompositeMode: id.CompositeMode,
		Filter:        id.filter,
		ColorM:        id.colorM,
		SourceRect:    srcRect,
		GeoM:          id.geom(w, h),
	}
	Q.add(drawReq{
		f:        ebiDrawF(tex, do, afterDrawCb),
		reqOrder: order,
	})
}

func (id *ImageDrawer) geom(w, h float64) (G ebiten.GeoM) {
	G.Translate(-w*id.pivot.X, -h*id.pivot.Y)
	G.Scale(id.r.Width/w, id.r.Height/h)
	G.Rotate(id.r.Ang * v2.Deg2Rad)
	G.Translate(id.r.Position.X, id.r.Position.Y)
	return G
}

func ebiDrawF(img *ebiten.Image, do *ebiten.DrawImageOptions, afterDrawCb ActImage) drawF {
	return func(dest *ebiten.Image) {
		dest.DrawImage(img, do)
		if afterDrawCb != nil {
			afterDrawCb(img)
		}
	}
}
