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

	dirty   bool
	offs    []image.Point
	rect    *image.Rectangle
	rectOff image.Point

	permanentTex bool
	permTex      Tex

	InterLineK float64
}

func (ts *TextSrc) GetSrcRectUID() (srcRect *image.Rectangle, uid uint64) {
	ts.recalc()
	if ts.permanentTex {
		uid = ts.permTex.uid
	}
	return ts.rect, uid
}

func (ts *TextSrc) GetSrcImage() *ebiten.Image {
	ts.recalc()

	if ts.permanentTex {
		return ts.permTex.image
	}

	w, h := ts.rect.Dx(), ts.rect.Dy()
	if w == 0 || h == 0 {
		return nil
	}
	tempImage := GetTempImage(w, h)
	ts.drawTextInto(tempImage)
	return tempImage
}

func (ts *TextSrc) drawTextInto(img *ebiten.Image) {
	if img == nil {
		panic("TextSrc.drawTextInto is not supposed to be called with nil ebiten.Image")
	}
	for i, s := range ts.strs {
		text.Draw(img, s.str, s.face,
			ts.offs[i].X-ts.rectOff.X, ts.offs[i].Y-ts.rectOff.Y, s.color)
	}
}
func (ts *TextSrc) recalc() {
	if !ts.dirty {
		return
	}
	ts.dirty = false

	var maxAdvance int
	for _, str := range ts.strs {
		if str.advance > maxAdvance {
			maxAdvance = str.advance
		}
	}
	var VOff, HOff int

	ts.offs = ts.offs[:0]

	for _, str := range ts.strs {
		switch str.align {
		case AlignLeft:
			HOff = 0
		case AlignCenter:
			HOff = (maxAdvance - str.advance) / 2
		case AlignRight:
			HOff = maxAdvance - str.advance
		default:
			//unknown align as default left and do not panic
			HOff = 0
		}
		ts.offs = append(ts.offs, image.Pt(HOff, VOff))
		VOff += str.strH
	}

	//rect containing all strings
	ts.rect = new(image.Rectangle)
	for i, s := range ts.strs {
		*ts.rect = ts.rect.Union(s.bound.Add(ts.offs[i]))
	}
	ts.rectOff = ts.rect.Min
	*ts.rect = ts.rect.Sub(ts.rect.Min)

	if ts.permanentTex {
		if ts.permTex.image != nil {
			PutPoolTex(ts.permTex)
		}
		w, h := ts.rect.Dx(), ts.rect.Dy()
		if w == 0 || h == 0 {
			ts.permTex = Tex{}
			return
		}
		ts.permTex = GetPoolTex(w, h)
		ts.drawTextInto(ts.permTex.image)
	}
}

func NewTextSrc(InterLineK float64, permanentTex bool) *TextSrc {
	res := &TextSrc{
		strs:         make([]textString, 0),
		offs:         make([]image.Point, 0),
		InterLineK:   InterLineK,
		permanentTex: permanentTex,
		dirty:        true,
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
		ts.strs = append(ts.strs, newTextString(str, face, clr, align, ts.InterLineK))
	}
}
