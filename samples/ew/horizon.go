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
}

func newHorizon() *horizon{
	return &horizon{
		periods: make([]horizonPeriod,0),
	}
}

func (h *horizon) Reset(position vec2.V2){
	h.position = position
	h.periods = h.periods[:0]
}

func (h *horizon) Add(obj HorizonCircler) {
	circle:=obj.HorizonCircle()
	h.periods = append(h.periods, horizonPeriod{
		anglePeriod:circle.FromPoint(h.position),
		distance:circle.Center.Sub(h.position).Len(),
		object: obj,
	})
}