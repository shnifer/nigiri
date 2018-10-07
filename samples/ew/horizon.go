package main

import "github.com/shnifer/nigiri/vec2"

type horizonPeriod struct{
	anglePeriod
	distance float64
	object HorizonCircler
}

type horizon struct{
	position vec2.V2
	periods []horizonPeriod
	angleMarks []float64
}

func newHorizon(position vec2.V2) *horizon{
	return &horizon{
		position:position,
		periods: make([]horizonPeriod,0),
		angleMarks: make([]float64, 0),
	}
}

func (h *horizon) Reset(position vec2.V2){
	h.position = position
	h.periods = h.periods[:0]
	h.angleMarks = h.angleMarks[:0]
}

func calcCircler (obj HorizonCircler, position vec2.V2) horizonPeriod{
	circle:=obj.HorizonCircle()
	return horizonPeriod{
		anglePeriod:circle.FromPoint(position),
		distance:circle.Center.Sub(position).Len(),
		object: obj,
	}
}

func (h *horizon) Add(obj HorizonCircler) {
	period := calcCircler(obj, h.position)
	if period.anglePeriod == EmptyAnglePeriod{
		return
	}

	h.periods = append(h.periods, period)
	h.addMark(period.Start)
	h.addMark(vec2.NormAng(period.End))
}