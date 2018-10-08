package main

import "github.com/shnifer/nigiri/vec2"

type Horizon struct{
	////objects by roles
	//targets []HorizonObject
	//obstacles []HorizonObject
	//blockers []HorizonObject

	//current source position
	point vec2.V2
	zone vec2.AnglePeriod

	//temporary arrays
	blockAngles []objectData
	targetAngles []objectData
	blockResolve []vec2.AnglePeriod
}

func NewHorizon () *Horizon {
	res := &Horizon{
		zone: vec2.FullAnglePeriod,
		targetAngles: make([]objectData,0),
		blockResolve: make([]vec2.AnglePeriod,0),
	}
	return res
}

func (h *Horizon) SetPointZone(p vec2.V2, zone vec2.AnglePeriod){
	h.point = p
	h.zone = zone
}

type  HorizonCalculus struct{
	Target HorizonObject
	Obstacles []HorizonObject
}

type objectData struct{
	object HorizonObject
	dist float64
	angles vec2.AnglePeriod
}

func (h *Horizon) Calculate(targets,obstacles,blockers []HorizonObject) []HorizonCalculus{
	h.blockAngles = h.blockAngles[:0]
	h.targetAngles = h.targetAngles[:0]

	var angles, cross vec2.AnglePeriod
	var dist float64
	var ok bool

	add:=func(arr *[]objectData, obj HorizonObject, dist float64, angles vec2.AnglePeriod){
		*arr = append(*arr, objectData{
			object: obj,
			dist:dist,
			angles:angles,
		})
	}

	for _,blocker:=range blockers{
		dist, angles=blocker.HorizonCircle().FromPoint(h.point)
		if cross, ok=h.zone.Intersect(angles);ok{
			add(&h.blockAngles, blocker, dist, cross)
		}
	}

	for _,target:=range targets{
		dist, angles=target.HorizonCircle().FromPoint(h.point)
		cross, ok=h.zone.Intersect(angles)
		if !ok{
			continue
		}

		h.blockResolve = h.blockResolve[:0]
		h.blockResolve = append(h.blockResolve, cross)

		for _,block:=range h.blockAngles{
			if block.dist>=dist{
				continue
			}
			l := len(h.blockResolve)
			i := 0
			for i<l{
				n, p1, p2 := h.blockResolve[i].Sub(block.angles)
				switch n {
				case 0:
					h.blockResolve[i] = h.blockResolve[l-1]
					h.blockResolve = h.blockResolve[:l-1]
					l --
				case 1:
					h.blockResolve[i] = p1
					i++
				case 2:
					h.blockResolve[i] = p1
					h.blockResolve = append(h.blockResolve, p2)
					i++
				}
			}
		}

		for _,res:=range h.blockResolve{
			add(&h.targetAngles, target, dist, res)
		}
	}

	return []HorizonCalculus{}
}