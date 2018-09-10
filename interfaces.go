package nigiri

import (
	"github.com/hajimehoshi/ebiten"
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
