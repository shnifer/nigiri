package main

type Diffuser interface {
	HorizonCircler
	DiffuserAlbedo(t EmiType) float64
}

type Shadower interface{
	HorizonCircler
	ShadowDensity(t EmiType) (density float64, absolute bool)
}

type DiffShadowBody struct{
	Circle
	Albedo float64
}

func (d *DiffShadowBody) HorizonCircle() Circle {
	return d.Circle
}
func (d *DiffShadowBody) ShadowDensity(t EmiType) (density float64, absolute bool) {
	return 0, true
}
func (d *DiffShadowBody) DiffuserAlbedo(t EmiType) float64 {
	return d.Albedo
}
