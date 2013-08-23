package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
)

type Floor struct {
	world.VisibleObject

	material *Material
}

func init() {
	world.Register("floor", world.Visible((*Floor)(nil)))
	world.RegisterSpawnFunc(WrapSpawnFunc(func(material *Material, s string) world.Visible {
		if s == "floor" {
			return &Floor{
				material: material,
			}
		}
		return nil
	}))
}

func (f *Floor) Name() string {
	return "floor"
}

func (f *Floor) Sprite() string {
	return "floor"
}

func (f *Floor) Colors() []string {
	return []string{"#888"}
}

func (f *Floor) SpritePos() (uint, uint) {
	return 0, 0
}

func (f *Floor) SpriteSize() (uint, uint) {
	return 64, 64
}

func (f *Floor) AnimationType() string {
	if f.Position() == nil {
		return ""
	}
	return "_fl"
}

func (f *Floor) Blocking() bool {
	return false
}

func (f *Floor) Save() (uint, interface{}, []world.ObjectLike) {
	return 2, uint(0), []world.ObjectLike{f.material}
}

func (f *Floor) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		// no fields in version 0
	case 1:
		f.material = &Material{}
		world.InitObject(f.material)
		f.material.Load(0, data.(map[string]interface{}), nil)
	case 2:
		f.material = attached[0].(*Material)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}
