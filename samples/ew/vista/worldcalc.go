package vista

import (
	"runtime"
	"sync"
)

var defWorkerPoolSize = runtime.NumCPU()
type WorkerPool struct{
	workers chan struct{}
	wg *sync.WaitGroup
}
func NewWorkerPool()WorkerPool{
	return WorkerPool{
		workers: make(chan struct{}, defWorkerPoolSize),
		wg: &sync.WaitGroup{},
	}
}
func (wp WorkerPool) Run(f func()){
	wp.workers <- struct{}{}
	wp.wg.Add(1)
	go func(){
		defer func(){<-wp.workers
		wp.wg.Done()}()

		f()
	}()
}
func (wp WorkerPool) Wait(){
	wp.wg.Wait()
}
func (wp WorkerPool) WaitClose(){
	wp.wg.Wait()
	close (wp.workers)
}

func (w *World) Calculate(){
	w.results = make(map[Watcher]WatchResult, len(w.watchers))

	vInd:=0
	for i:=range w.emitters{
		cone:=w.emitters[i].EmitterCone()
		w.vistas[vInd].SightCone = cone
		obj, _:=w.emitters[i].(Object)
		w.vistas[vInd].IgnoreSelf = obj
		vInd++
	}
	for i:=range w.watchers{
		cone:=w.watchers[i].WatchCone()
		w.vistas[vInd].SightCone = cone
		obj, _:=w.emitters[i].(Object)
		w.vistas[vInd].IgnoreSelf = obj
		vInd++
	}

	wp:=NewWorkerPool()
	for i:=0; i<vInd; i++{
		vista:=w.vistas[i]
		wp.Run(func(){
			vista.Calculate(w.objects)
		})
	}
	wp.Wait()


	//todo: CONTINUE WORK HERE

	wp.Wait()
}