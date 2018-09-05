package nigiri

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"sync"
	"time"
)

type tempTex struct {
	tex  *ebiten.Image
	used bool
	usedInCycle bool
}

const poolShrinkCycles = 100

type texPool struct {
	sync.Mutex
	p       []tempTex
	updateCounter int
	usedCount int
	maxUsed int
}

var ttPool *texPool
func init() {
	ttPool = newTexPool()
}

func GetTempTex(w, h int) *ebiten.Image {
	return ttPool.GetTex(w, h)
}

func PutTempTex(tex *ebiten.Image) {
	ttPool.PutTex(tex)
}

func newTexPool() *texPool {
	res := &texPool{
		p:       make([]tempTex, 0),
	}
	return res
}

func (pool *texPool) GetTempTex(w, h int) *ebiten.Image {
	pool.Lock()
	defer pool.Unlock()

	pool.usedCount++
	if pool.usedCount>pool.maxUsed{
		pool.maxUsed = pool.usedCount
	}

	var sw, sh int
	for i, v := range pool.p {
		if v.used || v.tex == nil {
			continue
		}
		sw, sh = v.tex.Size()
		if sw < w || sh < h {
			continue
		}

		pool.p[i].used = true
		pool.p[i].usedInCycle = true
		return pool.p[i].tex
	}

	tex, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	pool.insertNewTex(tex, true)
	//@@@
	log.Println("pool extended to len: ", len(pool.p))
	return tex
}

func (pool *texPool) PutTempTex(tex *ebiten.Image) {
	pool.Lock()
	defer pool.Unlock()

	pool.usedCount--
	for i := range pool.p {
		if pool.p[i].tex == tex {
			pool.p[i].used = false
			break
		}
	}
}

func (pool *texPool) GetPoolTex(w, h int) *ebiten.Image {
	pool.Lock()
	defer pool.Unlock()

	//todo: binary search
	var sw, sh int
	for i, v := range pool.p {
		if v.used || v.tex == nil {
			continue
		}
		sw, sh = v.tex.Size()
		if sw < w || sh < h {
			continue
		}

		pool.removeElement(i)
		return v.tex
	}

	tex, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	return tex
}

func (pool *texPool) PutPoolTex(tex *ebiten.Image) {
	pool.Lock()
	defer pool.Unlock()

	for i := range pool.p {
		if pool.p[i].tex == tex {
			break
		}
	}
	pool.insertNewTex(tex,false)
	//@@@
	log.Println("pool extended to len: ", len(pool.p))
}


func (pool *texPool) afterLoop(){
	pool.Lock()
	defer pool.Unlock()

	pool.usedCount = 0
	for i:=0; i<len(pool.p);i++{
		pool.p[i].used = false
	}
	pool.updateCounter++
	if pool.updateCounter< poolShrinkCycles {
		return
	}
	pool.updateCounter=0

	pool.checkShrink()

	pool.maxUsed = 0
	for i:=0; i<len(pool.p);i++{
		pool.p[i].usedInCycle = false
	}
}

//run under mutex
func (pool *texPool) checkShrink() {
	l := len(pool.p)

	now := time.Now().Unix()

	for i := 0; i < len(pool.p); {
		if pool.p[i].used || now-pool.p[i].last < poolShrinkCycles {
			i++
			continue
		}
		pool.p = append(pool.p[:i], pool.p[i+1:]...)
	}

	//@@@
	if len(pool.p) < l {
		log.Println("temp pool shrink to len: ", len(pool.p))
	}
}

func (pool *texPool) insertNewTex(tex *ebiten.Image, used bool) {
	var i int
	var sw, sh int
	w,h:=tex.Size()
	//todo: binary search
	for i = 0; i < len(pool.p); i++ {
		sw, sh = pool.p[i].tex.Size()
		if sw+sh >= w+h {
			break
		}
	}
	v := tempTex{
		tex:  tex,
		used: used,
		usedInCycle: true,
	}
	pool.p = append(pool.p, tempTex{})
	copy(pool.p[i+1:], pool.p[i:])
	pool.p[i] = v
}

//run under mutex
func (pool *texPool) removeElement(n int){
	pool.p = append(pool.p[:n], pool.p[n+1:]...)
}