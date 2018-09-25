package main

import (
	"github.com/shnifer/nigiri/vec2"
	"math/rand"
	"strconv"
)

type objectData struct {
	pos     vec2.V2
	radius  float64
	name    string
	r, g, b bool
	texType int
}

var Data []*objectData

func LoadData() {
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
}
