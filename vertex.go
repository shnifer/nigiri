package nigiri

import (
	"github.com/shnifer/nigiri/vec2"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

func Vertex(v2 vec2.V2) ebiten.Vertex {
	return ebiten.Vertex{
		DstX:   float32(v2.X),
		DstY:   float32(v2.Y),
		ColorR: 1,
		ColorG: 1,
		ColorB: 1,
		ColorA: 1,
	}
}
func VertexColor(v2 vec2.V2, clr color.Color) ebiten.Vertex {
	r,g,b,a:=clr.RGBA()
	max:=float32(0xFFFF)
	return ebiten.Vertex{
		DstX:   float32(v2.X),
		DstY:   float32(v2.Y),
		ColorR: float32(r)/max,
		ColorG: float32(g)/max,
		ColorB: float32(b)/max,
		ColorA: float32(a)/max,
	}
}
