package nigiri

import (
	"github.com/Shnifer/nigiri/v2"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

var defRectTex *ebiten.Image
var lineImgDrawer *ImageDrawer
var lineRect Rect

func init() {
	defRectTex, _ = ebiten.NewImage(10, 10, ebiten.FilterDefault)
	defRectTex.Fill(color.White)

	lineImgDrawer = NewImageDrawer(NewStatic(defRectTex, nil, "!defRectTex"), &lineRect)
}

type Line struct {
	CamTransform Transformer
	Layer        Layer
	From         v2.V2
	To           v2.V2
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

func NewLineExt(camTransform Transformer, layer Layer, from, to v2.V2, width float64, color color.Color) Line {
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
	v := v2.Sub(l.From, l.To)
	lineRect.Height = v.Len()
	lineRect.Angle = v.Dir()
	if l.CamTransform != nil {
		lineRect = l.CamTransform.TransformRect(lineRect)
	}
	lineRect.Width = l.Width

	Q.Add(lineImgDrawer)
}
