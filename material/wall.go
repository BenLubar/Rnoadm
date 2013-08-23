package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
)

type Wall struct {
	world.VisibleObject

	wood *WoodType
	stone *StoneType
	metal *MetalType
}

func init() {
	world.Register("wall", world.Visible((*Wall)(nil)))
	world.RegisterSpawnFunc(WrapSpawnFunc(func(wood *WoodType, stone *StoneType, metal *MetalType, s string) world.Visible {
		if s == "wall" {
			return &Wall{
				wood:wood,
				stone:stone,
				metal:metal,
			}
		}
		return nil
	}))
}

func (w *Wall) Name() string {
	return "wall"
}

func (w *Wall) Sprite() string {
	return "wall"
}

func (w *Wall) Colors() []string {
	if w.metal == nil {
		if w.stone == nil {
			if w.wood == nil {
				return []string{"#888"}
			} else {
				return []string{"", w.wood.BarkColor()}
			}
		} else {
			if w.wood == nil {
				return []string{"", "", w.stone.Color()}
			} else {
				return []string{"", "", w.stone.Color(), "", w.wood.BarkColor()}
			}
		}
	} else {
		if w.stone == nil {
			if w.wood == nil {
				return []string{"", "", "", w.metal.Color()}
			} else {
				return []string{"", "", "", w.metal.Color(), w.wood.BarkColor()}
			}
		} else {
			if w.wood == nil {
				return []string{"", "", w.stone.Color(), "", "", w.metal.Color()}
			} else {
				return []string{"", "", w.stone.Color(), "", w.wood.BarkColor(), w.metal.Color()}
			}
		}
	}
}

func containsWall(objects []world.ObjectLike) bool {
	for _, o := range objects {
		if _, ok := o.(IsWall); ok {
			return true
		}
	}
	return false
}

type IsWall interface {
	isWall() IsWall
}

func (w *Wall) isWall() IsWall {
	return w
}

func (w *Wall) SpritePos() (uint, uint) {
	var sides uint
	if pos := w.Position(); pos != nil {
		x, y := pos.Position()
		z := pos.Zone()

		if x > 0 && containsWall(z.Tile(x-1, y).Objects()) {
			sides |= 1
		}
		if x < 255 && containsWall(z.Tile(x+1, y).Objects()) {
			sides |= 2
		}
		if y > 0 && containsWall(z.Tile(x, y-1).Objects()) {
			sides |= 4
		}
		if y < 255 && containsWall(z.Tile(x, y+1).Objects()) {
			sides |= 8
		}
	}
	return sides, 0
}

func (w *Wall) SpriteSize() (uint, uint) {
	return 32, 64
}

func (w *Wall) Blocking() bool {
	return true
}

func (w *Wall) Save() (uint, interface{}, []world.ObjectLike) {
	materials := make(map[string]interface{})
	if w.wood != nil {
		materials["w"] = uint64(*w.wood)
	}
	if w.stone != nil {
		materials["s"] = uint64(*w.stone)
	}
	if w.metal != nil {
		materials["m"] = uint64(*w.metal)
	}
	return 1, materials, nil
}

func (w *Wall) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		// no fields in version 0
	case 1:
		materials := data.(map[string]interface{})
		if wood, ok := materials["w"].(uint64); ok {
			w.wood = (*WoodType)(&wood)
		}
		if stone, ok := materials["s"].(uint64); ok {
			w.stone = (*StoneType)(&stone)
		}
		if metal, ok := materials["m"].(uint64); ok {
			w.metal = (*MetalType)(&metal)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}
