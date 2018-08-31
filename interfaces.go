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
	GetSrcRect() (srcRect *image.Rectangle, tag string)
	GetSrcTex() (srcImage *ebiten.Image)
}

type Transformer interface {
	TransformRect(rect Rect) Rect
}

type TransformerF func(Rect) Rect

func (f TransformerF) TransformRect(rect Rect) Rect {
	return f(rect)
}
