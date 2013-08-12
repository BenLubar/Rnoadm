package world

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	var z1 Zone
	z1.X = 1
	z1.Y = -1
	z1.Tile(42, 1).Add(&Object{})
	z1.Tile(42, 1).Add(&VisibleObject{})
	z1.Tile(100, 200).Add(&Object{})

	var buf bytes.Buffer
	writeZone(&z1, &buf)
	z2 := readZone(&buf)

	if z1.X != z2.X {
		t.Errorf("X(%d) != X(%d)", z1.X, z2.X)
	}
	if z1.Y != z2.Y {
		t.Errorf("X(%d) != X(%d)", z1.Y, z2.Y)
	}
	for x := 0; x < 256; x++ {
		x8 := uint8(x)
		for y := 0; y < 256; y++ {
			y8 := uint8(y)
			t1 := z1.Tile(x8, y8)
			t2 := z2.Tile(x8, y8)
			o1 := t1.Objects()
			o2 := t2.Objects()
			if len(o1) != len(o2) {
				t.Errorf("(%d, %d) len(%d) != len(%d)", x, y, len(o1), len(o2))
			} else {
				for i := range o1 {
					if reflect.TypeOf(o1[i]) != reflect.TypeOf(o2[i]) {
						t.Errorf("(%d, %d)[%d] %T != %T", x, y, i, o1[i], o2[i])
					}
				}
			}
		}
	}
}
