package nigiri

type Transforms []RTransformer

func (t Transforms) TransformRect(rect Rect) Rect {
	for _, v := range t {
		if v == nil {
			continue
		}
		rect = v.TransformRect(rect)
		if rect.Empty() {
			return ZR
		}
	}
	return rect
}

func (t Transforms) lastClipper() Clipper {
	for i := len(t) - 1; i >= 0; i-- {
		if clipper, ok := t[i].(Clipper); ok {
			return clipper
		}
		if trans, ok := t[i].(Transforms); ok {
			clipper := trans.lastClipper()
			if clipper != nil {
				return clipper
			}
		}
	}
	return nil
}
