package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri/vec2"
	"image"
)

type Updater interface {
	Update(dt float64)
}

type DrawRequester interface {
	DrawReqs(Q *Queue)
}

type TexSrcer interface {
	GetSrcRectUID() (srcRect *image.Rectangle, uid uint64)
	GetSrcImage() *ebiten.Image
}

type RTransformer interface {
	TransformRect(rect Rect) Rect
}

type TransformerF func(in Rect) (out Rect, clip image.Rectangle)

func (f TransformerF) TransformRect(rect Rect) Rect {
	r, _ := f(rect)
	return r
}
func (f TransformerF) ClipRect() image.Rectangle {
	_, clip := f(ZR)
	return clip
}

type TriSrcer interface {
	GetVerticesIndices() ([]ebiten.Vertex, []uint16)
}

type VTransformer interface {
	ApplyPoint(vec2.V2) vec2.V2
}

//Apply t to both ends of vector (0,0)-(v) and takes delta t(v)-t(0,0)
func TransformVector(t VTransformer, v vec2.V2) vec2.V2 {
	return t.ApplyPoint(v).Sub(t.ApplyPoint(vec2.ZV))
}
func TransformUpVector(t VTransformer) vec2.V2 {
	return t.ApplyPoint(vec2.Up).Sub(t.ApplyPoint(vec2.ZV))
}

type Clipper interface {
	ClipRect() image.Rectangle
}

type OnMouser interface {
	OnMouse(x, y int) bool
}

type OnMouserF func(x, y int) bool

func (f OnMouserF) OnMouse(x, y int) bool {
	return f(x, y)
}
