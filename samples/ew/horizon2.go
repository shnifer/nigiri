package ew

import (
	"github.com/shnifer/nigiri/samples/ew/area"
	"github.com/shnifer/nigiri/vec2"
	"github.com/shnifer/nigiri/vec2/angle"
	"sort"
)

type ObjectData struct {
	Object HorizonObject
	Dist   float64
	Area   area.Area
}

type Horizon struct {
	//current source position
	Point   vec2.V2
	Zone    area.Area
	MaxDist float64

	//temporary arrays
	blockAreas []ObjectData
	targetAreas []ObjectData
	blocksOnTarget [][]ObjectData
}

func NewHorizon() *Horizon{
	return &Horizon{
		Zone: area.New(angle.FullPeriod, 180),
		blockAreas: make([]ObjectData,0),
		targetAreas: make([]ObjectData,0),
		blocksOnTarget:make([][]ObjectData,0),
	}
}

func (h *Horizon) CalculateTargets(targets, blockers []HorizonObject, ignoreSelf HorizonObject) {
	var (
		circle         Circle
		height, dist   float64
		angles angle.Period
	)

	h.blockAreas = h.blockAreas[:0]
	for _, blocker := range blockers {
		if blocker == ignoreSelf {
			continue
		}
		circle = blocker.HorizonCircle()
		dist, angles = circle.FromPoint(h.Point)

		if dist-circle.Radius > h.MaxDist && h.MaxDist > 0 {
			continue
		}
		height = angles.Wide()
		if height > h.Zone.Height {
			height = h.Zone.Height
		}
		if h.Zone.IsIntersect(angles) {
			addObjIntoArr(&h.blockAreas, blocker, dist, area.New(angles, height))
		}
	}
	v:=byDist(h.blockAreas)
	sort.Sort(v)

	blocksOver:=make([]ObjectData,0,20)

	h.targetAreas = h.targetAreas[:0]
	h.blocksOnTarget = h.blocksOnTarget[:0]
	mainLoop:
	for _, target := range targets {
		if target == ignoreSelf {
			continue
		}
		circle = target.HorizonCircle()
		dist, angles = circle.FromPoint(h.Point)
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
		targetArea :=area.New(angles, height)
		blocksOver = blocksOver[:0]
		for i:=range h.blockAreas{
			if h.blockAreas[i].Dist>=dist{
				break
			}
			if h.blockAreas[i].Area.IsIntersect(targetArea.Period){
				if h.blockAreas[i].Area.Contains(targetArea) {
					continue mainLoop
				}
				blocksOver = append(blocksOver, h.blockAreas[i])
			}
		}

		//add together, cz we synchronize this data by slice index
		addObjIntoArr(&h.targetAreas, target, dist, targetArea)
		if len(blocksOver)==0 {
			h.blocksOnTarget = append(h.blocksOnTarget, nil)
		} else {
			arr:=make([]ObjectData, len(blocksOver))
			copy(arr, blocksOver)
			h.blocksOnTarget = append(h.blocksOnTarget, arr)
		}
	}
}

func (h *Horizon) sortBlockAreas(i, j int) bool {
	return h.blockAreas[i].Dist < h.blockAreas[j].Dist
}

func addObjIntoArr(arr *[]ObjectData, obj HorizonObject, dist float64, area area.Area) {
	*arr = append(*arr, ObjectData{
		Object: obj,
		Dist:   dist,
		Area: area,
	})
}

type byDist []ObjectData
func (o byDist) Len() int {
	return len(o)
}
func (o byDist) Swap(i,j int) {
	o[i],o[j] = o[j],o[i]
}
func (o byDist) Less(i,j int) bool{
	return o[i].Dist<o[j].Dist
}