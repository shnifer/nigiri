package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

var rectIndices []uint16
func init(){
	rectIndices = []uint16{0,1,2,2,3,0}
}

type RectDrawer struct{
	Rect
	*TriDrawer
	cam RTransformer
}

func DrawRect(img *ebiten.Image, rect Rect, clr color.Color){
	const MaxColor = 0xffff

	var v []ebiten.Vertex
	corners:=rect.Corners()
	for _,v2:=range corners{
		v = append(v, Vertex(v2))
	}
	dto:=getDto()
	*dto = ebiten.DrawTrianglesOptions{}
	dto.ColorM.Reset()
	r, g, b, a := clr.RGBA()
	dto.ColorM.Scale(float64(r)/MaxColor,
		float64(g)/MaxColor,
		float64(b)/MaxColor,
		float64(a)/MaxColor)

	img.DrawTriangles(v, rectIndices, triDefTex.image, dto)
	putDto(dto)
}

func (rd RectDrawer) GetVerticesIndices() (v []ebiten.Vertex, i []uint16) {
	rect:=rd.Rect
	if rd.cam!=nil{
		rect = rd.cam.TransformRect(rect)
	}
	corners:=rect.Corners()
	for _,v2:=range corners{
		v = append(v, Vertex(v2))
	}
	return v, rectIndices
}

func NewRectDrawer (r Rect, layer Layer, cam RTransformer)RectDrawer{
	res:=RectDrawer{
		Rect: r,
		cam: cam,
	}
	res.TriDrawer = NewTriDrawer(res ,layer, nil)
	return res
}