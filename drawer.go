package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri/vec2"
	"image/color"
)

type Drawer struct {
	Src           TexSrcer
	Layer         Layer
	Transform     Transformer
	ChangeableTex bool

	compositeMode ebiten.CompositeMode
	//Color
	color  color.Color
	alpha  float64
	colorM ebiten.ColorM

	filter ebiten.Filter

	//tagSuffix depends on Color and filter and CompositeMode
	drawTag uint64

	//just a temp to not alloc, rewritten each draw
	r Rect
}

//Chache for DrawImageOptions. 96 bytes each
var doCache []*ebiten.DrawImageOptions

func init() {
	doCache = make([]*ebiten.DrawImageOptions, 0)
}
func getDo() *ebiten.DrawImageOptions {
	if len(doCache) == 0 {
		return new(ebiten.DrawImageOptions)
	}
	v := doCache[len(doCache)-1]
	doCache = doCache[:len(doCache)-1]
	return v
}
func putDo(do *ebiten.DrawImageOptions) {
	doCache = append(doCache, do)
}

func NewDrawer(src TexSrcer, layer Layer, transform ...Transformer) *Drawer {
	res := &Drawer{
		Src:       src,
		Transform: Transforms(transform),
		Layer:     layer,
		color:     color.White,
		alpha:     1,
	}
	res.calcDrawTag()
	return res
}

func (id *Drawer) SetSmooth(smooth bool) {
	if smooth {
		id.filter = ebiten.FilterLinear
	} else {
		id.filter = ebiten.FilterDefault
	}
	id.calcDrawTag()
}

func (id *Drawer) CompositeMode() ebiten.CompositeMode {
	return id.compositeMode
}

func (id *Drawer) SetCompositeMode(mode ebiten.CompositeMode) {
	if id.compositeMode == mode {
		return
	}
	id.compositeMode = mode
	id.calcDrawTag()
}

func (id *Drawer) calcColorM() {
	const MaxColor = 0xffff
	id.colorM.Reset()
	r, g, b, a := id.color.RGBA()
	id.colorM.Scale(id.alpha*float64(r)/MaxColor,
		id.alpha*float64(g)/MaxColor,
		id.alpha*float64(b)/MaxColor,
		id.alpha*float64(a)/MaxColor)
}

func (id *Drawer) ColorAlpha() (color.Color, float64) {
	return id.color, id.alpha
}

func (id *Drawer) SetColor(color color.Color) {
	if id.color == color {
		return
	}
	id.color = color
	id.calcColorM()
	id.calcDrawTag()
}

func (id *Drawer) SetAlpha(alpha float64) {
	if id.alpha == alpha {
		return
	}
	id.alpha = alpha
	id.calcColorM()
	id.calcDrawTag()
}

func (id *Drawer) calcDrawTag() {
	r, g, b, a := id.color.RGBA()
	const k = 0xff
	br, bg, bb, ba := uint64(r/k), uint64(g/k), uint64(b/k), uint64(a/k)
	bAlpha := uint64(id.alpha * k)
	bFilterAndComposite := uint64(id.compositeMode) + uint64(16*id.filter)
	id.drawTag = bFilterAndComposite<<40 + bAlpha<<32 + br<<24 + bg<<16 + bb<<8 + ba<<0
}

func (id *Drawer) DrawReqs(Q *Queue) {
	if id.Src == nil {
		return
	}

	srcRect, uid := id.Src.GetSrcRectUID()
	if srcRect == nil {
		return
	}
	w, h := float64(srcRect.Dx()), float64(srcRect.Dy())
	if w == 0 || h == 0 {
		return
	}
	id.r = NewRect(w, h, vec2.ZV)
	if id.Transform != nil {
		id.r = id.Transform.TransformRect(id.r)
	}
	if id.r.Empty() {
		return
	}

	order := reqOrder{
		layer:   id.Layer,
		uid:     uid,
		drawTag: id.drawTag,
	}

	do := getDo()
	*do = ebiten.DrawImageOptions{
		CompositeMode: id.compositeMode,
		Filter:        id.filter,
		ColorM:        id.colorM,
		SourceRect:    srcRect,
		GeoM:          id.geom(w, h),
	}
	if id.ChangeableTex {
		tex := id.Src.GetSrcImage()
		if tex == nil {
			return
		}
		Q.add(drawReq{
			f:        texDrawF(tex, do),
			reqOrder: order,
		})
	} else {
		Q.add(drawReq{
			f:        srcDrawF(id.Src, do),
			reqOrder: order,
		})
	}

}

func (id *Drawer) geom(w, h float64) (G ebiten.GeoM) {
	G.Translate(-w*id.r.pivot.X, -h*id.r.pivot.Y)
	G.Scale(id.r.Width/w, id.r.Height/h)
	G.Rotate(-id.r.Angle * vec2.Deg2Rad)
	G.Translate(id.r.Position.X, id.r.Position.Y)
	return G
}

func texDrawF(tex *ebiten.Image, do *ebiten.DrawImageOptions) drawF {
	return func(dest *ebiten.Image) {
		dest.DrawImage(tex, do)
		putDo(do)
		PutTempImage(tex)
	}
}

func srcDrawF(src TexSrcer, do *ebiten.DrawImageOptions) drawF {
	return func(dest *ebiten.Image) {
		tex := src.GetSrcImage()
		if tex == nil {
			return
		}
		dest.DrawImage(tex, do)
		putDo(do)
		PutTempImage(tex)
	}
}
