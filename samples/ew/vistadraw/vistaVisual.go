package main

import (
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/samples/ew/vista"
	"golang.org/x/image/colornames"
	"github.com/shnifer/nigiri/samples/ew/area"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"math"
	"github.com/shnifer/nigiri/vec2"
)

type VistaResultsSprite struct{
	nigiri.Tex
	*nigiri.Sprite
	scaleW, scaleH, maxH float64
}

func NewVistaResultsSprite(scaleW, scaleH, maxH float64, layer nigiri.Layer, cam *nigiri.Camera) *VistaResultsSprite{
	w,h:=int(scaleW*360), int(maxH*scaleH)
	res:=&VistaResultsSprite{
		scaleW: scaleW,
		scaleH:scaleH,
		maxH: maxH,
		Tex: nigiri.GetPoolTex(w,h),
	}
	res.Sprite = nigiri.NewSprite( res ,layer, cam.Local())
	return res
}

func (v *VistaResultsSprite) Take(res []vista.Result){
	image:=v.Tex.GetSrcImage()
	image.Fill(colornames.Gray)

	for _,rec:=range res{
		v.drawArea(image, rec.Target.Area, colornames.White)
	}
	for _,rec:=range res {
		for _, block := range rec.Blockers {
			v.drawArea(image, block.Area, colornames.Black)
		}
	}
}

func (v *VistaResultsSprite) drawArea(image *ebiten.Image, a area.Area, clr color.Color) {
	height:=math.Min(a.Height, v.maxH) * v.scaleH
	drawRect:=func(a1,a2 float64) {
		w:=(a2-a1)*v.scaleW
		rect:=nigiri.NewRect(w, height, vec2.BotRight)
		rect.Position = vec2.V(v.scaleW*360-a1*v.scaleW, v.maxH*v.scaleH)
		nigiri.DrawRect(image, rect, clr)
	}
	if a.Period.Start()<=a.Period.End() {
		drawRect(a.Period.Start(), a.Period.End())
	} else {
		drawRect(a.Period.Start(), 360)
		drawRect(0, a.Period.End())
	}
}