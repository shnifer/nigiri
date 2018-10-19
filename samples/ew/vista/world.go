package vista

type World struct{
	objects []Object
	emitters []Emitter
	watchers []Watcher
	vistas []*Vista

	exists map[interface{}]struct{}
}

type Spectre = string

type Emitter interface {
	EmitterCone() SightCone
	EmitSpectre() Spectre
	EmitDir(dir float64)
}

type Watcher interface {
	WatchCone() SightCone
	WantedSpectre() Spectre
}

func NewWorld() *World{
	return &World{
		objects: make([]Object,0),
		emitters: make([]Emitter,0),
		watchers: make([]Watcher,0),
		vistas: make([]*Vista,0),
		exists: make(map[interface{}]struct{}),
	}
}

func (w *World) Add(x interface{}){
	if _,exist:=w.exists[x]; exist{
		return
	}
	w.exists[x] = struct{}{}

	if obj, ok:=x.(Object);ok{
		w.addObject(obj)
	}
	if emi, ok:=x.(Emitter);ok{
		w.addEmitter(emi)
	}
	if wch, ok:=x.(Watcher);ok{
		w.addWatcher(wch)
	}
}

func (w *World) Del(x interface{}){
	if _,exist:=w.exists[x]; !exist{
		return
	}
	delete(w.exists, x)

	if obj, ok:=x.(Object);ok{
		w.delObject(obj)
	}
	if emi, ok:=x.(Emitter);ok{
		w.delEmitter(emi)
	}
	if wch, ok:=x.(Watcher);ok{
		w.delWatcher(wch)
	}
}

func (w *World) addObject(obj Object){
	w.objects = append(w.objects, obj)
}

func (w *World) delObject(obj Object){
	l:=len(w.objects)
	for i:=range w.objects{
		if w.objects[i]==obj{
			w.objects[i] = w.objects[l-1]
			w.objects = w.objects[:l-1]
		}
	}
}

func (w*World) addEmitter(emi Emitter){
	w.emitters = append(w.emitters, emi)
	w.vistas = append(w.vistas, New())
}

func (w *World) delEmitter(emi Emitter){
	l:=len(w.emitters)
	for i:=range w.emitters{
		if w.emitters[i]==emi{
			w.emitters[i] = w.emitters[l-1]
			w.emitters = w.emitters[:l-1]
		}
	}
	w.vistas[len(w.vistas)-1].ClearTempSlices(true)
	w.vistas[len(w.vistas)-1].IgnoreSelf = nil
	w.vistas = w.vistas[:l-1]
}

func (w*World) addWatcher(wch Watcher){
	w.watchers = append(w.watchers, wch)
	w.vistas = append(w.vistas, New())
}

func (w *World) delWatcher(wch Watcher){
	l:=len(w.watchers)
	for i:=range w.watchers{
		if w.watchers[i]==wch{
			w.watchers[i] = w.watchers[l-1]
			w.watchers = w.watchers[:l-1]
		}
	}
	w.vistas[len(w.vistas)-1].ClearTempSlices(true)
	w.vistas[len(w.vistas)-1].IgnoreSelf = nil
	w.vistas = w.vistas[:l-1]
}

func (w *World) isObjEmitter(obj Object) bool{
	_,ok:=obj.(Emitter)
	return ok
}

//todo: better vistasCount and clear management