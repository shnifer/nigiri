package angle

//start and end angle of period,
//both are supposed to be [0;360)
//[a,a) is a line-reduced ray, if not isFull
//isFull is 360 degree full circle
type Period struct {
	start, end float64
	isFull     bool
}

var EmptyPeriod = Period{0, 0, false}
var FullPeriod = Period{0, 0, true}

func NewPeriod(start, end float64) Period {
	return newPeriod(Norm(start), Norm(end))
}

func newPeriod(start, end float64) Period {
	return Period{
		start: start,
		end:   end}
}

func (a Period) IsRay() bool {
	return a.start == a.end && !a.isFull
}
func (a Period) IsFull() bool {
	return a.isFull
}
func (a Period) Get() (start, end float64) {
	return a.start, a.end
}
func (a Period) Start() float64 {
	return a.start
}
func (a Period) End() float64 {
	return a.end
}
func (a Period) Wide() float64 {
	if a.isFull {
		return 360
	}
	if a.start > a.end {
		return 360 - (a.start - a.end)
	}
	return a.end - a.start
}
func (a Period) Medium() float64 {
	//if a.isFull {
	//	return 180
	//}
	//if a.IsRay() {
	//	return a.start
	//}
	//if a.start > a.end {
	//	return Norm((a.start+a.end)/2 + 180)
	//}
	//return (a.start + a.end) / 2
	return Norm(a.start + a.Wide()/2)
}
func (a Period) MedPart(alpha float64) float64{
	return Norm(a.start+a.Wide()*alpha)
}

//is dir within Period [start;end)
//Ray contains only one value of dir == a.start == a.end
//Full contains any direction
//Dir MUST be NORMED
func (a Period) Has(dir float64) bool {
	return a.isFull || a.start==dir||a.start <= dir && dir < a.end ||
		a.start>a.end && (dir>=a.start || dir<a.end)
	// this is slower and not inline
	//if a.isFull {
	//	return true
	//}
	//dir = Norm(dir)
	//if a.IsRay() {
	//	return dir == a.start
	//}
	//if a.start < a.end {
	//	return dir >= a.start && dir < a.end
	//} else {
	//	return dir >= a.start || dir < a.end
	//}
}

//HasIn is Has without a.start point, so for period it is (start;end)
//Rays have nothing within
//Dir MUST be NORMED
func (a Period) HasIn(dir float64) bool {
return a.isFull || a.start==a.end && a.start==dir|| a.start < dir && dir < a.end ||
	a.start>a.end && (dir>a.start || dir<a.end)
	// this is slower and not inline
	//	if a.isFull {
	//		return true
	//	}
	//	if a.IsRay() {
	//		return false
	//	}
	//	dir = Norm(dir)
	//	if a.start < a.end {
	//		return dir > a.start && dir < a.end
	//	} else {
	//		return dir > a.start || dir < a.end
	//	}
}

func (a Period) IsIntersect(b Period) bool{
	return a.Has(b.start) || b.Has(a.start)
	//this is much slower somehow
	//as,ae := a.start, a.end
	//bs,be := b.start, b.end
	//return as>ae && (bs>=as || bs<ae || be>as || be<ae) ||
	//	as <= bs && bs < ae ||
	//	as < be && be < ae ||
	//		a.isFull || b.isFull
}

//Intersect returns number of intersection (0-2) and their values
//Rays may intersect equal Ray or period containing ray's direction, result is ray
//Periods touching one start-end point do not intersect in it,
//so intersect results on non-ray periods can't be a ray
func (a Period) Intersect(b Period) (n int, r1, r2 Period) {
	if a.isFull {
		return 1, b, EmptyPeriod
	}
	if b.isFull {
		return 1, a, EmptyPeriod
	}
	if a.IsRay() {
		if b.Has(a.start) {
			return 1, a, EmptyPeriod
		} else {
			return 0, EmptyPeriod, EmptyPeriod
		}
	}
	if b.IsRay() {
		if a.Has(b.start) {
			return 1, b, EmptyPeriod
		} else {
			return 0, EmptyPeriod, EmptyPeriod
		}
	}
	if a.Has(b.start) && b.Has(a.start) {
		return 2, newPeriod(b.start, a.end), newPeriod(a.start, b.end)
	}
	if a.Has(b.start) {
		return 1, newPeriod(b.start, a.end), EmptyPeriod
	}
	if b.Has(a.start) {
		return 1, newPeriod(a.start, b.end), EmptyPeriod
	}
	return 0, EmptyPeriod, EmptyPeriod
}

//Sub subtracts b from a period, returning number of, and parts
//Ray subtracted from equal ray deletes it.
//Ray subtracted from period is no-op.
func (a Period) Sub(b Period) (n int, c, d Period) {
	if b.isFull {
		return 0, EmptyPeriod, EmptyPeriod
	}
	if b.IsRay() {
		if a == b {
			return 0, EmptyPeriod, EmptyPeriod
		} else {
			return 1, a, EmptyPeriod
		}
	}
	if a.IsRay() {
		if b.Has(a.start) {
			return 0, EmptyPeriod, EmptyPeriod
		} else {
			return 1, a, EmptyPeriod
		}
	}
	if a.isFull {
		return 1, newPeriod(b.end, b.start), EmptyPeriod
	}

	//both a and b here are periods, not rays or full
	if a.HasIn(b.start) && a.HasIn(b.end) {
		return 2, newPeriod(a.start, b.start), newPeriod(b.end, a.end)
	}
	if a.HasIn(b.end) {
		return 1, newPeriod(b.end, a.end), EmptyPeriod
	}
	if a.HasIn(b.start) {
		return 1, newPeriod(a.start, b.start), EmptyPeriod
	}
	if b.Has(a.start){
		return 0, EmptyPeriod, EmptyPeriod
	}
	return 1, a, EmptyPeriod
}

func (a Period) Split(b Period) (n int, r1, r2, r3 Period) {
	intersectN, i1, i2 := a.Intersect(b)
	SubN, s1, s2 := a.Sub(b)
	n = intersectN + SubN

	if n == 1 {
		return 1, a, EmptyPeriod, EmptyPeriod
	}

	if intersectN == 1 {
		return n, i1, s1, s2
	} else {
		return n, i1, i2, s1
	}
}

//put angle in degs in [0;360) range
func Norm(angle float64) float64 {
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
func NormRange(ang1, ang2 float64) (float64, float64) {
	start, end := ang1, ang2
	if start > end {
		start, end = end, start
	}
	d := end - start
	start = Norm(start)
	if d == 0 {
		return start, start
	}
	d = Norm(d)
	if d == 0 {
		d = 360
	}
	return start, start + d
}
