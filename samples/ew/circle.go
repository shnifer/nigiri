package main

import (
	"github.com/shnifer/nigiri/vec2"
	"math"
)

type Circle struct{
	Center vec2.V2
	Radius float64
}

type anglePeriod struct{
	Start, End float64
}

var EmptyAnglePeriod = anglePeriod{0,0}
var FullAnglePeriod = anglePeriod{0,360}

func rayAnglePeriod(midDir, halfAng float64) anglePeriod{
	if halfAng==0{
		return EmptyAnglePeriod
	}
	if halfAng>=180{
		return FullAnglePeriod
	}
	start:=vec2.NormAng(midDir - halfAng)
	return anglePeriod{
		Start:start,
		End:start+2*halfAng,
	}
}

func (c Circle) FromPoint(p vec2.V2) anglePeriod{
	if c.Radius<=0 {
		return EmptyAnglePeriod
	}
	V:=c.Center.Sub(p)
	dist :=V.Len()
	if  dist < c.Radius {
		return FullAnglePeriod
	}
	halfAng:=math.Asin(dist /c.Radius)
	return  rayAnglePeriod(V.Dir(), halfAng)
}