package main

import (
	"github.com/shnifer/nigiri/vec2"
)

type Horizon struct {
	////objects by roles
	//targets []HorizonObject
	//obstacles []HorizonObject
	//blockers []HorizonObject

	//current source position
	point   vec2.V2
	zone    vec2.AnglePeriod
	maxDist float64

	//temporary arrays
	blockAngles    []ObjectData
	targetAngles   []ObjectData
	obstacleAngles []ObjectData
	blockResolve   []vec2.AnglePeriod

	result []HorizonCalculus
}

func NewHorizon() *Horizon {
	res := &Horizon{
		zone:           vec2.FullAnglePeriod,
		blockAngles:    make([]ObjectData, 0),
		targetAngles:   make([]ObjectData, 0),
		obstacleAngles: make([]ObjectData, 0),
		blockResolve:   make([]vec2.AnglePeriod, 0),
		result:         make([]HorizonCalculus, 0),
	}
	return res
}

func (h *Horizon) SetPointZoneDist(p vec2.V2, zone vec2.AnglePeriod, maxDist float64) {
	h.point = p
	h.zone = zone
	h.maxDist = maxDist
}

type ObjectPart struct {
	Dir       float64
	PartStart float64
	PartEnd   float64
}
func TakeObjectPart(object, part vec2.AnglePeriod) ObjectPart{
	total:=object.Wide()
	if total == 0{
		return ObjectPart{
			PartStart: 0,
			PartEnd: 1,
			Dir: object.Medium(),
		}
	}
	start1,_:=object.Get()
	start2,_:=part.Get()
	startOff:=vec2.NewAnglePeriod(start1, start2).Wide()
	medOff:=part.Wide()
	return ObjectPart{
		PartStart:startOff/total,
		PartEnd:(startOff+medOff)/total,
		Dir: object.Medium(),
	}
}

type HorizonObjectPart struct {
	ObjectData
	ObjectPart
}

type HorizonCalculus struct {
	Target    HorizonObjectPart
	Obstacles []HorizonObjectPart
}

type ObjectData struct {
	Object HorizonObject
	Dist   float64
	Angles vec2.AnglePeriod
}

func (h *Horizon) Calculate(targets, obstacles, blockers []HorizonObject) []HorizonCalculus {
	h.calculateTemporary(targets, obstacles, blockers)
	h.calculateResult()

	return h.result
}

func (h *Horizon) calculateResult() {
	h.result = h.result[:0]
	for _, target := range h.targetAngles {
		var obs []HorizonObjectPart
		for _, obstacle := range h.obstacleAngles {
			if obstacle.Dist >= target.Dist {
				continue
			}
			cross, is := obstacle.Angles.Intersect(target.Angles)
			if !is {
				continue
			}
			_, total:=obstacle.Object.HorizonCircle().FromPoint(h.point)
			obs = append(obs, HorizonObjectPart{
				ObjectData: obstacle,
				ObjectPart: TakeObjectPart(total, cross),
			})
		}
		_, total:=target.Object.HorizonCircle().FromPoint(h.point)
		h.result = append(h.result, HorizonCalculus{
			Obstacles: obs,
			Target: HorizonObjectPart{
				ObjectData: target,
				ObjectPart: TakeObjectPart(total, target.Angles),
			},
		})
	}
}

func (h *Horizon) calculateTemporary(targets, obstacles, blockers []HorizonObject) {
	h.blockAngles = h.blockAngles[:0]
	h.targetAngles = h.targetAngles[:0]
	h.obstacleAngles = h.obstacleAngles[:0]

	var angles, cross vec2.AnglePeriod
	var dist float64
	var ok bool

	add := func(arr *[]ObjectData, obj HorizonObject, dist float64, angles vec2.AnglePeriod) {
		*arr = append(*arr, ObjectData{
			Object: obj,
			Dist:   dist,
			Angles: angles,
		})
	}

	for _, blocker := range blockers {
		dist, angles = blocker.HorizonCircle().FromPoint(h.point)
		if dist > h.maxDist && h.maxDist > 0 {
			continue
		}
		if cross, ok = h.zone.Intersect(angles); ok {
			add(&h.blockAngles, blocker, dist, cross)
		}
	}
	for _, target := range targets {
		dist, angles = target.HorizonCircle().FromPoint(h.point)
		if dist > h.maxDist && h.maxDist > 0 {
			continue
		}
		cross, ok = h.zone.Intersect(angles)
		if !ok {
			continue
		}

		h.applyBlockOnResolve(dist, cross)

		for _, res := range h.blockResolve {
			add(&h.targetAngles, target, dist, res)
		}
	}

	for _, obstacle := range obstacles {
		dist, angles = obstacle.HorizonCircle().FromPoint(h.point)
		if dist > h.maxDist && h.maxDist > 0 {
			continue
		}
		cross, ok = h.zone.Intersect(angles)
		if !ok {
			continue
		}

		for i,target:=range h.targetAngles{
			if target.Dist<=dist{
				continue
			}
			cross,ok:=target.Angles.Intersect(cross)
			if !ok{
				continue
			}
			start,end:=cross.Get()
			h.splitTarget(i, start)
			h.splitTarget(i, vec2.NormAng(end))
		}

		add(&h.obstacleAngles, obstacle, dist, cross)
	}
}

func (h *Horizon) applyBlockOnResolve(dist float64, cross vec2.AnglePeriod) {
	h.blockResolve = h.blockResolve[:0]
	h.blockResolve = append(h.blockResolve, cross)

	for _, block := range h.blockAngles {
		if block.Dist >= dist {
			continue
		}

		l := len(h.blockResolve)
		i := 0
		for i < l {
			n, p1, p2 := h.blockResolve[i].Sub(block.Angles)
			switch n {
			case 0:
				h.blockResolve[i] = h.blockResolve[l-1]
				h.blockResolve = h.blockResolve[:l-1]
				l--
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
}

func (h *Horizon) splitTarget(i int, angle float64) {
	start, end:=h.targetAngles[i].Angles.Get()
	if start==angle || vec2.NormAng(end)==angle{
		return
	}
	h.targetAngles[i].Angles = vec2.NewAnglePeriod(start, angle+360)
	h.targetAngles = append(h.targetAngles, h.targetAngles[i])
	h.targetAngles[len(h.targetAngles)-1].Angles = vec2.NewAnglePeriod(angle, end+360)
}