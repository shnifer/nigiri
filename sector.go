package nigiri

/*import (
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
}

const (
	defSectorLen      = 1000
	sectorSmallDeg = 0.716 //tan(7.125)=1/80
	sectorMidDeg   = 7.125 //tan(7.125)=1/8
	sectorBigDeg   = 45
)

var sectorLen int

var smallDegreeTex Tex
var midDegreeTex Tex
var bigDegreeTex Tex

func init() {
	SetSectorLen(defSectorLen)
}

func SetSectorLen(len int) {
	sectorLen = len
	smallDegreeTex = degreeTex(sectorSmallDeg)
	midDegreeTex = degreeTex(sectorMidDeg)
	bigDegreeTex = degreeTex(sectorBigDeg)
}

func degreeTex(ang float64) Tex {
	h := sectorLen
	w := 1 + int(float64(sectorLen)*math.Tan(ang*vec2.Deg2Rad))
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

func (s Sector) DrawReqs(Q *Queue) {
	s.sprite.Position = s.Center
	s.sprite.Scaler = NewScaler(s.Radius / float64(sectorLen))

	start, end := vec2.NormAngRange(s.StartAng, s.EndAng)

	var size, step float64
	left:=end-start
	switch {
	case left >= sectorBigDeg:
		size = sectorBigDeg
		step = size * 0.9
		s.sprite.Src = bigDegreeTex
	case left >= sectorMidDeg:
		size = sectorMidDeg
		step = size * 0.95
		s.sprite.Src = midDegreeTex
	default:
		size = sectorSmallDeg
		step = size * 0.99
		s.sprite.Src = smallDegreeTex
	}

	ang := start
	for ang+size<=end{
		s.sprite.Angle = ang
		ang += step
		Q.Add(s.sprite)
	}

	if left>sectorSmallDeg {
		s.sprite.ScaleFactor.X *= (-1)
		s.sprite.Angle = end
		Q.Add(s.sprite)
	}
}
*/
