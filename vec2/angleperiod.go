package vec2

//start and end angle of period,
//both are supposed to be [0;360)
//[a,a) is a line-reduced ray, if not isFull
//isFull is 360 degree full circle
type AnglePeriod struct {
	start, end float64
	isFull     bool
}

var EmptyAnglePeriod = AnglePeriod{0, 0, false}
var FullAnglePeriod = AnglePeriod{0, 0, true}

func NewAnglePeriod(start, end float64) AnglePeriod {
	return newAnglePeriod(NormAng(start), NormAng(end))
}

func newAnglePeriod(start, end float64) AnglePeriod {
	return AnglePeriod{
		start: start,
		end:   end}
}

func (a AnglePeriod) IsRay() bool {
	return a.start == a.end && !a.isFull
}
func (a AnglePeriod) IsFull() bool {
	return a.isFull
}
func (a AnglePeriod) Get() (start, end float64) {
	return a.start, a.end
}
func (a AnglePeriod) Start() float64 {
	return a.start
}
func (a AnglePeriod) End() float64 {
	return a.end
}
func (a AnglePeriod) Wide() float64 {
	if a.isFull {
		return 360
	}
	if a.start > a.end {
		return 360 - (a.start - a.end)
	}
	return a.end - a.start
}
func (a AnglePeriod) Medium() float64 {
	//if a.isFull {
	//	return 180
	//}
	//if a.IsRay() {
	//	return a.start
	//}
	//if a.start > a.end {
	//	return NormAng((a.start+a.end)/2 + 180)
	//}
	//return (a.start + a.end) / 2
	return NormAng(a.start + a.Wide()/2)
}
func (a AnglePeriod) MedPart(alpha float64) float64{
	return NormAng(a.start+a.Wide()*alpha)
}

//is dir within AnglePeriod [start;end)
//Ray contains only one value of dir == a.start == a.end
//Full contains any direction
//Dir MUST be NORMED
func (a AnglePeriod) Has(dir float64) bool {
	//if a.isFull {
	//	return true
	//}
	//dir = NormAng(dir)
	//if a.IsRay() {
	//	return dir == a.start
	//}
	//if a.start < a.end {
	//	return dir >= a.start && dir < a.end
	//} else {
	//	return dir >= a.start || dir < a.end
	//}
	return a.isFull || a.start <= dir && dir <= a.end ||
		a.start>a.end && (dir>=a.start || dir<a.end)
}

//HasIn is Has without a.start point, so for period it is (start;end)
//Rays have nothing within
//Dir MUST be NORMED
func (a AnglePeriod) HasIn(dir float64) bool {
//	if a.isFull {
//		return true
//	}
//	if a.IsRay() {
//		return false
//	}
//	dir = NormAng(dir)
//	if a.start < a.end {
//		return dir > a.start && dir < a.end
//	} else {
//		return dir > a.start || dir < a.end
//	}
return a.isFull || a.start==a.end || a.start < dir && dir < a.end ||
	a.start>a.end && (dir>a.start || dir<a.end)
}

func (a AnglePeriod) IsIntersect(b AnglePeriod) bool{
	return a.Has(b.start) || b.Has(a.start)
}

//Intersect returns number of intersection (0-2) and their values
//Rays may intersect equal Ray or period containing ray's direction, result is ray
//Periods touching one start-end point do not intersect in it,
//so intersect results on non-ray periods can't be a ray
func (a AnglePeriod) Intersect(b AnglePeriod) (n int, r1, r2 AnglePeriod) {
	if a.isFull {
		return 1, b, EmptyAnglePeriod
	}
	if b.isFull {
		return 1, a, EmptyAnglePeriod
	}
	if a.IsRay() {
		if b.Has(a.start) {
			return 1, a, EmptyAnglePeriod
		} else {
			return 0, EmptyAnglePeriod, EmptyAnglePeriod
		}
	}
	if b.IsRay() {
		if a.Has(b.start) {
			return 1, b, EmptyAnglePeriod
		} else {
			return 0, EmptyAnglePeriod, EmptyAnglePeriod
		}
	}
	if a.Has(b.start) && b.Has(a.start) {
		return 2, newAnglePeriod(b.start, a.end), newAnglePeriod(a.start, b.end)
	}
	if a.Has(b.start) {
		return 1, newAnglePeriod(b.start, a.end), EmptyAnglePeriod
	}
	if b.Has(a.start) {
		return 1, newAnglePeriod(a.start, b.end), EmptyAnglePeriod
	}
	return 0, EmptyAnglePeriod, EmptyAnglePeriod
}

//Sub subtracts b from a period, returning number of, and parts
//Ray subtracted from equal ray deletes it.
//Ray subtracted from period is no-op.
func (a AnglePeriod) Sub(b AnglePeriod) (n int, c, d AnglePeriod) {
	if b.isFull {
		return 0, EmptyAnglePeriod, EmptyAnglePeriod
	}
	if b.IsRay() {
		if a == b {
			return 0, EmptyAnglePeriod, EmptyAnglePeriod
		} else {
			return 1, a, EmptyAnglePeriod
		}
	}
	if a.IsRay() {
		if b.Has(a.start) {
			return 0, EmptyAnglePeriod, EmptyAnglePeriod
		} else {
			return 1, a, EmptyAnglePeriod
		}
	}
	if a.isFull {
		return 1, newAnglePeriod(b.end, b.start), EmptyAnglePeriod
	}

	//both a and b here are periods, not rays or full
	if a.HasIn(b.start) && a.HasIn(b.end) {
		return 2, newAnglePeriod(a.start, b.start), newAnglePeriod(b.end, a.end)
	}
	if a.HasIn(b.end) {
		return 1, newAnglePeriod(b.end, a.end), EmptyAnglePeriod
	}
	if a.HasIn(b.start) {
		return 1, newAnglePeriod(a.start, b.start), EmptyAnglePeriod
	}
	if b.Has(a.start){
		return 0, EmptyAnglePeriod, EmptyAnglePeriod
	}
	return 1, a, EmptyAnglePeriod
}

func (a AnglePeriod) Split(b AnglePeriod) (n int, r1, r2, r3 AnglePeriod) {
	intersectN, i1, i2 := a.Intersect(b)
	SubN, s1, s2 := a.Sub(b)
	n = intersectN + SubN

	if n == 1 {
		return 1, a, EmptyAnglePeriod, EmptyAnglePeriod
	}

	if intersectN == 1 {
		return n, i1, s1, s2
	} else {
		return n, i1, i2, s1
	}
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
	d := end - start
	start = NormAng(start)
	if d == 0 {
		return start, start
	}
	d = NormAng(d)
	if d == 0 {
		d = 360
	}
	return start, start + d
}
