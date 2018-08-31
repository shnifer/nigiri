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

type ActImage func(srcImage *ebiten.Image)

type TexSrcer interface {
	GetTexSrc() (srcImage *ebiten.Image, srcRect *image.Rectangle, tag string, afterDraw ActImage)
}

type Transformer interface {
	TransformRect(rect Rect) Rect
}
type TransformerF func(Rect) Rect

func (f TransformerF) TransformRect(rect Rect) Rect {
	return f(rect)
}
