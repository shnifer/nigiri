package main

import (
	"github.com/shnifer/nigiri/vec2"
	"math"
)

type Circle struct{
	Center vec2.V2
	Radius float64
}

func (c Circle) FromPoint(p vec2.V2) (dist float64, period vec2.AnglePeriod) {
	V:=c.Center.Sub(p)
	dist =V.Len()
	if c.Radius<=0 {
		return dist, vec2.EmptyAnglePeriod
	}
	var halfAng float64
	if  dist < c.Radius {
		halfAng = 180
	} else {
		halfAng = math.Asin(dist / c.Radius)
	}
	return  dist, rayAnglePeriod(V.Dir(), halfAng)
}

func rayAnglePeriod(midDir, halfAng float64) vec2.AnglePeriod {
	if halfAng<0{
		halfAng = 0
	}
	if halfAng>180{
		halfAng = 180
	}
	return vec2.NewAnglePeriod(midDir-halfAng, midDir+halfAng)
}