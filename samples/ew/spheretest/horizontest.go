package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri"
	_ "image/png"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
	"log"
	"github.com/shnifer/nigiri/vec2"
)

var Q *nigiri.Queue
var S *Sphere

func mainLoop(win *ebiten.Image, dt float64) error {
	//if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
	//	x, y := ebiten.CursorPosition()
	//	point := C.UnApplyPoint(vec2.V(float64(x), float64(y)))
	//	horizon.SetPointZoneDist(point, vec2.FullAnglePeriod, 0)
	//}
	if ebiten.IsRunningSlowly() {
		return nil
	}
	Q.Clear()
	Q.Add(S)
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v",ebiten.CurrentFPS()))
	return nil
}

func main() {
	Q = nigiri.NewQueue()

	sprite:=nigiri.NewSprite(nigiri.CircleTex(),0)
	sprite.Scaler = nigiri.NewFixedScaler(20,20)

	S = &Sphere{
		Center: vec2.V(400,400),
		Radius: 200,
		PointSprite:  sprite,
	}

	ebiten.SetVsyncEnabled(false)
	err:=nigiri.Run(mainLoop, 800, 800, 1, "TEST")
	if err!=nil{
		log.Println("ERROR: ",err)
	}
}
