package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri/vec2"
	"image/color"
)

var defRectTex Tex
var lineImgDrawer *ImageDrawer
var lineRect Rect

func init() {
	img,_ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	img.Fill(color.White)
	defRectTex = newTex(img)

	lineImgDrawer = NewImageDrawer(defRectTex, &lineRect)
}

type Line struct {
	CamTransform Transformer
	Layer        Layer
	From         vec2.V2
	To           vec2.V2
	Width        float64
	Color        color.Color
}

func NewLine(camTransform Transformer, layer Layer) Line {
	return Line{
		CamTransform: camTransform,
		Width:        1,
		Color:        color.White,
		Layer:        layer,
	}
}

func NewLineExt(camTransform Transformer, layer Layer, from, to vec2.V2, width float64, color color.Color) Line {
	return Line{
		CamTransform: camTransform,
		Width:        width,
		Color:        color,
		From:         from,
		To:           to,
		Layer:        layer,
	}
}

func (l Line) DrawReqs(Q *Queue) {
	lineImgDrawer.SetColor(l.Color)
	lineImgDrawer.Layer = l.Layer

	lineRect.Position = l.From
	v := vec2.Sub(l.From, l.To)
	lineRect.Height = v.Len()
	lineRect.Angle = v.Dir()
	if l.CamTransform != nil {
		lineRect = l.CamTransform.TransformRect(lineRect)
	}
	lineRect.Width = l.Width

	Q.Add(lineImgDrawer)
}
