package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"sort"
)

type Queue struct {
	reqs []drawReq
}

type drawF func(dest *ebiten.Image)

type Layer float32

type reqOrder struct {
	layer    Layer
	uid uint64
}

func (r reqOrder) less(s reqOrder) bool {
	if r.layer != s.layer {
		return r.layer < s.layer
	} else {
		return r.uid < s.uid
	}
}

type drawReq struct {
	reqOrder
	f drawF
}

func (Q *Queue) less(i, j int) bool {
	return Q.reqs[i].less(Q.reqs[j].reqOrder)
}

func NewQueue() *Queue {
	return &Queue{
		reqs: make([]drawReq, 0),
	}
}

func (Q *Queue) Run(dest *ebiten.Image) {
	sort.SliceStable(Q.reqs, Q.less)
	for _, req := range Q.reqs {
		req.f(dest)
	}
}

func (Q *Queue) Clear() {
	Q.reqs = Q.reqs[:0]
}

//for use from outside package
func (Q *Queue) Add(d DrawRequester) {
	d.DrawReqs(Q)
}

//for use from primitives, sprites and texts
func (Q *Queue) add(drawReq drawReq) {
	Q.reqs = append(Q.reqs, drawReq)
}

func (Q *Queue) Len() int {
	return len(Q.reqs)
}
