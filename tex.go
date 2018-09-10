package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"log"
)

type Tex struct {
	image *ebiten.Image
	rect  *image.Rectangle
	uid   uint64
}

func (t Tex) GetSrcRectUID() (srcRect *image.Rectangle, uid uint64) {
	return t.rect, t.uid
}

func (t Tex) GetSrcImage() *ebiten.Image {
	return t.image
}

func NewTex(img *ebiten.Image) Tex {
	if img == nil {
		return Tex{}
	}
	rect := new(image.Rectangle)
	rect.Max.X, rect.Max.Y = img.Size()
	return Tex{
		image: img,
		rect:  rect,
		uid:   UID(),
	}
}

func (t Tex) Size() (x, y int) {
	return t.image.Size()
}

var uid uint64

func UID() uint64 {
	uid++
	log.Println("UID", uid)
	return uid
}
