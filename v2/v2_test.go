package v2

import (
	"github.com/Shnifer/magellan/v2"
	"math"
	"testing"
)

func Test_DirInDir(t *testing.T) {
	for i := 0; i < 360; i++ {
		ang := float64(i)
		v := v2.InDir(ang)
		if math.Abs(v.Dir()-ang) > 0.0001 {
			t.Errorf("for ang %v generated vector %v but it's Dir is %v", ang, v, v.Dir())
		}
	}
}
