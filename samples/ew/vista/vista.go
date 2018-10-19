package vista

import (
	"github.com/shnifer/nigiri/vec2"
	"github.com/shnifer/nigiri/vec2/angle"
	"sort"
)

type Object interface {
	VistaCircle() Circle
	VistaTypes() (isObstacle, isTarget, isBlocker bool)
}

type ObjectData struct {
	Object Object
	Dist   float64
	Area   Area
}

type Result struct {
	Target     ObjectData
	Blockers   []ObjectData
	Obstacles  []ObjectData
	MainPeriod angle.Period
}

type SightCone struct{
	Position vec2.V2
	Zone Area
	MaxDist float64
}

type Vista struct {
	//current source position
	SightCone
	IgnoreSelf Object

	//temporary arrays
	circleData []circleDat
	targetMainPeriods []angle.Period
	blockAreas []ObjectData
	targetAreas    []ObjectData
	blocksOnTarget [][]ObjectData
	obstacleAreas []ObjectData
	result        []Result
}

type circleDat struct {
	block, target, obstacle bool
	dist                    float64
	period                  angle.Period
}

func New() *Vista {
	return &Vista{
		SightCone: SightCone{
			Zone: NewArea(angle.FullPeriod, 180),
		},
		blockAreas:     make([]ObjectData, 0),
		targetAreas:    make([]ObjectData, 0),
		blocksOnTarget: make([][]ObjectData, 0),
		obstacleAreas:  make([]ObjectData, 0),
		result:         make([]Result, 0),
		circleData:     make([]circleDat, 0),
		targetMainPeriods: make([]angle.Period,0),
	}
}

