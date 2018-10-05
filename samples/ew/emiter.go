package main

import "github.com/shnifer/nigiri/vec2"

type EmiType int
const(
	EMI_HEAT EmiType = iota
	EMI_LIGHT
	EMI_HYPERLINK
)

const EmiDirCount = 360

type EmitData struct{
	Type EmiType
	Power float64
	Signature string
}

type HorizonCircler interface {
	HorizonCircle() Circle
}

type Emitter interface {
	HorizonCircler
	Emits(dir float64) []EmitData
}

type LightEmitter struct{
	Circle
	Dir float64
	Signature string
	MaxPower float64
	PowerK [EmiDirCount]float64
}

func NewLightEmitter(maxPower float64, powerK [EmiDirCount]float64, signature string) *LightEmitter{
	res := &LightEmitter{
		Dir:0,
		Signature: signature,
		MaxPower: maxPower,
		PowerK: powerK,
	}
	return res
}

func (l *LightEmitter) HorizonCircle() Circle{
	return l.Circle
}

func (l *LightEmitter) Emits(dir float64) []EmitData {
	if l.MaxPower<=0{
		return nil
	}
	ang:=vec2.NormAng(dir - l.Dir)
	angN:=int(ang*EmiDirCount/360)
	k:=l.PowerK[angN]
	if k<=0{
		return nil
	}
	return []EmitData{
		EmitData{
			Type: EMI_LIGHT,
			Signature: l.Signature,
			Power: l.MaxPower * k,
		},
	}
}