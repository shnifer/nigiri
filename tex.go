package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"errors"
	"image"
)

type Tex struct {
	image *ebiten.Image
	rect *image.Rectangle
	uid uint64
}

func (t Tex) GetSrcRectUID() (srcRect *image.Rectangle, uid uint64) {
	return t.rect, t.uid
}

func (t Tex) GetSrcImage() *ebiten.Image{
	return t.image
}

func NewTex (img *ebiten.Image) (Tex, error) {
	if img == nil {
		return Tex{}, errors.New("newTex can't use nil image")
	}
	return newTex(img), nil
}

func newTex(img *ebiten.Image) Tex{
	if img == nil{
		panic(errors.New("newTex can't use nil image"))
	}
	rect:=new(image.Rectangle)
	rect.Max.X, rect.Max.Y = img.Size()
	return Tex{
		image: img,
		rect: rect,
		uid: genUID(),
	}
}

func (t Tex) Size() (x,y int){
	return t.image.Size()
}

var uid uint64
func genUID() uint64{
	uid++
	log.Println("genUID", uid)
	return uid
}