package main

import (
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
	"golang.org/x/image/colornames"
	"image/color"
)

const (
	particleSize = 25
)

type visualObject struct {
	data         *objectData
	mainSprite   nigiri.Sprite
	flagParticle nigiri.Sprite
	caption      nigiri.TextSprite
}

func (vo *visualObject) DrawReqs(Q *nigiri.Queue) {
	dat := vo.data
	vo.mainSprite.FixedH = dat.radius * 2
	vo.mainSprite.FixedW = dat.radius * 2
	vo.mainSprite.Position = dat.pos
	Q.Add(vo.mainSprite)

	r := vo.mainSprite.GetRect()
	r.Angle = 0
	corners := r.Corners()
	topLeft := corners[0]
	off := 0.0
	drawFlag := func(clr color.Color) {
		vo.flagParticle.SetColor(clr)
		vo.flagParticle.Position = topLeft.AddMul(vec2.V(particleSize, 0), off)
		Q.Add(vo.flagParticle)
		off += 1
	}
	if vo.data.r {
		drawFlag(colornames.Red)
	}
	if vo.data.g {
		drawFlag(colornames.Green)
	}
	if vo.data.b {
		drawFlag(colornames.Blue)
	}

	vo.caption.SetText(vo.data.name)
	vo.caption.Position = corners[2].Add(corners[3]).Mul(0.5)
	Q.Add(vo.caption)
}

func (vo *visualObject) Update(dt float64) {
	vo.mainSprite.Update(dt)
	vo.flagParticle.Update(dt)
	if vo.data.texType == 0 {
		vo.mainSprite.Angle += dt * 10
	}
}

func NewVisualObject(data *objectData, cam MyCam) *visualObject {
	var src nigiri.TexSrcer
	if data.texType == 0 {
		src, _ = nigiri.GetTex("HUD_Ship.png")
	} else {
		tex, _ := nigiri.GetTex("planet_ani.png")
		src, _ = nigiri.NewFrameTexSrc(tex, 64, 64, 19,
			nigiri.AnimateFrameCycle(10))
	}
	mainSprite := nigiri.NewSprite(src, 1, cam.Phys())
	mainSprite.Scaler = nigiri.NewFixedScaler(0, 0)
	mainSprite.Pivot = vec2.Center
	mainSprite.SetSmooth(true)

	particle, _ := nigiri.GetTex("particle.png")
	particleSprite := nigiri.NewSprite(particle, 2, cam.Local())
	particleSprite.Pivot = vec2.Center
	particleSprite.Scaler = nigiri.NewFixedScaler(particleSize, particleSize)

	captionSprite := nigiri.NewTextSprite(1.2, false, 2,
		cam.Local())
	face, _ := nigiri.GetFace("furore.ttf", 16)
	captionSprite.DefFace = face

	res := &visualObject{
		data:         data,
		mainSprite:   mainSprite,
		flagParticle: particleSprite,
		caption:      captionSprite,
	}

	return res
}
