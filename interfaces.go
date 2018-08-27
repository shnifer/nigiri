package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

type Updater interface{
	Update(dt float64)
}

type DrawRequester interface {
	DrawReqs(Q *Queue)
}

type RectScaler interface {
	RectScale(inW,inH int) (outW,outH float64)
}

type SpriteSrcer interface{
	Updater
	GetSpriteSrc() (srcImage *ebiten.Image, srcRect *image.Rectangle, tag string)
}
