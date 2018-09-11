package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri/vec2"
	"image/color"
)

var defRectTex Tex
var lineImgDrawer *Drawer
var lineRect Rect

func init() {
	img, _ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	img.Fill(color.White)
	defRectTex = NewTex(img)

	lineImgDrawer = NewDrawer(defRectTex, 0, &lineRect)
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

type Polyline struct {
	CamTransform Transformer
	Layer        Layer
	Points []vec2.V2
	Closed bool
	Width        float64
	Color        color.Color
}

func NewPolyline(camTransform Transformer, layer Layer) Polyline {
	return Polyline{
		CamTransform: camTransform,
		Width:        1,
		Color:        color.White,
		Layer:        layer,
		Points: make([]vec2.V2,0),
	}
}

func NewPolylineExt(camTransform Transformer, layer Layer, points []vec2.V2,
	width float64, color color.Color) Polyline {
	return Polyline{
		CamTransform: camTransform,
		Width:        width,
		Color:        color,
		Points: points,
		Layer:        layer,
	}
}

func (l Polyline) DrawReqs(Q *Queue) {
	if l.Points == nil{
		return
	}

	lineImgDrawer.SetColor(l.Color)
	lineImgDrawer.Layer = l.Layer

	var to vec2.V2
	for i, from:=range l.Points {
		if i==len(l.Points)-1{
			if l.Closed{
				to = l.Points[0]
			} else {
				break
			}
		} else {
			to = l.Points[i+1]
		}
		lineRect.Position = from
		v := vec2.Sub(from, to)
		lineRect.Height = v.Len()
		lineRect.Angle = v.Dir()
		if l.CamTransform != nil {
			lineRect = l.CamTransform.TransformRect(lineRect)
		}
		lineRect.Width = l.Width

		Q.Add(lineImgDrawer)
	}
}