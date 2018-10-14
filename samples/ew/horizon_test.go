package ew

import (
	"testing"
	"github.com/shnifer/nigiri/vec2"
	"math/rand"
)

func objectsSetup(count int, objectSize, areaRadius float64) []HorizonObject{
	HorizonObjects := make([]HorizonObject,count)
	//setup
	for i:=0; i<count; i++{
		circle:=Circle{Center: vec2.RandomInCircle(areaRadius), Radius: (1+rand.Float64())*objectSize}
		HorizonObjects[i] = DiffShadowBody{Circle: circle}
	}
	return HorizonObjects
}

func runCalc(b *testing.B, count int, objectSize, areaRadius, zoneHeight float64) {
	HorizonObjects := objectsSetup(count , objectSize, areaRadius)

	horizon:=NewHorizon()
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
