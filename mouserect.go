package nigiri

import "image"

type MouseRect struct {
	MouseActive bool
	CatchRect   image.Rectangle
	childs      []OnMouser
	selfClick   OnMouserF
}

func NewClickRect(selfF OnMouserF) MouseRect {
	return MouseRect{
		MouseActive: true,
		childs:      make([]OnMouser, 0),
		selfClick:   selfF,
	}
}

func (c *MouseRect) AddChild(child OnMouser) {
	for _, v := range c.childs {
		if v == child {
			return
		}
	}
	c.childs = append(c.childs, child)
}

func (c *MouseRect) DeleteChild(child OnMouser) {
	for i, v := range c.childs {
		if v == child {
			c.childs = append(c.childs[:i], c.childs[i+1:]...)
			return
		}
	}
}

func (c MouseRect) OnMouse(x, y int) bool {
	if !c.MouseActive {
		return false
	}
	if !image.Pt(x, y).In(c.CatchRect) {
		return false
	}
	for _, child := range c.childs {
		if child != nil && child.OnMouse(x, y) {
			return true
		}
	}
	if c.selfClick != nil {
		return c.selfClick(x, y)
	}
	return false
}
