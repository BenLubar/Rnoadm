package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
)

type Floor struct {
	world.VisibleObject

	wood *WoodType
	stone *StoneType
	metal *MetalType
}

func init() {
	world.Register("floor", world.Visible((*Floor)(nil)))
	world.RegisterSpawnFunc(WrapSpawnFunc(func(wood *WoodType, stone *StoneType, metal *MetalType, s string) world.Visible {
		if s == "floor" {
			return &Floor{
				wood:wood,
				stone:stone,
				metal:metal,
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
	materials := make(map[string]interface{})
	if f.wood != nil {
		materials["w"] = uint64(*f.wood)
	}
	if f.stone != nil {
		materials["s"] = uint64(*f.stone)
	}
	if f.metal != nil {
		materials["m"] = uint64(*f.metal)
	}
	return 1, materials, nil
}

func (f *Floor) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		// no fields in version 0
	case 1:
		materials := data.(map[string]interface{})
		if wood, ok := materials["w"].(uint64); ok {
			f.wood = (*WoodType)(&wood)
		}
		if stone, ok := materials["s"].(uint64); ok {
			f.stone = (*StoneType)(&stone)
		}
		if metal, ok := materials["m"].(uint64); ok {
			f.metal = (*MetalType)(&metal)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}
