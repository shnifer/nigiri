package main

import (
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
	"math"
)

const RadiusN = 8
const LevelsN = 3

type Sphere struct{
	Values [RadiusN][LevelsN] float64
	Center vec2.V2
	Radius float64
	PointSprite *nigiri.Sprite
}

func levelAng(n int) float64{
	m:=float64(n)
	return m*180/(2*m+1)
}

func (s *Sphere) DrawReqs(Q *nigiri.Queue) {
	for rN:=0; rN<RadiusN; rN++{
		for lN:=0; lN<LevelsN; lN++{
			dir:=float64(rN)/RadiusN*360
			dst:=math.Cos(levelAng(lN)*vec2.Deg2Rad)*s.Radius
			s.PointSprite.Position = s.Center.AddMul(vec2.InDir(dir), dst)
			Q.Add(s.PointSprite)
		}
	}
}


