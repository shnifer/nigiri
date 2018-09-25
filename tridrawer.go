package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri/vec2"
	"image"
	"image/color"
)

type TriDrawer struct {
	Src           TriSrcer
	Cam           VTransformer
	Clipper       Clipper
	ChangeableSrc bool

	DrawOptions

	//temporary, do no mem alloc
	v []ebiten.Vertex
}

//Chache for DrawTrianglesOptions.
var dtoCache []*ebiten.DrawTrianglesOptions

//Use another texture for another uid
var triDefTex Tex

func init() {
	dtoCache = make([]*ebiten.DrawTrianglesOptions, 0)
	img, _ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	img.Fill(color.White)
	triDefTex = NewTex(img)
}

func getDto() *ebiten.DrawTrianglesOptions {
	if len(dtoCache) == 0 {
		return new(ebiten.DrawTrianglesOptions)
	}
	v := dtoCache[len(dtoCache)-1]
	dtoCache = dtoCache[:len(dtoCache)-1]
	return v
}
func putDto(dto *ebiten.DrawTrianglesOptions) {
	dtoCache = append(dtoCache, dto)
}

func NewTriDrawer(src TriSrcer, layer Layer, cam VTransformer) *TriDrawer {
	clipper, _ := cam.(Clipper)
	res := &TriDrawer{
		Src:         src,
		Cam:         cam,
		Clipper:     clipper,
		DrawOptions: NewDrawOptions(layer),
		v:           make([]ebiten.Vertex, 0),
	}
	return res
}

func (s *TriDrawer) DrawReqs(Q *Queue) {
	order := reqOrder{
		layer:   s.Layer,
		uid:     triDefTex.uid,
		drawTag: s.drawTag,
	}

	if s.ChangeableSrc {
		v, i, dto := s.getParams(s.Src)
		if len(i) == 0 {
			return
		}
		V := make([]ebiten.Vertex, len(v))
		copy(V, v)
		I := make([]uint16, len(i))
		copy(I, i)
		Q.add(drawReq{
			reqOrder: order,
			f:        s.drawTriF(V, I, dto),
		})
	} else {
		Q.add(drawReq{
			reqOrder: order,
			f:        s.drawSrcF(s.Src),
		})
	}
}

func (s *TriDrawer) getParams(src TriSrcer) ([]ebiten.Vertex,
	[]uint16, *ebiten.DrawTrianglesOptions) {
	dto := getDto()
	*dto = ebiten.DrawTrianglesOptions{
		ColorM:        s.colorM,
		Filter:        s.filter,
		CompositeMode: s.compositeMode,
	}

	v, i := src.GetVerticesIndices()
	if s.Cam == nil {
		//for j:=range v {
		//	v[j].SrcX, v[j].SrcY = 5, 5
		//}
		return v, i, dto
	}

	s.v = append(s.v[:0], v...)
	var p vec2.V2
	for j, V := range s.v {
		//s.v[j].SrcX, s.v[j].SrcY = 5, 5
		p = s.Cam.ApplyV2(vec2.V(float64(V.DstX), float64(V.DstY)))
		s.v[j].DstX, s.v[j].DstY = float32(p.X), float32(p.Y)
	}

	if s.Clipper == nil || len(s.v) == 0 {
		return s.v, i, dto
	}

	clipRect := s.Clipper.ClipRect()
	if clipRect.Empty() {
		return s.v, i, dto
	}

	rect := func(v ebiten.Vertex) image.Rectangle {
		return image.Rect(int(v.DstX), int(v.DstY),
			int(v.DstX)+1, int(v.DstY)+1)
	}

	r := rect(s.v[0])
	for i := 1; i < len(s.v); i++ {
		r = r.Union(rect(s.v[i]))
	}
	//clipped out, use len(i)=0 as marker
	if !clipRect.Overlaps(r) {
		return s.v, []uint16{}, dto
	}
	return s.v, i, dto
}

func (s *TriDrawer) drawTriF(v []ebiten.Vertex, i []uint16, dto *ebiten.DrawTrianglesOptions) drawF {
	return func(dest *ebiten.Image) {
		dest.DrawTriangles(v, i, triDefTex.image, dto)
		putDto(dto)
	}
}

func (s *TriDrawer) drawSrcF(src TriSrcer) drawF {
	return func(dest *ebiten.Image) {
		v, i, dto := s.getParams(src)
		if len(i) > 0 {
			dest.DrawTriangles(v, i, triDefTex.image, dto)
		}
		putDto(dto)
	}
}
