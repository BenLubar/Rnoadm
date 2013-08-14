package hero

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
)

type EquippableSlot uint16

const (
	SlotHead EquippableSlot = iota
	SlotShirt
	SlotPants
	SlotFeet

	equippableSlotCount
)

type Equippable struct {
	world.VisibleObject

	slot EquippableSlot
	kind uint64

	customColors []string

	wearer *Hero // not saved
}

func init() {
	world.Register("equip", world.Visible((*Equippable)(nil)))
}

func (e *Equippable) Save() (uint, interface{}, []world.ObjectLike) {
	colors := make([]interface{}, len(e.customColors))
	for i, c := range e.customColors {
		colors[i] = c
	}
	return 0, map[string]interface{}{
		"s": uint16(e.slot),
		"k": e.kind,
		"c": colors,
	}, nil
}

func (e *Equippable) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		dataMap := data.(map[string]interface{})

		e.slot = EquippableSlot(dataMap["s"].(uint16))
		e.kind = dataMap["k"].(uint64)
		colors := dataMap["c"].([]interface{})
		if len(colors) > 0 {
			e.customColors = make([]string, len(colors))
			for i, c := range colors {
				e.customColors[i] = c.(string)
			}
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (e *Equippable) Sprite() string {
	return equippables[e.slot][e.kind].sprite
}

func (e *Equippable) Colors() []string {
	defaultColors := equippables[e.slot][e.kind].colors
	colors := make([]string, len(defaultColors))
	copy(colors, defaultColors)
	for i, c := range e.customColors {
		if i >= len(colors) {
			break
		}
		if c != "" {
			colors[i] = c
		}
	}
	return colors
}

func (e *Equippable) AnimationType() string {
	if e.wearer == nil {
		return ""
	}
	return e.wearer.AnimationType()
}

func (e *Equippable) SpritePos() (uint, uint) {
	if e.wearer == nil {
		return 16, 0
	}
	return e.wearer.SpritePos()
}
