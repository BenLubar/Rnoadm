package hero

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/material"
	"github.com/BenLubar/Rnoadm/world"
)

type EquipSlot uint16

const (
	SlotHead EquipSlot = iota
	SlotShirt
	SlotPants
	SlotFeet
	SlotShoulders
	SlotChest
	SlotArms
	SlotHands
	SlotWaist
	SlotLegs
	SlotMainHand
	SlotOffHand
	SlotPickaxe
	SlotHatchet

	equipSlotCount
)

var equipFacing = [4][]EquipSlot{
	{ // stage front
		SlotHead,
		SlotShirt,
		SlotPants,
		SlotFeet,
		SlotShoulders,
		SlotChest,
		SlotArms,
		SlotHands,
		SlotWaist,
		SlotLegs,
		SlotMainHand,
		SlotOffHand,
		SlotPickaxe,
		SlotHatchet,
	},
	{ // stage back
		SlotHead,
		SlotShirt,
		SlotPants,
		SlotFeet,
		SlotShoulders,
		SlotChest,
		SlotArms,
		SlotHands,
		SlotWaist,
		SlotLegs,
		SlotMainHand,
		SlotOffHand,
		SlotPickaxe,
		SlotHatchet,
	},
	{ // stage left
		SlotHead,
		SlotShirt,
		SlotPants,
		SlotFeet,
		SlotShoulders,
		SlotChest,
		SlotArms,
		SlotHands,
		SlotWaist,
		SlotLegs,
		SlotMainHand,
		SlotOffHand,
		SlotPickaxe,
		SlotHatchet,
	},
	{ // stage right
		SlotHead,
		SlotShirt,
		SlotPants,
		SlotFeet,
		SlotShoulders,
		SlotChest,
		SlotArms,
		SlotHands,
		SlotWaist,
		SlotLegs,
		SlotMainHand,
		SlotOffHand,
		SlotPickaxe,
		SlotHatchet,
	},
}

type Equip struct {
	world.VisibleObject

	slot EquipSlot
	kind uint64

	customColors []string

	wood  *material.WoodType
	stone *material.StoneType
	metal *material.MetalType

	wearer *Hero // not saved
}

var _ world.Item = (*Equip)(nil)

func init() {
	world.Register("equip", world.Visible((*Equip)(nil)))

	world.RegisterSpawnFunc(material.WrapSpawnFunc(func(wood *material.WoodType, stone *material.StoneType, metal *material.MetalType, s string) world.Visible {
		for t := range equippables {
			for i, e := range equippables[t] {
				if e.name == s {
					haveWood := false
					haveStone := false
					haveMetal := false
					for _, c := range e.colors {
						if c == woodColor {
							haveWood = true
						}
						if c == stoneColor {
							haveStone = true
						}
						if c == metalColor {
							haveMetal = true
						}
					}
					if !haveWood && wood != nil {
						continue
					}
					if !haveStone && stone != nil {
						continue
					}
					if !haveMetal && metal != nil {
						continue
					}
					return &Equip{
						slot:  EquipSlot(t),
						kind:  uint64(i),
						wood:  wood,
						stone: stone,
						metal: metal,
					}
				}
			}
		}
		return nil
	}))
}

func (e *Equip) Save() (uint, interface{}, []world.ObjectLike) {
	colors := make([]interface{}, len(e.customColors))
	for i, c := range e.customColors {
		colors[i] = c
	}
	materials := make(map[string]interface{})
	if e.wood != nil {
		materials["w"] = uint64(*e.wood)
	}
	if e.stone != nil {
		materials["s"] = uint64(*e.stone)
	}
	if e.metal != nil {
		materials["m"] = uint64(*e.metal)
	}
	return 0, map[string]interface{}{
		"s": uint16(e.slot),
		"k": e.kind,
		"c": colors,
		"m": materials,
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
		materials, _ := dataMap["m"].(map[string]interface{})
		if wood, ok := materials["w"].(uint64); ok {
			e.wood = (*material.WoodType)(&wood)
		}
		if stone, ok := materials["s"].(uint64); ok {
			e.stone = (*material.StoneType)(&stone)
		}
		if metal, ok := materials["m"].(uint64); ok {
			e.metal = (*material.MetalType)(&metal)
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
	for i, c := range colors {
		switch c {
		case woodColor:
			if e.wood == nil {
				colors[i] = ""
			} else {
				colors[i] = e.wood.BarkColor()
			}
		case stoneColor:
			if e.stone == nil {
				colors[i] = ""
			} else {
				colors[i] = e.stone.Color()
			}
		case metalColor:
			if e.metal == nil {
				colors[i] = ""
			} else {
				colors[i] = e.metal.Color()
			}
		}
	}
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
	var animation string
	if e.wearer != nil {
		animation = e.wearer.AnimationType()
	}
	if a, ok := equippables[e.slot][e.kind].animationOverrides[animation]; ok {
		animation = a
	}
	return animation
}

func (e *Equip) SpritePos() (uint, uint) {
	if e.wearer == nil {
		return 12, 0
	}
	return e.wearer.SpritePos()
}

func (e *Equip) SpriteSize() (uint, uint) {
	return equippables[e.slot][e.kind].width, equippables[e.slot][e.kind].height
}

func (e *Equip) Volume() uint64 {
	return 1 // TODO
}

func (e *Equip) Weight() uint64 {
	return 0 // TODO
}

func (e *Equip) AdminOnly() bool {
	return equippables[e.slot][e.kind].adminOnly
}

func (e *Equip) Actions() []string {
	actions := e.VisibleObject.Actions()
	if e.Position() == nil {
		if e.wearer == nil {
			actions = append(actions, "equip")
		} else {
			actions = append(actions, "unequip")
		}
	}
	return actions
}
