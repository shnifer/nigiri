package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/nigiri"
	_ "image/png"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
	"log"
	"github.com/shnifer/nigiri/vec2"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/shnifer/nigiri/samples/ew/vista"
	"github.com/shnifer/nigiri/samples/ew/vista/vistautils"
	"golang.org/x/image/colornames"
)

var Q *nigiri.Queue
var SolidObjects []*SolidObject
var Vista *vista.Vista
var C *nigiri.Camera
var Objects[]vista.Object
var ResSprite *VistaResultsSprite
var ViewDrawer *vistautils.ViewSectorDrawer

func mainLoop(win *ebiten.Image, dt float64) error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		point := C.UnApplyPoint(vec2.V(float64(x), float64(y)))
		Vista.Point = point
		ViewDrawer.Point = Vista.Point
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}
	res:=Vista.Calculate(Objects, nil)
	ResSprite.Take(res)
	Q.Clear()
	for _,rec:=range res{
		ViewDrawer.Target = rec.Target
		Q.Add(ViewDrawer)
	}
	for _,obj:=range SolidObjects{
		Q.Add(obj)
	}
	Q.Add(ResSprite)
	Q.Run(win)
	ebitenutil.DebugPrint(win, fmt.Sprintf("FPS: %v",ebiten.CurrentFPS()))
	return nil
}

func main() {
	Vista =vista.New()
	C = nigiri.NewCamera()
	C.SetCenter(vec2.V(400,400))

	Q = nigiri.NewQueue()

	ViewDrawer = vistautils.NewViewSectorDrawer(-1,C)
	ViewDrawer.Color = colornames.Yellow

	ResSprite = NewVistaResultsSprite(2,10, 20, 2, C)
	ResSprite.Pivot = vec2.BotLeft
	ResSprite.Position = vec2.V(40,760)
	ResSprite.SetAlpha(0.8)

	SolidObjects = make([]*SolidObject,0)
	Objects = make([]vista.Object,0)

	body1:= NewSolidObject(vista.Circle{Center: vec2.ZV, Radius: 50})
	body2:= NewSolidObject(vista.Circle{Center: vec2.V(-100,0), Radius: 30})

	Objects = append(Objects, body1)
	SolidObjects = append(SolidObjects, body1)
	Objects = append(Objects, body2)
	SolidObjects = append(SolidObjects, body2)

	ebiten.SetVsyncEnabled(false)
	err:=nigiri.Run(mainLoop, 800, 800, 1, "TEST")
	if err!=nil{
		log.Println("ERROR: ",err)
	}
}
