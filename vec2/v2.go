package vec2

import (
	"fmt"
	"math"
	"math/rand"
)

//float64 vector2
type V2 struct {
	X float64
	Y float64
}

const Deg2Rad = math.Pi / 180
const Rad2Deg = 180 / math.Pi

//Generators

//zero vector
var ZV V2

//Corners relative vectors
var TopLeft V2
var TopMid V2
var TopRight V2
var MidLeft V2
var Center V2
var MidRight V2
var BotLeft V2
var BotMid V2
var BotRight V2

func init() {
	TopLeft = V2{X: 0, Y: 0}
	TopMid = V2{X: 0.5, Y: 0}
	TopRight = V2{X: 1, Y: 0}
	MidLeft = V2{X: 0, Y: 0.5}
	Center = V2{X: 0.5, Y: 0.5}
	MidRight = V2{X: 1, Y: 0.5}
	BotLeft = V2{X: 0, Y: 1}
	BotMid = V2{X: 0.5, Y: 1}
	BotRight = V2{X: 1, Y: 1}
}

//RandomOrt returns a random vector with len = 1
func RandomOrt() V2 {
	a := rand.Float64() * 2 * math.Pi
	return V2{math.Sin(a), math.Cos(a)}
}

//RandomInCircle returns a random vector in circle with radius R
func RandomInCircle(R float64) V2 {
	if R == 0 {
		return V2{}
	}
	ort := RandomOrt()
	dist := math.Sqrt(rand.Float64() * (R * R))
	return Mul(ort, dist)
}

//InDir return an ort vector in direction of angle degrees
//0 angle is up (0,-1), positive direction is counterclockwise
//for world coords primary, use for screen coords with caution (because of Y axis)
func InDir(angle float64) V2 {
	a := angle * Deg2Rad
	s, c := math.Sincos(a)
	return V2{X: -s, Y: -c}
}

//Operations

//procedure syntax

//AddMul returns a new vector = a+b*t
func AddMul(a, b V2, t float64) V2 {
	return Add(a, Mul(b, t))
}

//Rotate returns a new vector equal to V rotated by angle degrees
func Rotate(V V2, angle float64) V2 {
	if angle == 0 {
		return V
	}
	a := angle * Deg2Rad
	sin, cos := math.Sincos(a)
	return V2{
		X: V.X*cos + V.Y*sin,
		Y: V.Y*cos - V.X*sin,
	}
}

//Rotate returns a new vector equal to V rotated by 90 degrees
func Rotate90(a V2) V2 {
	return V2{
		X: a.Y,
		Y: -a.X,
	}
}

//ApplyOnTransform translate vector V moving by pos and turning by angle degrees
func ApplyOnTransform(V, pos V2, angle float64) V2 {
	return Add(pos, Rotate(V, angle))
}

//Add returns vector equal to sum a+b
func Add(a, b V2) V2 {
	return V2{a.X + b.X, a.Y + b.Y}
}

//Sub returns vector equal to a-b
func Sub(a, b V2) V2 {
	return V2{a.X - b.X, a.Y - b.Y}
}

//Mul returns vector a multiplied by t
func Mul(a V2, t float64) V2 {
	return V2{a.X * t, a.Y * t}
}

//Mul returns vector a.x*b.x; a.y*b.y
func MulXY(a V2, b V2) V2 {
	return V2{a.X * b.X, a.Y * b.Y}
}

//Len return length of vector a
func Len(a V2) float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

//LenSqr return square length of vector a
func LenSqr(a V2) float64 {
	return a.X*a.X + a.Y*a.Y
}

//Normed returns normed copy of vector a with len = 1
func Normed(a V2) V2 {
	if a.X == 0 && a.Y == 0 {
		return a
	}
	K := 1 / Len(a)
	return Mul(a, K)
}

func Dir(v V2) float64 {
	if v == ZV {
		return 0
	}
	if v.Y == 0 {
		if v.X > 0 {
			return 270
		} else {
			return 90
		}
	}
	a := math.Atan(v.X/v.Y) * Rad2Deg
	if v.Y > 0 {
		a += 180
	}
	if a < 0 {
		a += 360
	}
	return a
}

//method syntax

func (a V2) Add(b V2) V2 {
	return Add(a, b)
}

func (a V2) Sub(b V2) V2 {
	return Sub(a, b)
}

func (a V2) Mul(t float64) V2 {
	return Mul(a, t)
}

func (a V2) MulXY(b V2) V2 {
	return MulXY(a, b)
}

func (a V2) Len() float64 {
	return Len(a)
}

func (a V2) LenSqr() float64 {
	return LenSqr(a)
}

func (a V2) Normed() V2 {
	return Normed(a)
}

func (a *V2) DoNorm() {
	*a = Normed(*a)
}

func (a V2) Rotate(angle float64) V2 {
	return Rotate(a, angle)
}

func (a V2) Rotate90() V2 {
	return Rotate90(a)
}

func (v V2) ApplyOnTransform(pos V2, angle float64) V2 {
	return ApplyOnTransform(v, pos, angle)
}

func (a V2) AddMul(b V2, t float64) V2 {
	return AddMul(a, b, t)
}

func (a *V2) DoAddMul(b V2, t float64) {
	*a = AddMul(*a, b, t)
}

func (a V2) Dir() float64 {
	return Dir(a)
}

func (a V2) String() string {
	return fmt.Sprintf("[%v, %v]", a.X, a.Y)
}

func RotateF(ang float64) func(V2) V2 {
	var sin, cos float64
	if ang == 0 {
		sin, cos = 0, 1
	} else {
		sin, cos = math.Sincos(ang * Deg2Rad)
	}
	return func(v V2) V2 {
		return V2{
			X: v.X*cos - v.Y*sin,
			Y: v.Y*cos + v.X*sin,
		}
	}
}

func V(x, y float64) V2 {
	return V2{X: x, Y: y}
}

//put angle in degs in [0;360) range
func NormAng(angle float64) float64 {
	if angle < 0 {
		a := float64(int(-angle/360) + 1)
		return angle + 360*a
	}
	if angle >= 360 {
		a := float64(int(angle / 360))
		return angle - 360*a
	}
	return angle
}

//normalize start angle in [0;360) and end in [start; start+360)
//so always end > start. End value itself may be more than 360
func NormAngRange(start, end float64) (float64, float64) {
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		a := float64(int(-start/360) + 1)
		start += 360 * a
	}
	if start >= 360 {
		a := float64(int(start / 360))
		start -= 360 * a
	}
	if end < start {
		a := float64(int((start-end)/360) + 1)
		end += 360 * a
	}
	if end >= start+360 {
		a := float64(int((end - start) / 360))
		end -= 360 * a
	}
	return start, end
}
