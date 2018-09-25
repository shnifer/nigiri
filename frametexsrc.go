package nigiri

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	"image"
)

type FrameUpdF func(src *FrameTexSrc, dt float64)

type FrameTexSrc struct {
	tex          Tex
	sizeX, sizeY int
	cols         int
	count        int
	//we do NOT suppose to call this TexSrc on one frame with different rect
	rect   *image.Rectangle
	frameN int
	updF   FrameUpdF
}

func NewFrameTexSrc(tex Tex, sizeX, sizeY, count int, updF FrameUpdF) (*FrameTexSrc, error) {
	x, y := tex.Size()
	cols, rows := x/sizeX, y/sizeY
	if cols*rows < count {
		return nil, errors.New(fmt.Sprintf(
			"not enough space for %v frames size %vx%v on image with size %vx%v",
			count, sizeX, sizeY, x, y))
	}
	res := &FrameTexSrc{
		tex:   tex,
		sizeX: sizeX,
		sizeY: sizeY,
		cols:  cols,
		count: count,
		rect:  new(image.Rectangle),
		updF:  updF,
	}
	res.calcRect()
	return res, nil
}

func (f *FrameTexSrc) GetSrcRectUID() (srcRect *image.Rectangle, uid uint64) {
	return f.rect, f.tex.uid
}

func (f *FrameTexSrc) GetSrcImage() *ebiten.Image {
	return f.tex.image
}

func (f *FrameTexSrc) Count() int {
	return f.count
}
func (f *FrameTexSrc) FrameN() int {
	return f.frameN
}
func (f *FrameTexSrc) SetFrameN(n int) {
	if n < 0 {
		n = 0
	}
	if n > f.count-1 {
		n = f.count - 1
	}
	if n == f.frameN {
		return
	}
	f.frameN = n
	f.calcRect()
}
func (f *FrameTexSrc) calcRect() {
	col, row := f.frameN%f.cols, f.frameN/f.cols
	*f.rect = image.Rect(col*f.sizeX, row*f.sizeY,
		col*f.sizeX+f.sizeX, row*f.sizeY+f.sizeY)
}
func (f *FrameTexSrc) Update(dt float64) {
	if f.updF != nil {
		f.updF(f, dt)
	}
}

func AnimateFrameCycle(fps float64) FrameUpdF {
	var t float64
	step := 1.0 / fps
	return func(src *FrameTexSrc, dt float64) {
		t += dt
		n := int(t / step)
		if n == 0 {
			return
		}
		t -= step * float64(n)
		frame, count := src.FrameN(), src.Count()
		src.SetFrameN((frame + n) % count)
	}
}

func AnimateFrameOnce(fps float64) FrameUpdF {
	var t float64
	step := 1.0 / fps
	return func(src *FrameTexSrc, dt float64) {
		frame, count := src.FrameN(), src.Count()
		if frame == count-1 {
			return
		}
		t += dt
		n := int(t / step)
		if n == 0 {
			return
		}
		t -= step * float64(n)
		src.SetFrameN(frame + n)
	}
}
