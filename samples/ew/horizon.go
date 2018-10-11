package ew

import (
	"github.com/shnifer/nigiri/vec2"
	"log"
	"sort"
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

//Part is assumed to be totally inside object
func TakeObjectPart(object, part vec2.AnglePeriod) ObjectPart {
	total := object.Wide()
	if total == 0 {
		return ObjectPart{
			PartStart: 0,
			PartEnd:   1,
			Dir:       object.Medium(),
		}
	}
	start1, _ := object.Get()
	start2, _ := part.Get()
	startOff := vec2.NewAnglePeriod(start1, start2).Wide()
	medOff := part.Wide()
	return ObjectPart{
		PartStart: startOff / total,
		PartEnd:   (startOff + medOff) / total,
		Dir:       object.Medium(),
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

func (h *Horizon) Calculate(targets, obstacles, blockers []HorizonObject,
	ignoreSelf HorizonObject) []HorizonCalculus {

	h.calculateTemporary(targets, obstacles, blockers, ignoreSelf)
	h.calculateResult()

	return h.result
}

func (h *Horizon) calculateResult() {
	var cN int
	var r1 vec2.AnglePeriod
	h.result = h.result[:0]
	for _, target := range h.targetAngles {
		var obs []HorizonObjectPart
		for _, obstacle := range h.obstacleAngles {
			if obstacle.Dist >= target.Dist {
				continue
			}
			cN, r1, _ = obstacle.Angles.Intersect(target.Angles)
			if cN == 0 {
				continue
			}
			_, total := obstacle.Object.HorizonCircle().FromPoint(h.point)
			obs = append(obs, HorizonObjectPart{
				ObjectData: obstacle,
				ObjectPart: TakeObjectPart(total, r1),
			})
			if cN == 2 {
				log.Println("There is not waiting for obstacles 2-cross target parts")
			}
		}
		_, total := target.Object.HorizonCircle().FromPoint(h.point)
		h.result = append(h.result, HorizonCalculus{
			Obstacles: obs,
			Target: HorizonObjectPart{
				ObjectData: target,
				ObjectPart: TakeObjectPart(total, target.Angles),
			},
		})
	}
}

func addObjIntoArr(arr *[]ObjectData, obj HorizonObject, dist float64, angles vec2.AnglePeriod) {
	*arr = append(*arr, ObjectData{
		Object: obj,
		Dist:   dist,
		Angles: angles,
	})
}

func (h *Horizon) calculateTemporary(targets, obstacles, blockers []HorizonObject,
	ignoreSelf HorizonObject) {

	h.blockAngles = h.blockAngles[:0]
	h.targetAngles = h.targetAngles[:0]
	h.obstacleAngles = h.obstacleAngles[:0]

	var rN int
	var r1, r2 vec2.AnglePeriod
	var angles vec2.AnglePeriod
	var circle Circle
	var dist float64

	for _, blocker := range blockers {
		if blocker == ignoreSelf {
			continue
		}
		circle = blocker.HorizonCircle()
		dist, angles = circle.FromPoint(h.point)
		if dist-circle.Radius > h.maxDist && h.maxDist > 0 {
			continue
		}
		rN, r1, r2 = h.zone.Intersect(angles)
		switch rN {
		case 0:
			continue
		case 1:
			addObjIntoArr(&h.blockAngles, blocker, dist, r1)
		case 2:
			addObjIntoArr(&h.blockAngles, blocker, dist, r1)
			addObjIntoArr(&h.blockAngles, blocker, dist, r2)
		}
	}
	sort.Slice(h.blockAngles, func(i,j int) bool{
		return h.blockAngles[i].Dist<h.blockAngles[j].Dist
	})

	for _, target := range targets {
		if target == ignoreSelf {
			continue
		}
		circle = target.HorizonCircle()
		dist, angles = circle.FromPoint(h.point)
		if dist-circle.Radius > h.maxDist && h.maxDist > 0 {
			continue
		}
		rN, r1, r2 = h.zone.Intersect(angles)
		switch rN {
		case 0:
			continue
		case 1:
			h.applyBlockOnTarget(target, dist, r1)
		case 2:
			h.applyBlockOnTarget(target, dist, r1)
			h.applyBlockOnTarget(target, dist, r2)
		}
	}

	for _, obstacle := range obstacles {
		if obstacle == ignoreSelf {
			continue
		}
		circle = obstacle.HorizonCircle()
		dist, angles = circle.FromPoint(h.point)
		if dist-circle.Radius > h.maxDist && h.maxDist > 0 {
			continue
		}
		rN, r1, r2 = h.zone.Intersect(angles)
		switch rN {
		case 0:
			continue
		case 1:
			h.applyObstacle(obstacle, dist, r1)
		case 2:
			h.applyObstacle(obstacle, dist, r1)
			h.applyObstacle(obstacle, dist, r2)
		}
	}
}

func (h *Horizon) applyBlockOnTarget(target HorizonObject, dist float64, angles vec2.AnglePeriod) {
	h.blockResolve = h.blockResolve[:0]
	h.blockResolve = append(h.blockResolve, angles)
	l:=1//len(h.blockResolve)
	for iBlock :=range h.blockAngles{
		if l==0 || h.blockAngles[iBlock].Dist >= dist {
			break
		}

		i := 0
		for i < l {
			if !h.blockResolve[i].IsIntersect(h.blockAngles[iBlock].Angles){
				i++
				continue
			}
			n, p1, p2 := h.blockResolve[i].Sub(h.blockAngles[iBlock].Angles)
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

	for _, res := range h.blockResolve {
		addObjIntoArr(&h.targetAngles, target, dist, res)
	}
}

func (h *Horizon) applyObstacle(obstacle HorizonObject, dist float64, angles vec2.AnglePeriod) {
	for i := 0; i < len(h.targetAngles); i++ {
		target := h.targetAngles[i]
		if target.Dist <= dist {
			continue
		}
		h.splitTargetWithAngles(i, angles)
	}

	addObjIntoArr(&h.obstacleAngles, obstacle, dist, angles)
}

func (h *Horizon) splitTargetWithAngles(i int, angles vec2.AnglePeriod) {
	target := h.targetAngles[i].Angles
	sN, r1, r2, r3 := target.Split(angles)
	if sN <= 1 {
		return
	}
	h.targetAngles[i].Angles = r1
	if sN >= 2 {
		h.targetAngles = append(h.targetAngles, h.targetAngles[i])
		h.targetAngles[len(h.targetAngles)-1].Angles = r2
	}
	if sN == 3 {
		h.targetAngles = append(h.targetAngles, h.targetAngles[i])
		h.targetAngles[len(h.targetAngles)-1].Angles = r3
	}
}
