package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri/vec2"
	"image"
	"math"
)

type Camera struct {
	center vec2.V2
	pos    vec2.V2
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

func (c *Camera) Center() vec2.V2 {
	return c.center
}
func (c *Camera) Pos() vec2.V2 {
	return c.pos
}
func (c *Camera) Scale() float64 {
	return c.scale
}
func (c *Camera) Angle() float64 {
	return c.ang
}
func (c *Camera) ClipRect() image.Rectangle {
	return c.clipRect
}

func (c *Camera) Translate(delta vec2.V2) {
	c.SetPos(c.pos.Add(delta))
}

func (c *Camera) Rotate(ang float64) {
	c.SetAng(c.ang + ang)
}

func (c *Camera) MulScale(scaleK float64) {
	c.SetScale(c.scale * scaleK)
}

func (c *Camera) SetCenter(v vec2.V2) {
	if v == c.center {
		return
	}
	c.center = v
	c.dirty = true
}

func (c *Camera) SetPos(v vec2.V2) {
	if v == c.pos {
		return
	}
	c.pos = v
	c.dirty = true
}

func (c *Camera) SetScale(v float64) {
	if v == c.scale {
		return
	}
	c.scale = v
	c.dirty = true
}

func (c *Camera) SetAng(v float64) {
	if v == c.ang {
		return
	}
	c.ang = v
	c.dirty = true
}

func (c *Camera) SetClipRect(rect image.Rectangle) {
	c.clipRect = rect
}

func (c *Camera) phys(rect Rect) Rect {
	rect.Width *= c.scale
	rect.Height *= c.scale
	rect.Angle += c.ang
	rect.Position = c.applyV2(rect.Position)
	if c.ClippedRect(rect) {
		rect = ZR
	}
	return rect
}

func (c *Camera) Phys() Transformer {
	return TransformerF(c.phys)
}

func (c *Camera) noRot(rect Rect) Rect {
	rect.Width *= c.scale
	rect.Height *= c.scale
	rect.Position = c.applyV2(rect.Position)
	if c.ClippedRect(rect) {
		rect = ZR
	}
	return rect
}

func (c *Camera) NoRot() Transformer {
	return TransformerF(c.noRot)
}

func (c *Camera) noScale(rect Rect) Rect {
	rect.Angle += c.ang
	rect.Position = c.applyV2(rect.Position)
	if c.ClippedRect(rect) {
		rect = ZR
	}
	return rect
}

func (c *Camera) NoScale() Transformer {
	return TransformerF(c.noScale)
}

func (c *Camera) local(rect Rect) Rect {
	if c.ClippedRect(rect) {
		rect = ZR
	}
	return rect
}

func (c *Camera) Local() Transformer {
	return TransformerF(c.local)
}

func (c *Camera) mark(rect Rect) Rect {
	rect.Position = c.applyV2(rect.Position)
	if c.ClippedRect(rect) {
		rect = ZR
	}
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
	c.posG.Rotate(-c.ang * vec2.Deg2Rad)
	c.posG.Translate(c.center.X, c.center.Y)
}

func (c *Camera) Apply(x, y float64) (float64, float64) {
	c.calcPosG()
	return c.posG.Apply(x, y)
}

func (c *Camera) applyV2(v vec2.V2) vec2.V2 {
	c.calcPosG()
	x, y := c.posG.Apply(v.X, v.Y)
	return vec2.V2{X: x, Y: y}
}

func (c *Camera) inClipRect(v vec2.V2) bool {
	x, y := int(v.X), int(v.Y)
	return image.Pt(int(x), int(y)).In(c.clipRect)
}

func (c *Camera) ClippedRect(rect Rect) bool {
	if c.clipRect.Empty() {
		return false
	}
	if c.inClipRect(rect.Position) {
		return false
	}
	px := math.Max(rect.pivot.X, 1-rect.pivot.X) * rect.Width
	py := math.Max(rect.pivot.Y, 1-rect.pivot.Y) * rect.Height
	dr := int(vec2.V2{X: px, Y: py}.Len())
	x, y := int(rect.Position.X), int(rect.Position.Y)
	cr := image.Rect(x-dr, y-dr, x+dr, y+dr)
	return cr.Intersect(c.clipRect).Empty()
}
