package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

type StaticSrc struct {
	img *ebiten.Image
	rect *image.Rectangle
	tag string
}

func NewStatic(img *ebiten.Image, rect *image.Rectangle, tag string) StaticSrc {
	if rect==nil{
		rect = new(image.Rectangle)
		rect.Max.X, rect.Max.Y = img.Size()
	}
	return StaticSrc{
		img:img,
		rect: rect,
		tag: tag,
	}
}

func (s StaticSrc) GetSpriteSrc() (srcImage *ebiten.Image, srcRect *image.Rectangle, tag string) {
	return s.img, s.rect, s.tag
}