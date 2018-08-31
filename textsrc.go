package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"strings"
)

type Align byte

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

type textString struct {
	str   string
	face  font.Face
	color color.Color
	align Align

	strH    int
	bound   image.Rectangle
	advance int
}

func newTextString(str string, face font.Face, color color.Color, align Align, interlinesK float64) textString {
	strH := int(float64(face.Metrics().Height.Ceil()) * interlinesK)

	b, a := font.BoundString(face, str)
	r := image.Rect(b.Min.X.Floor(), b.Min.Y.Floor(), b.Max.X.Ceil(), b.Max.Y.Ceil())

	return textString{
		str:   str,
		face:  face,
		color: color,
		align: align,

		strH:    strH,
		bound:   r,
		advance: a.Round(),
	}
}

type TextSrc struct {
	strs []textString

	dirty bool
	offs  []image.Point

	InterLineK float64
	Layer      Layer
}

func (ts *TextSrc) GetTexSrc() (srcImage *ebiten.Image, srcRect *image.Rectangle, tag string, cb ActImage) {
	if ts.dirty {
		ts.dirty = false
		ts.recalcOffsets()
	}

	var rect image.Rectangle
	for i, s := range ts.strs {
		rect = rect.Union(s.bound.Add(ts.offs[i]))
	}

	w, h := rect.Dx(), rect.Dy()
	srcR := image.Rect(0, 0, w, h)
	tempImage := GetTempTex(w, h)
	for i, s := range ts.strs {
		text.Draw(tempImage, s.str, s.face,
			ts.offs[i].X, ts.offs[i].Y, s.color)
	}
	return tempImage, &srcR, "", PutTempTex
}

func (ts *TextSrc) recalcOffsets() {
	var maxA int
	for _, str := range ts.strs {
		if str.advance > maxA {
			maxA = str.advance
		}
	}
	var VOff, HOff int

	ts.offs = ts.offs[:0]

	for _, str := range ts.strs {
		switch str.align {
		case AlignLeft:
			HOff = 0
		case AlignCenter:
			HOff = (maxA - str.advance) / 2
		case AlignRight:
			HOff = maxA - str.advance
		default:
			//unknown align as default left and do not panic
			HOff = 0
		}
		ts.offs = append(ts.offs, image.Pt(HOff, VOff))
		VOff += str.strH
	}
}

func NewTextSrc(InterLineK float64, layer Layer) *TextSrc {
	res := &TextSrc{
		strs:       make([]textString, 0),
		offs:       make([]image.Point, 0),
		InterLineK: InterLineK,
		Layer:      layer,
	}
	return res
}

func (ts *TextSrc) Reset() {
	ts.strs = ts.strs[:0]
	ts.offs = ts.offs[:0]
}

func (ts *TextSrc) SetText(text string, face font.Face, align Align, clr color.Color) {
	ts.Reset()
	ts.AddText(text, face, align, clr)
}

func (ts *TextSrc) AddText(text string, face font.Face, align Align, clr color.Color) {
	ts.dirty = true
	strs := strings.Split(text, "\n")

	for _, str := range strs {
		ts.strs = append(ts.strs, textString{
			str:   str,
			face:  face,
			color: clr,
			align: align,
		})
	}
}
