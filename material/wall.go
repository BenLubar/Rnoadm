package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
)

type Wall struct {
	world.VisibleObject

	material *Material
}

func init() {
	world.Register("wall", world.Visible((*Wall)(nil)))
	world.RegisterSpawnFunc(WrapSpawnFunc(func(material *Material, s string) world.Visible {
		if s == "wall" {
			return &Wall{
				material: material,
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
	wood, stone, metal := w.material.Get()
	if metal == nil {
		if stone == nil {
			if wood == nil {
				return []string{"#888"}
			} else {
				return []string{"", wood.BarkColor()}
			}
		} else {
			if wood == nil {
				return []string{"", "", stone.Color()}
			} else {
				return []string{"", "", stone.Color(), "", wood.BarkColor()}
			}
		}
	} else {
		if stone == nil {
			if wood == nil {
				return []string{"", "", "", metal.Color()}
			} else {
				return []string{"", "", "", metal.Color(), wood.BarkColor()}
			}
		} else {
			if wood == nil {
				return []string{"", "", stone.Color(), "", "", metal.Color()}
			} else {
				return []string{"", "", stone.Color(), "", wood.BarkColor(), metal.Color()}
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
	return 2, uint(0), []world.ObjectLike{w.material}
}

func (w *Wall) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		// no fields in version 0
	case 1:
		w.material = &Material{}
		world.InitObject(w.material)
		w.material.Load(0, data.(map[string]interface{}), nil)
	case 2:
		w.material = attached[0].(*Material)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}
