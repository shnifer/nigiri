package nigiri

/*
import (
	"github.com/hajimehoshi/ebiten"

	"image"
)

type StaticSrc struct {
	img  *ebiten.Image
	rect *image.Rectangle
	tag  string
}

func NewStatic(img *ebiten.Image, rect *image.Rectangle, tag string) StaticSrc {
	if rect == nil {
		rect = new(image.Rectangle)
		rect.Max.X, rect.Max.Y = img.Size()
	}
	return StaticSrc{
		img:  img,
		rect: rect,
		tag:  tag,
	}
}

func (s StaticSrc) GetSrcRect() (srcRect *image.Rectangle, tag string) {
	return s.rect, s.tag
}

func (s StaticSrc) GetSrcTex() (srcImage *ebiten.Image) {
	return s.img
}
*/
