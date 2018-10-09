package main

import (
	"github.com/shnifer/nigiri/vec2"
	"math"
)

const minSinAlphaOptimisation = 0.1

type Circle struct {
	Center vec2.V2
	Radius float64
}

func (c Circle) FromPoint(p vec2.V2) (dist float64, period vec2.AnglePeriod) {
	V := c.Center.Sub(p)
	dist = V.Len()
	if dist == 0 {
		return dist, vec2.FullAnglePeriod
	} else if c.Radius <= 0 {
		return dist, rayAnglePeriod(V.Dir(), 0)
	} else if dist < c.Radius {
		return dist, vec2.FullAnglePeriod
	} else if dist == c.Radius {
		return dist, rayAnglePeriod(V.Dir(), 90)
	}

	var halfAng float64
	alpha := c.Radius / dist
	if alpha < minSinAlphaOptimisation {
		halfAng = alpha * vec2.Rad2Deg
	} else {
		halfAng = math.Asin(alpha) * vec2.Rad2Deg
	}
	return dist, rayAnglePeriod(V.Dir(), halfAng)
}

func rayAnglePeriod(midDir, halfAng float64) vec2.AnglePeriod {
	if halfAng < 0 {
		halfAng = 0
	}
	if halfAng >= 180 {
		return vec2.FullAnglePeriod
	}
	return vec2.NewAnglePeriod(midDir-halfAng, midDir+halfAng)
}
