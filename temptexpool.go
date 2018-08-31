package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"sync"
	"time"
)

type tempTex struct {
	tex  *ebiten.Image
	last int64
	used bool
}

type tempTexPool struct {
	sync.Mutex
	p       []tempTex
	periodS int64
}

var ttPool *tempTexPool

func newTempTexPool() *tempTexPool {
	res := &tempTexPool{
		p:       make([]tempTex, 0),
		periodS: 5,
	}
	go func() {
		for range time.Tick(time.Second) {
			res.checkLast()
		}
	}()
	return res
}

func (pool *tempTexPool) checkLast() {
	pool.Lock()
	defer pool.Unlock()

	now := time.Now().Unix()

	for i := 0; i < len(pool.p); {
		if pool.p[i].used || now-pool.p[i].last < pool.periodS {
			i++
			continue
		}
		pool.p = append(pool.p[:i], pool.p[i+1:]...)
	}
}

func init() {
	ttPool = newTempTexPool()
}

func GetTempTex(w, h int) *ebiten.Image {
	return ttPool.GetTex(w, h)
}

func PutTempTex(tex *ebiten.Image) {
	ttPool.PutTex(tex)
}

func (pool *tempTexPool) GetTex(w, h int) *ebiten.Image {
	pool.Lock()
	defer pool.Unlock()

	var sw, sh int
	for i, v := range pool.p {
		if v.used || v.tex == nil {
			continue
		}
		sw, sh = v.tex.Size()
		if sw < w || sh < h {
			continue
		}
		if (sw + sh) > 2*(w+h) {
			break
		}

		pool.p[i].used = true
		return pool.p[i].tex
	}

	tex, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	var i int
	for i = 0; i < len(pool.p); i++ {
		sw, sh = pool.p[i].tex.Size()
		if sw+sh >= w+h {
			break
		}
	}
	v := tempTex{
		tex:  tex,
		last: time.Now().Unix(),
		used: true,
	}
	pool.p = append(pool.p, tempTex{})
	copy(pool.p[i+1:], pool.p[i:])
	pool.p[i] = v
	//@@@
	log.Println("temp pool len: ", len(pool.p))
	return tex
}

func (pool *tempTexPool) PutTex(tex *ebiten.Image) {
	pool.Lock()
	defer pool.Unlock()

	for i, v := range pool.p {
		if v.tex == tex {
			pool.p[i].used = false
			pool.p[i].last = time.Now().Unix()
			break
		}
	}
}
