package hero

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
)

type EquipSlot uint16

const (
	SlotHead EquipSlot = iota
	SlotShirt
	SlotPants
	SlotFeet

	equipSlotCount
)

type Equip struct {
	world.VisibleObject

	slot EquipSlot
	kind uint64

	customColors []string

	wearer *Hero // not saved
}

func init() {
	world.Register("equip", world.Visible((*Equip)(nil)))

	world.RegisterSpawnFunc(func(s string) world.Visible {
		for t := range equippables {
			for i, e := range equippables[t] {
				if e.name == s {
					return world.InitObject(&Equip{
						slot: EquipSlot(t),
						kind: uint64(i),
					}).(world.Visible)
				}
			}
		}
		return nil
	})
}

func (e *Equip) Save() (uint, interface{}, []world.ObjectLike) {
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

func (e *Equip) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		dataMap := data.(map[string]interface{})

		e.slot = EquipSlot(dataMap["s"].(uint16))
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

func (e *Equip) Name() string {
	return equippables[e.slot][e.kind].name
}

func (e *Equip) Sprite() string {
	return equippables[e.slot][e.kind].sprite
}

func (e *Equip) Colors() []string {
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

func (e *Equip) AnimationType() string {
	if e.wearer == nil {
		return ""
	}
	return e.wearer.AnimationType()
}

func (e *Equip) SpritePos() (uint, uint) {
	if e.wearer == nil {
		return 12, 0
	}
	return e.wearer.SpritePos()
}
