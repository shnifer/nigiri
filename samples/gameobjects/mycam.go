package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri"
	"github.com/shnifer/nigiri/vec2"
)

type MyCam struct {
	*nigiri.Camera
}

func (C MyCam) Update(dt float64) {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		C.Translate(vec2.InDir(C.Angle()).Rotate90().Mul(1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		C.Translate(vec2.InDir(C.Angle()).Rotate90().Mul(-1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		C.Translate(vec2.InDir(C.Angle()).Mul(1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		C.Translate(vec2.InDir(C.Angle()).Mul(-1 / C.Scale()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		C.Rotate(1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		C.Rotate(-1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		C.MulScale(1.05)
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		C.MulScale(1 / 1.05)
	}
}
