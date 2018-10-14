package vista

import (
	"testing"
	"github.com/shnifer/nigiri/vec2"
	"math/rand"
)

type testRB struct{
	Circle
}

func (b testRB) VistaTypes() (isObstacle, isTarget, isBlocker bool) {
	return false, true, true
}

func (b testRB) VistaCircle() Circle {
	return b.Circle
}
type testCloud struct{
	Circle
}

func (c testCloud) VistaTypes() (isObstacle, isTarget, isBlocker bool) {
	return true, false, false
}

func (c testCloud) VistaCircle() Circle {
	return c.Circle
}

func objectsSetup(count int, objectSize, areaRadius float64, obsCount int, obsSize float64) []Object {
	HorizonObjects := make([]Object,0)
	//setup
	for i:=0; i<count; i++{
		circle:=Circle{Center: vec2.RandomInCircle(areaRadius), Radius: (1+rand.Float64())*objectSize}
		HorizonObjects = append(HorizonObjects, testRB{Circle: circle})
	}
	for i:=0; i<obsCount; i++{
		circle:=Circle{Center: vec2.RandomInCircle(areaRadius), Radius: (1+rand.Float64())*obsSize}
		HorizonObjects = append(HorizonObjects, testCloud{Circle: circle})
	}
	return HorizonObjects
}

func runCalc(b *testing.B, count int, objectSize, areaRadius, zoneHeight float64) {
	HorizonObjects := objectsSetup(count , objectSize, areaRadius, count/5, objectSize*2)

	horizon:= New()
	horizon.Zone.Height = zoneHeight

	var ang [100]vec2.V2
	for i:=0;i<100; i++{
		ang[i] = vec2.InDir(float64(i)/float64(b.N)*360).Mul(areaRadius*0.6)
	}

	b.ResetTimer()
	for i:=0;i<b.N; i++{
		horizon.Point = ang[i%100]
		res := horizon.Calculate(HorizonObjects,  nil)
		_=res
	}
}

func BenchmarkHorizon2_Calc1(b *testing.B) {
	runCalc(b, 300, 20, 300, 180)
}

func BenchmarkHorizon2_Calc3(b *testing.B) {
	runCalc(b, 300, 20, 300, 1)
}

func BenchmarkHorizon2_Calc4(b *testing.B) {
	runCalc(b, 300, 20, 300, 0.01)
}

func BenchmarkHorizon2_Calc5(b *testing.B) {
	runCalc(b, 100, 20, 300, 180)
}

func BenchmarkHorizon2_Calc2(b *testing.B) {
	runCalc(b, 300, 20, 1000, 180)
}
