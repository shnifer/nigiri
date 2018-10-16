package vista

import "github.com/shnifer/nigiri/vec2/angle"

type Area struct {
	angle.Period
	//total height, from 0-180
	//Full angle is supposed to have 180
	//and Ray to have 0
	Height float64
}

func NewArea(period angle.Period, height float64) Area{
	if height>180{
		height = 180
	} else if height<0{
		height = 0
	}
	return Area{
		Period: period,
		Height: height,
	}
}

func (a Area) Contains (b Area) bool{
	return a.Height>=b.Height && a.Period.Contains(b.Period)
}