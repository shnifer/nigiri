package nigiri

import (
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"image"
)

type Camera struct {
	center v2.V2
	pos    v2.V2
	scale  float64
	ang    float64
	dirty  bool
	posG   ebiten.GeoM

	clipRect image.Rectangle
}

func NewCamera() *Camera {
	res := &Camera{
		scale:    1,
		clipRect: image.ZR,
		dirty:    true,
	}
	return res
}

func (c *Camera) phys(rect Rect) Rect {
	rect.Width *= c.scale
	rect.Height *= c.scale
	rect.Ang += c.ang
	rect.Pos = c.applyV2(rect.Pos)
	return rect
}

func (c *Camera) Phys() Transformer {
	return TransformerF(c.phys)
}

func (c *Camera) noRot(rect Rect) Rect {
	rect.Width *= c.scale
	rect.Height *= c.scale
	rect.Pos = c.applyV2(rect.Pos)
	return rect
}

func (c *Camera) NoRot() Transformer {
	return TransformerF(c.noRot)
}

func (c *Camera) noScale(rect Rect) Rect {
	rect.Ang += c.ang
	rect.Pos = c.applyV2(rect.Pos)
	return rect
}

func (c *Camera) NoScale() Transformer {
	return TransformerF(c.noScale)
}

func (c *Camera) mark(rect Rect) Rect {
	rect.Pos = c.applyV2(rect.Pos)
	return rect
}

func (c *Camera) Mark() Transformer {
	return TransformerF(c.mark)
}

func (c *Camera) calcPosG() {
	if !c.dirty {
		return
	}
	c.dirty = false
	c.posG.Reset()
	c.posG.Translate(-c.pos.X, -c.pos.Y)
	c.posG.Scale(c.scale, c.scale)
	c.posG.Rotate(-c.ang)
	c.posG.Translate(c.center.X, c.center.Y)
}

func (c *Camera) Apply(x, y float64) (float64, float64) {
	c.calcPosG()
	return c.posG.Apply(x, y)
}

func (c *Camera) applyV2(v v2.V2) v2.V2 {
	c.calcPosG()
	x, y := c.posG.Apply(v.X, v.Y)
	return v2.V2{X: x, Y: y}
}

func (c *Camera) inClipRect(v v2.V2) bool {
	x, y := int(v.X), int(v.Y)
	return image.Pt(int(x), int(y)).In(c.clipRect)
}

func (c *Camera) ClipRect(rect Rect) bool {
	if c.clipRect == image.ZR {
		return false
	}
	if c.inClipRect(rect.Pos) {
		return false
	}
	dr := int(rect.Height + rect.Width)
	x, y := int(rect.Pos.X), int(rect.Pos.Y)
	cr := image.Rect(x-dr, y-dr, y+dr, y+dr)
	return cr.Intersect(c.clipRect).Empty()
}
