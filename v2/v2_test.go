package v2

import (
	"math"
	"testing"
)

func Test_DirInDir(t *testing.T) {
	for i := 0; i < 360; i++ {
		ang := float64(i)
		v := InDir(ang)
		if math.Abs(v.Dir()-ang) > 0.0001 {
			t.Errorf("for ang %v generated vector %v but it's Dir is %v", ang, v, v.Dir())
		}
	}
}

func Test_RotateInDir(t *testing.T) {
	rangs := [...]float64{0, 0.01, 1, 45, 90, 180, 360, 721, -0.01, -1, -45, -90, -180, -360}
	for _, rang := range rangs {
		for i := 0; i < 360; i++ {
			ang := float64(i)
			v := InDir(ang)
			v = v.Rotate(rang)
			d := NormAng(v.Dir() - ang - rang)

			if d > 0.0001 && d < 359.999 {
				t.Errorf("for ang %v rotated by %v vector %v but it's Dir is %v", ang, rang, v, v.Dir())
			}
		}
	}
}

func TestV2_Rotate90(t *testing.T) {

	for i := 0; i < 360; i++ {
		sAng := float64(i)
		v := InDir(sAng)
		v = v.Rotate90()
		wait := NormAng(sAng + 90)
		d := wait - v.Dir()
		if d > 0.0001 && d < 359.99 {
			t.Errorf("for start vector with ang %v result of rotate90 is %v needed %v\n", sAng, v.Dir(), wait)
		}
	}
}

func TestDir(t *testing.T) {
	v := V(100, 0)
	if v.Dir() != 270 {
		t.Errorf("{100,0}.InDir waited 270 got %v", v.Dir())
	}
	v = V(-100, 0)
	if v.Dir() != 90 {
		t.Errorf("{-100,0}.InDir waited 90 got %v", v.Dir())
	}
}
