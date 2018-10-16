package main

import (
	"testing"
	"math"
	"github.com/shnifer/nigiri/vec2"
)


func sin1(alpha,D,R float64) float64{
	v:=math.Sin(alpha*vec2.Deg2Rad)*D/R
	if v>1 {
		return 90-alpha
	}
	return math.Asin(v)*vec2.Rad2Deg - alpha
}

func Test_sin1(t *testing.T) {
	D:= 10.0
	R:= 4.0
	N:=10

	maxAngle:=math.Asin(R/D)

	for i:=0;i<N;i++{
		angle:=maxAngle*float64(i)/float64(N-1)
		beta:=sin1(angle, D,R)
		dir:=vec2.V(0,-D).AddMul(vec2.InDir(180-beta), R).Dir()
		if math.Abs(dir-angle)>0.0001 {
			t.Errorf("for angle %v != result %v",angle, dir)
		}
	}
}

func Benchmark_sin1(b *testing.B) {
	D:= 10.0
	R:= 4.0

	maxAngle:=math.Asin(R/D)
	step:=maxAngle/float64(b.N-1)
	s:=0.0
	for i:=0;i<b.N;i++{
		angle:=step*float64(i)
		beta:=sin1(angle, D,R)
		s+=beta
	}
}