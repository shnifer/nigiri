package main

import (
	"testing"
	"github.com/shnifer/nigiri/vec2"
	"math/rand"
)

func BenchmarkHorizon_Calc1(b *testing.B) {
	const objectCount = 300
	HorizonObjects := make([]HorizonObject,objectCount)
	//setup
	for i:=0; i<objectCount; i++{
		circle:=Circle{Center: vec2.RandomInCircle(800), Radius: rand.Float64()*50+10}
		HorizonObjects[i] = NewSolidObject(circle)
	}

	horizon:=NewHorizon()

	b.ResetTimer()
	for i:=0;i<b.N; i++{
		horizon.SetPointZoneDist(vec2.InDir(float64(i)/float64(b.N)*360).Mul(1000), vec2.FullAnglePeriod, 0)
		hRes:=horizon.Calculate(HorizonObjects, nil, HorizonObjects, nil)
		_ = hRes
	}
}

func BenchmarkHorizon_Calc2(b *testing.B) {
	const objectCount = 300
	HorizonObjects := make([]HorizonObject,objectCount)
	//setup
	for i:=0; i<objectCount; i++{
		circle:=Circle{Center: vec2.RandomInCircle(10000), Radius: rand.Float64()*50+10}
		HorizonObjects[i] = NewSolidObject(circle)
	}

	horizon:=NewHorizon()

	b.ResetTimer()
	for i:=0;i<b.N; i++{
		horizon.SetPointZoneDist(vec2.InDir(float64(i)/float64(b.N)*360).Mul(1000), vec2.FullAnglePeriod, 0)
		hRes:=horizon.Calculate(HorizonObjects, nil, HorizonObjects, nil)
		_ = hRes
	}
}