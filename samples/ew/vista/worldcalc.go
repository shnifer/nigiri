package vista

import "sync"

func (w *World) Calculate(){
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

	wg:=&sync.WaitGroup{}
	for i:=0; i<vInd; i++{
		wg.Add(1)
		go func (i int) {
			defer wg.Done()
			w.vistas[i].Calculate(w.objects)
		}(i)
	}
	wg.Wait()


}
