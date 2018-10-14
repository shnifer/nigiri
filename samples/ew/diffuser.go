package ew

import "github.com/shnifer/nigiri/samples/ew/vista"

type Diffuser interface {
	vista.Object
	DiffuserAlbedo(t EmiType) float64
}

type Shadower interface {
	vista.Object
	ShadowDensity(t EmiType) (density float64)
	ShadowBlock() bool
}

type DiffShadowBody struct {
	vista.Circle
	Albedo float64
}

func (d DiffShadowBody) VistaTypes() (isObstacle, isTarget, isBlocker bool) {
	return false, true, true
}

func (d DiffShadowBody) VistaCircle() vista.Circle {
	return d.Circle
}
func (d DiffShadowBody) ShadowDensity(t EmiType) (density float64) {
	return 0
}
func (d DiffShadowBody) ShadowBlock() bool {
	return true
}

func (d DiffShadowBody) DiffuserAlbedo(t EmiType) float64 {
	return d.Albedo
}
