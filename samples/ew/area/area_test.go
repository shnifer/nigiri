package area

import (
	"testing"
	"github.com/shnifer/nigiri/vec2/angle"
)

func (a Area) Contains1 (b Area) bool{
	return a.Height>=b.Height && a.Period.Contains(b.Period)
}

func (a Area) Contains2 (b Area) bool{
	return a.Height>=b.Height &&(a.IsFull() || !b.IsFull())&& a.Has(b.Start()) && (a.Has(b.End()) || a.End()==b.End())
}

func (a Area) Contains3 (b Area) bool{
	if b.IsFull(){
		return a.IsFull() && a.Height>=b.Height
	}
	return a.Height>=b.Height && a.Has(b.Start()) && (a.Has(b.End()) || a.End()==b.End())
}

func (a Area) Contains4 (b Area) bool{
	if a.Height<b.Height{
		return false
	}
	if b.IsFull(){
		return a.IsFull()
	}
	return a.Has(b.Start()) && (a.Has(b.End()) || a.End()==b.End())
}

func vsSetup() []Area{
	vs:=make([]Area,0)
	for i:=0; i<10; i++{
		for j :=0; j <10; j++{
			per:=angle.NewPeriod(float64(i)*30, float64(j)*30)
			vs = append(vs, New(per, per.Wide()))
		}
	}
	vs = append(vs, New(angle.FullPeriod, 180))
	return vs
}

func BenchmarkArea_Contains1(b *testing.B) {
	vs := vsSetup()
	b.ResetTimer()
	var res bool
	_=res
	for m :=0; m <b.N; m++{
		for i:=range vs{
			for j:=range vs{
				res = vs[i].Contains1(vs[j])
			}
		}
	}
}

func BenchmarkArea_Contains2(b *testing.B) {
	vs := vsSetup()
	b.ResetTimer()
	var res bool
	_=res
	for m :=0; m <b.N; m++{
		for i:=range vs{
			for j:=range vs{
				res = vs[i].Contains2(vs[j])
			}
		}
	}
}

func BenchmarkArea_Contains3(b *testing.B) {
	vs := vsSetup()
	b.ResetTimer()
	var res bool
	_=res
	for m :=0; m <b.N; m++{
		for i:=range vs{
			for j:=range vs{
				res = vs[i].Contains3(vs[j])
			}
		}
	}
}

func BenchmarkArea_Contains4(b *testing.B) {
	vs := vsSetup()
	b.ResetTimer()
	var res bool
	_=res
	for m :=0; m <b.N; m++{
		for i:=range vs{
			for j:=range vs{
				res = vs[i].Contains4(vs[j])
			}
		}
	}
}