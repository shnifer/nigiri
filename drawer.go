package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri/vec2"
	"image"
)

type Drawer struct {
	Src           TexSrcer
	Transform     Transformer
	ChangeableTex bool

	DrawOptions

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
		Src:         src,
		Transform:   Transforms(transform),
		DrawOptions: NewDrawOptions(layer),
	}
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

func (id *Drawer) GetRect() Rect {
	srcRect, _ := id.Src.GetSrcRectUID()
	return id.calcRect(srcRect)
}

func (id *Drawer) calcRect(srcRect *image.Rectangle) Rect {
	if srcRect == nil {
		return ZR
	}
	w, h := float64(srcRect.Dx()), float64(srcRect.Dy())
	if w == 0 || h == 0 {
		return ZR
	}
	id.r = NewRect(w, h, vec2.ZV)
	if id.Transform != nil {
		id.r = id.Transform.TransformRect(id.r)
	}
	if id.r.Empty() {
		return ZR
	}
	return id.r
}

func (id *Drawer) DrawReqs(Q *Queue) {
	if id.Src == nil {
		return
	}

	srcRect, uid := id.Src.GetSrcRectUID()
	w, h := float64(srcRect.Dx()), float64(srcRect.Dy())
	id.r = id.calcRect(srcRect)

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
