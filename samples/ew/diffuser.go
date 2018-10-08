package main

type Diffuser interface {
	HorizonObject
	DiffuserAlbedo(t EmiType) float64
}

type Shadower interface{
	HorizonObject
	ShadowDensity(t EmiType) (density float64)
	ShadowBlock() bool
}

type DiffShadowBody struct{
	Circle
	Albedo float64
}

func (d *DiffShadowBody) HorizonCircle() Circle {
	return d.Circle
}
func (d *DiffShadowBody) ShadowDensity(t EmiType) (density float64) {
	return 0
}
func (d *DiffShadowBody) ShadowBlock() bool{
	return true
}

func (d *DiffShadowBody) DiffuserAlbedo(t EmiType) float64 {
	return d.Albedo
}