func (h *Vista) Calculate(objects []Object) {
	var (
		circle       Circle
		height, dist float64
		angles       angle.Period
		rec          Result
	)

	h.ClearTempSlices(true)
	h.circleData = h.circleData[:0]
	for i := range objects {
		o, t, b := objects[i].VistaTypes()
		h.circleData = append(h.circleData, circleDat{
			obstacle: o,
			target:   t,
			block:    b,
		})
	}

	getCircleData := func(ind int) {
		if h.circleData[ind].dist != 0 {
			dist, angles = h.circleData[ind].dist, h.circleData[ind].period
		} else {
			dist, angles = circle.FromPoint(h.Position)
			h.circleData[ind].dist = dist
			h.circleData[ind].period = angles
		}
	}

	for ind, blocker := range objects {
		if !h.circleData[ind].block {
			continue
		}
		if blocker == h.IgnoreSelf{
			continue
		}

		circle = blocker.VistaCircle()
		getCircleData(ind)

		if dist-circle.Radius > h.MaxDist && h.MaxDist > 0 {
			continue
		}
		height = angles.Wide()
		if height > h.Zone.Height {
			height = h.Zone.Height
		}
		if h.Zone.IsIntersect(angles) {
			addObjIntoArr(&h.blockAreas, blocker, dist, NewArea(angles, height))
		}
	}
	sort.Sort(byDist(h.blockAreas))

	h.targetMainPeriods = h.targetMainPeriods[:0]
	blocksOver := make([]ObjectData, 0, 20)
	parts := make([]Area, 0)
mainLoop:
	for ind, target := range objects {
		if !h.circleData[ind].target {
			continue
		}
		if target == h.IgnoreSelf {
			continue
		}
		circle = target.VistaCircle()

		getCircleData(ind)
		if dist-circle.Radius > h.MaxDist && h.MaxDist > 0 {
			continue
		}
		height = angles.Wide()
		if height > h.Zone.Height {
			height = h.Zone.Height
		}
		if !h.Zone.IsIntersect(angles) {
			continue
		}
		targetArea := NewArea(angles, height)
		blocksOver = blocksOver[:0]
		for i := range h.blockAreas {
			if h.blockAreas[i].Dist >= dist {
				break
			}
			if h.blockAreas[i].Area.IsIntersect(targetArea.Period) {
				if h.blockAreas[i].Area.Contains(targetArea) {
					continue mainLoop
				}
				blocksOver = append(blocksOver, h.blockAreas[i])
			}
		}

		n, zones := targetArea.Intersect(h.Zone.Period)
		parts = parts[:0]
		for i := 0; i < n; i++ {
			parts = append(parts, targetArea)
			parts[i].Period = zones[i]
		}

		for _, block := range blocksOver {
			if block.Area.Height < targetArea.Height {
				continue
			}
			i, l := 0, len(parts)
			for i < l {
				n, chunks := parts[i].Sub(block.Area.Period)
				switch n {
				case 0:
					parts[i] = parts[l-1]
					parts = parts[:l-1]
					l--
				case 1:
					parts[i].Period = chunks[0]
					i++
				case 2:
					parts[i].Period = chunks[0]
					parts = append(parts, parts[i])
					parts[l].Period = chunks[1]
					l++
					i++
				}
			}
		}

		for _, part := range parts {
			//add together, cz we synchronize this data by slice index
			addObjIntoArr(&h.targetAreas, target, dist, part)
			if len(blocksOver) == 0 {
				h.blocksOnTarget = append(h.blocksOnTarget, nil)
			} else {
				arr := make([]ObjectData, len(blocksOver))
				copy(arr, blocksOver)
				h.blocksOnTarget = append(h.blocksOnTarget, arr)
			}
			//register primary targetInd
			h.targetMainPeriods = append(h.targetMainPeriods, angles)
		}
	}

	for ind, obstacle := range objects {
		if !h.circleData[ind].obstacle {
			continue
		}
		if obstacle == h.IgnoreSelf {
			continue
		}
		circle = obstacle.VistaCircle()
		getCircleData(ind)

		if dist-circle.Radius > h.MaxDist && h.MaxDist > 0 {
			continue
		}
		height = angles.Wide()
		if height > h.Zone.Height {
			height = h.Zone.Height
		}
		if h.Zone.IsIntersect(angles) {
			addObjIntoArr(&h.obstacleAreas, obstacle, dist, NewArea(angles, height))
		}
	}
	sort.Sort(byDist(h.obstacleAreas))

partLoop:
	for ind, targetData := range h.targetAreas {
		blockers := h.blocksOnTarget[ind]
		partArea := targetData.Area

		var blocksOver []ObjectData
		for _, block := range blockers {
			if block.Area.IsIntersect(partArea.Period) {
				if block.Area.Contains(partArea) {
					continue partLoop
				}
				blocksOver = append(blocksOver, block)
			}
		}

		var obsOver []ObjectData
	obsLoop:
		for i := range h.obstacleAreas {
			if h.obstacleAreas[i].Dist >= targetData.Dist {
				break
			}
			if !h.obstacleAreas[i].Area.IsIntersect(partArea.Period) {
				continue
			}
			for _, block := range blocksOver {
				if block.Area.Contains(h.obstacleAreas[i].Area) {
					continue obsLoop
				}
			}
			obsOver = append(obsOver, h.obstacleAreas[i])
		}

		rec = Result{
			Blockers:   blocksOver,
			Target:     targetData,
			MainPeriod: h.targetMainPeriods[ind],
			Obstacles:  obsOver,
		}
		h.result = append(h.result, rec)
	}
	h.ClearTempSlices(false)
}

func (h *Vista) Result()[]Result{
	return h.result
}

func (h *Vista) sortBlockAreas(i, j int) bool {
	return h.blockAreas[i].Dist < h.blockAreas[j].Dist
}

func addObjIntoArr(arr *[]ObjectData, obj Object, dist float64, area Area) {
	*arr = append(*arr, ObjectData{
		Object: obj,
		Dist:   dist,
		Area:   area,
	})
}

type byDist []ObjectData

func (o byDist) Len() int {
	return len(o)
}
func (o byDist) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
func (o byDist) Less(i, j int) bool {
	return o[i].Dist < o[j].Dist
}
func (h *Vista) ClearTempSlices(resultsToo bool){
	for i:=0; i<len(h.blockAreas); i++{
		h.blockAreas[i].Object = nil
	}
	for i:=0; i<len(h.targetAreas); i++{
		h.targetAreas[i].Object = nil
	}
	for i:=0; i<len(h.blocksOnTarget); i++{
		h.blocksOnTarget[i] = nil
	}
	for i:=0; i<len(h.obstacleAreas); i++{
		h.obstacleAreas[i].Object = nil
	}
	h.blockAreas = h.blockAreas[:0]
	h.targetAreas = h.targetAreas[:0]
	h.blocksOnTarget = h.blocksOnTarget[:0]
	h.obstacleAreas = h.obstacleAreas[:0]

	if resultsToo {
		for i:=0; i<len(h.result); i++{
			h.result[i].Target.Object = nil
			h.result[i].Obstacles = nil
			h.result[i].Blockers = nil
		}
		h.result = h.result[:0]
	}
}