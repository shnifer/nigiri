package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri/vec2"
	"github.com/shnifer/nigiri/vec2/angle"
)

const sectorAngleStep = 5

type SectorParams struct {
	Center   vec2.V2
	Radius   float64
	StartAng float64
	EndAng   float64
}

func (p SectorParams) New(layer Layer, camera *Camera) *Sector {
	return NewSector(p, layer, camera)
}

type Sector struct {
	SectorParams
	*TriDrawer
	calced SectorParams
	v      []ebiten.Vertex
	i      []uint16
}

func NewSector(p SectorParams, layer Layer, vTransformer VTransformer) *Sector {
	res := &Sector{
		SectorParams: p,
		calced:       p,
		v:            make([]ebiten.Vertex, 0),
		i:            make([]uint16, 0),
	}
	res.TriDrawer = NewTriDrawer(res, layer, vTransformer)
	res.recalc()
	return res
}

func (s *Sector) recalc() {
	s.v = s.v[:0]
	s.i = s.i[:0]

	s.v = append(s.v, Vertex(s.Center))
	start, end := angle.NormRange(s.StartAng, s.EndAng)
	s.v = append(s.v, Vertex(s.Center.AddMul(vec2.InDir(start), s.Radius)))
	st_int := int(start)
	st_int = st_int - (st_int % sectorAngleStep) + sectorAngleStep
	for a := float64(st_int); a < end; a += sectorAngleStep {
		s.v = append(s.v, Vertex(s.Center.AddMul(vec2.InDir(a), s.Radius)))
	}
	s.v = append(s.v, Vertex(s.Center.AddMul(vec2.InDir(end), s.Radius)))

	for n := 1; n < len(s.v)-1; n++ {
		s.i = append(s.i, 0, uint16(n), uint16(n+1))
	}
}

func (s *Sector) GetVerticesIndices() (v []ebiten.Vertex, i []uint16) {
	if s.calced != s.SectorParams {
		s.calced = s.SectorParams
		s.recalc()
	}
	return s.v, s.i
}
