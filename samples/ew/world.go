package ew

type World struct {
	emitters  []Emitter
	diffusers []Diffuser
	shadowers []Shadower
}

func NewWorld() *World {
	return &World{
		emitters:  make([]Emitter, 0),
		diffusers: make([]Diffuser, 0),
		shadowers: make([]Shadower, 0),
	}
}

func (w *World) AddObject(object interface{}) {
	if obj, ok := object.(Emitter); ok {
		found := false
		for _, v := range w.emitters {
			if v == obj {
				found = true
				break
			}
		}
		if !found {
			w.emitters = append(w.emitters, obj)
		}
	}
	if obj, ok := object.(Diffuser); ok {
		found := false
		for _, v := range w.diffusers {
			if v == obj {
				found = true
				break
			}
		}
		if !found {
			w.diffusers = append(w.diffusers, obj)
		}
	}
	if obj, ok := object.(Shadower); ok {
		found := false
		for _, v := range w.shadowers {
			if v == obj {
				found = true
				break
			}
		}
		if !found {
			w.shadowers = append(w.shadowers, obj)
		}
	}
}

func (w *World) DelObject(object interface{}) {
	if obj, ok := object.(Emitter); ok {
		found := -1
		for i, v := range w.emitters {
			if v == obj {
				found = i
				break
			}
		}
		if found >= 0 {
			w.emitters = append(w.emitters[:found], w.emitters[found+1:]...)
		}
	}
	if obj, ok := object.(Diffuser); ok {
		found := -1
		for i, v := range w.diffusers {
			if v == obj {
				found = i
				break
			}
		}
		if found >= 0 {
			w.diffusers = append(w.diffusers[:found], w.diffusers[found+1:]...)
		}
	}
	if obj, ok := object.(Shadower); ok {
		found := -1
		for i, v := range w.shadowers {
			if v == obj {
				found = i
				break
			}
		}
		if found >= 0 {
			w.shadowers = append(w.shadowers[:found], w.shadowers[found+1:]...)
		}
	}
}
