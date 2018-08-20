package main

import (
	"github.com/hajimehoshi/ebiten"
	"sort"
)

type Queue struct{
	reqs []drawReq
	cam *Camera
}

type drawF func (dest *ebiten.Image)

type reqOrder struct {
	Layer int
	GroupTag string
}

func (r reqOrder) less(s reqOrder) bool{
	if r.Layer!=s.Layer{
		return r.Layer<s.Layer
	} else {
		return r.GroupTag<s.GroupTag
	}
}

type drawReq struct{
	reqOrder
	f drawF
}

func (Q *Queue) less(i,j int) bool{
	return Q.reqs[i].less(Q.reqs[j].reqOrder)
}

func NewQueue(cam *Camera) *Queue {
	return &Queue{
		reqs: make([]drawReq ,0),
		cam: cam,
	}
}
func (Q *Queue) Run(dest *ebiten.Image){
	sort.SliceStable(Q.reqs, Q.less)
	for _,req:=range Q.reqs{
		req.f(dest)
	}
}
func (Q *Queue) Clear(){
	Q.reqs = Q.reqs[:0]
}
func (Q *Queue) SetCam(cam *Camera) {
	Q.cam = cam
}
func (Q *Queue) Cam() *Camera{
	return Q.cam
}

type DrawRequester interface {
	DrawReqs(Q *Queue)
}

//for use from outside package
func (Q *Queue) Add(d DrawRequester){
	d.DrawReqs(Q)
}

//for use from primitives, sprites and texts
func (Q *Queue) add(f drawF ,order reqOrder){
	Q.reqs = append(Q.reqs, drawReq{
		f:f,
		reqOrder: order,
	})
}
