package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/Shnifer/nigiri"
)

var Q *nigiri.Queue
var S *nigiri.Sprite

func mainLoop(win *ebiten.Image) error{
	Q.Clear()
	Q.Add(S)
	Q.Run(win)
	return nil
}

func main(){
	nigiri.StartProfile("sprite")
	defer nigiri.StopProfile("sprite")

	Q = nigiri.NewQueue(nil)
	nigiri.SetTexLoader(nigiri.FileTexLoader("samples"))
	tex, err:=nigiri.GetTex("HUD_Ship.png")
	if err!=nil{
		panic(err)
	}
	src:=nigiri.NewStatic(tex, nil, "tag")
	S = nigiri.NewSprite(src, nigiri.Scaler{ScaleFactor:1})
	S.SetSmooth(true)
	ebiten.Run(mainLoop, 400,400,1,"TEST")
}