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

type Transformer interface {
	TransformRect(rect Rect) Rect
}

type TransformerF func(Rect) Rect

func (f TransformerF) TransformRect(rect Rect) Rect {
	return f(rect)
}

type TriSrcer interface {
	GetVerticesIndices() ([]ebiten.Vertex, []uint16)
}

type VTransformer interface {
	ApplyV2(vec2.V2) vec2.V2
}

type Clipper interface {
	ClipRect() image.Rectangle
}
