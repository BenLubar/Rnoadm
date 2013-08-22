package material

import (
	"github.com/BenLubar/Rnoadm/world"
)

func WrapSpawnFunc(f func(*WoodType, *StoneType, *MetalType, string) world.Visible) func(string) world.Visible {
	return func(s string) world.Visible {
		prefix := func(p string) bool {
			return len(s) > len(p) && s[:len(p)] == p && s[len(p)] == ' '
		}

		var wood *WoodType
		var stone *StoneType
		var metal *MetalType
	find:
		for {
			if v := f(wood, stone, metal, s); v != nil {
				return v
			}
			if wood == nil {
				for i := range woodTypes {
					t := WoodType(i)
					if prefix(t.Name()) {
						wood = &t
						s = s[len(t.Name())+1:]
						continue find
					}
				}
			}
			if stone == nil {
				for i := range stoneTypes {
					t := StoneType(i)
					if prefix(t.Name()) {
						stone = &t
						s = s[len(t.Name())+1:]
						continue find
					}
				}
			}
			if metal == nil {
				for i := range metalTypes {
					t := MetalType(i)
					if prefix(t.Name()) {
						metal = &t
						s = s[len(t.Name())+1:]
						continue find
					}
				}
			}
			return nil
		}
	}
}
