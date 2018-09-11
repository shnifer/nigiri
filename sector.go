package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri/vec2"
	"image/color"
	"math"
)

type Sector struct {
	sprite Sprite

	Center   vec2.V2
	Radius   float64
	StartAng float64
	EndAng   float64

	maxRadius float64
}

const (
	sectorLen      = 1000
	sectorSmallDeg = 1
	sectorMidDeg   = 8
	sectorBigDeg   = 60
)

var smallDegreeTex Tex
var midDegreeTex Tex
var bigDegreeTex Tex

func init() {
	smallDegreeTex = degreeTex(sectorSmallDeg)
	midDegreeTex = degreeTex(sectorMidDeg)
	bigDegreeTex = degreeTex(sectorBigDeg)
}

func degreeTex(ang float64) Tex {
	h := sectorLen
	w := 1 + int(sectorLen*math.Tan(ang*vec2.Deg2Rad))
	img, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)

	p := make([]byte, w*h*4)
	dw := w * 4
	tan := math.Tan(ang * vec2.Deg2Rad)
	var ix, iy int
	white := []byte{255, 255, 255, 255}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if x*x+y*y > sectorLen*sectorLen {
				continue
			}
			if x == 0 || (y > 0 && float64(x)/float64(y) < tan) {
				ix = w - 1 - x
				iy = h - 1 - y
				for i := 0; i < 4; i++ {
					copy(p[4*ix+iy*dw:], white)
				}
			}
		}
	}
	img.ReplacePixels(p)
	return NewTex(img)
}

func NewSector(layer Layer, transformers ...Transformer) Sector {
	sprite := NewSprite(nil, layer, Transforms(transformers))
	sprite.Pivot = vec2.BotRight
	return Sector{
		sprite: sprite,
	}
}

func (s Sector) SetColor(clr color.Color) {
	s.sprite.SetColor(clr)
}
func (s Sector) SetAlpha(a float64) {
	s.sprite.SetAlpha(a)
}
func (s Sector) ColorAlpha() (clr color.Color, a float64) {
	return s.sprite.ColorAlpha()
}

func (s *Sector) ReduceRadiusCam(cam *Camera){
	if cam.clipRect.Empty(){
		s.maxRadius = 0
		return
	}
	sx,sy:=float64(cam.clipRect.Dx()), float64(cam.clipRect.Dy())
	L:=vec2.V(sx,sy).Len()

	s.maxRadius = L / cam.scale
}


func (s Sector) DrawReqs(Q *Queue) {
	s.sprite.Position = s.Center
	scale:=s.Radius / sectorLen
	if s.maxRadius>0 && s.Radius>s.maxRadius{
		scale = s.maxRadius / sectorLen
	}
	s.sprite.Scaler = NewScaler(scale)

	start, end := vec2.NormAngRange(s.StartAng, s.EndAng)

	ang := start
	step := float64(sectorSmallDeg)
	var left float64
loop:
	for {
		left = end - ang
		switch {
		case left >= sectorBigDeg:
			step = sectorBigDeg * 0.9
			s.sprite.Src = bigDegreeTex
		case left >= sectorMidDeg:
			step = sectorMidDeg * 0.95
			s.sprite.Src = midDegreeTex
		case left >= sectorSmallDeg:
			step = sectorSmallDeg * 0.99
			s.sprite.Src = smallDegreeTex
		default:
			break loop
		}
		s.sprite.Angle = ang
		ang += step
		Q.Add(s.sprite)
	}

	lastPart := end - sectorSmallDeg
	if lastPart > start {
		s.sprite.Src = smallDegreeTex
		s.sprite.Angle = lastPart
		Q.Add(s.sprite)
	}
}
