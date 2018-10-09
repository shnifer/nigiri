package vec2

//start and end angle of period,
//both are supposed to be [0;360)
//[a,a) is a line-reduced ray, if not isFull
//isFull is 360 degree full circle
type AnglePeriod struct{
	start, end float64
	isFull bool
}

var EmptyAnglePeriod = AnglePeriod{0,0, false}
var FullAnglePeriod = AnglePeriod{0,0, true}

func NewAnglePeriod(start, end float64) AnglePeriod {
	return AnglePeriod{
		start:NormAng(start),
		end:NormAng(end)}
}

func (a AnglePeriod) IsEmpty() bool{
	return a.start == a.end && !a.isFull
}
func (a AnglePeriod) IsFull() bool{
	return a.isFull
}
func (a AnglePeriod) Get() (start,end float64){
	return a.start, a.end
}
func (a AnglePeriod) Start() float64{
	return a.start
}
func (a AnglePeriod) End() float64{
	return a.end
}

//is dir within AnglePeriod [start-end)
func (a AnglePeriod) Has(dir float64) bool{
	if a.isFull{
		return true
	}
	dir = NormAng(dir)
	if a.start<a.end{
		return dir>=a.start && dir<a.end
	} else {
		return dir>=a.start || dir<a.end
	}
}

func (a AnglePeriod) Intersect (b AnglePeriod) (intersection AnglePeriod, is bool){
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