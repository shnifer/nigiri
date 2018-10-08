package vec2

import "log"

type AnglePeriod struct{
	start, end float64
}

var EmptyAnglePeriod = AnglePeriod{0,0}
var FullAnglePeriod = AnglePeriod{0,360}

func NewAnglePeriod(start, end float64) AnglePeriod {
	s,e:=NormAngRange(start, end)
	return AnglePeriod{start:s, end:e}
}

func (a AnglePeriod) IsEmpty() bool{
	return a.start == a.end
}
func (a AnglePeriod) IsFull() bool{
	return a.start + 360 == a.end
}
func (a AnglePeriod) Get() (start,end float64){
	return a.start, a.end
}
func (a AnglePeriod) Medium() float64{
	return NormAng((a.start+a.end)/2)
}
func (a AnglePeriod) Wide() float64{
	return a.end-a.start
}

//is dir within AnglePeriod [start-end)
func (a AnglePeriod) has(dir float64) bool{
	dir = NormAng(dir)
	if dir >= a.start && dir<a.end {
		return true
	}
	dir+=360
	if dir >= a.start && dir<a.end {
		return true
	}
	return false
}

func (a AnglePeriod) Intersect (b AnglePeriod) (intersection AnglePeriod, is bool){
	//fixme: implement for two end cross
	if a==FullAnglePeriod{
		return b, true
	}
	if b==FullAnglePeriod{
		return a, true
	}
	if a.has(b.start){
		return NewAnglePeriod(b.start, a.end+360), true
	}
	if b.has(a.start){
		return NewAnglePeriod(a.start, b.end+360), true
	}
	return EmptyAnglePeriod, false
}

func (a AnglePeriod) Sub (b AnglePeriod) (n int, c,d AnglePeriod){
	cross, is:= a.Intersect(b)
	log.Println("a", a)
	log.Println("b", b)
	log.Println("c", cross)
	if !is {
		return 1, a, EmptyAnglePeriod
	}
	if cross == a {
		return 0, EmptyAnglePeriod, EmptyAnglePeriod
	}
	if cross.start == a.start{
		return 1, NewAnglePeriod(cross.end, a.end), EmptyAnglePeriod
	}
	if NormAng(cross.end) == NormAng(a.end) {
		return 1, NewAnglePeriod(a.start, cross.start+360), EmptyAnglePeriod
	}
	return 2, NewAnglePeriod(a.start, cross.start+360),NewAnglePeriod(cross.end, a.end+360)
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

//normalize start angle in [0;360) and end in [start; start+360]
//so always end >= start. End value itself may be more than 360
func NormAngRange(ang1, ang2 float64) (float64, float64) {
	start, end := ang1, ang2
	if start > end {
		start, end = end, start
	}
	d:=end-start
	start = NormAng(start)
	if d == 0 {
		return start, start
	}
	d = NormAng(d)
	if d==0{
		d = 360
	}
	return start, start+d
}