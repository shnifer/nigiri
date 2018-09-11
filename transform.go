package nigiri

type Transforms []Transformer

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
