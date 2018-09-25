package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

type DrawOptions struct {
	Layer Layer

	compositeMode ebiten.CompositeMode
	//Color
	color  color.Color
	alpha  float64
	colorM ebiten.ColorM

	filter ebiten.Filter

	//tagSuffix depends on Color and filter and CompositeMode
	drawTag uint64
}

func NewDrawOptions(layer Layer) DrawOptions {
	res := DrawOptions{
		color: color.White,
		alpha: 1,
		Layer: layer,
	}
	res.calcDrawTag()
	return res
}

func (do DrawOptions) CompositeMode() ebiten.CompositeMode {
	return do.compositeMode
}

func (do *DrawOptions) SetCompositeMode(mode ebiten.CompositeMode) {
	if do.compositeMode == mode {
		return
	}
	do.compositeMode = mode
	do.calcDrawTag()
}

func (do *DrawOptions) calcColorM() {
	const MaxColor = 0xffff
	do.colorM.Reset()
	r, g, b, a := do.color.RGBA()
	do.colorM.Scale(do.alpha*float64(r)/MaxColor,
		do.alpha*float64(g)/MaxColor,
		do.alpha*float64(b)/MaxColor,
		do.alpha*float64(a)/MaxColor)
}

func (do DrawOptions) ColorAlpha() (color.Color, float64) {
	return do.color, do.alpha
}

func (do *DrawOptions) SetColor(color color.Color) {
	if do.color == color {
		return
	}
	do.color = color
	do.calcColorM()
	do.calcDrawTag()
}

func (do *DrawOptions) SetAlpha(alpha float64) {
	if do.alpha == alpha {
		return
	}
	do.alpha = alpha
	do.calcColorM()
	do.calcDrawTag()
}

func (do *DrawOptions) calcDrawTag() {
	r, g, b, a := do.color.RGBA()
	const k = 0xff
	br, bg, bb, ba := uint64(r/k), uint64(g/k), uint64(b/k), uint64(a/k)
	bAlpha := uint64(do.alpha * k)
	bFilterAndComposite := uint64(do.compositeMode) + uint64(16*do.filter)
	do.drawTag = bFilterAndComposite<<40 + bAlpha<<32 + br<<24 + bg<<16 + bb<<8 + ba<<0
}
