package main

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"github.com/Shnifer/nigiri/v2"
)

type Camera struct{
	posG     *ebiten.GeoM
	ang      float64
	scale    float64
	clipRect image.Rectangle
}

func NewCamera()*Camera{
	res := &Camera{
		posG: new(ebiten.GeoM),
		scale:1,
		clipRect: image.ZR,
	}
	return res
}

func (c *Camera) Apply(x,y float64) (float64,float64){
	return c.posG.Apply(x,y)
}

func (c *Camera) applyV2(v v2.V2) (v2.V2){
	x,y:=c.posG.Apply(v.X,v.Y)
	return v2.V2{X:x,Y:y}
}

func (c *Camera) inClipRect(v v2.V2) bool{
	x,y:=c.posG.Apply(v.X, v.Y)
	return image.Pt(int(x), int(y)).In(c.clipRect)
}

func (c *Camera) isClipped(r rect, o GOpts) bool{
	if c.clipRect==image.ZR{
		return false
	}
	if c.inClipRect(r.pos){
		return false
	}
	corners:=r.corners(r, o)
}