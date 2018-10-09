package main

import (
	"github.com/shnifer/nigiri/vec2"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type objectData struct {
	pos     vec2.V2
	radius  float64
	name    string
	r, g, b bool
	texType int
}

var DataMu sync.Mutex
var Data []*objectData

func StartData() {
	Data = make([]*objectData, 0)

	for i := 0; i < 16; i++ {
		objDat := &objectData{
			pos:     vec2.RandomInCircle(300),
			name:    strconv.Itoa(i),
			radius:  rand.Float64()*20 + 30,
			r:       i&1 > 0,
			g:       i&2 > 0,
			b:       i%4 > 0,
			texType: (i & 8) / 8,
		}
		Data = append(Data, objDat)
	}

	go DataProgress()
}

func DataProgress() {
	var n int
	for {
		time.Sleep(time.Second / 30)
		n++
		refreshData(n)
	}
}

func refreshData(n int) {
	DataMu.Lock()
	defer DataMu.Unlock()

	for i, v := range Data {
		dist := v.pos.Len()
		Data[i].pos = v.pos.Rotate(10 / dist)
		if n%20 == 0 {
			Data[i].r = !v.r
		}
		if n%30 == 0 {
			Data[i].g = !v.g
		}
		if n%50 == 0 {
			Data[i].b = !v.b
		}
	}
}
