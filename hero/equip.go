package hero

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/material"
	"github.com/BenLubar/Rnoadm/world"
	"math/big"
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

var _ [equipSlotCount]struct{} = [world.TuningEquipSlotCount]struct{}{}

var equipFacing = [4][]EquipSlot{
	{ // stage front
		SlotShirt,
		SlotPants,
		SlotChest,
		SlotArms,
		SlotShoulders,
		SlotLegs,
		SlotWaist,
		SlotHands,
		SlotFeet,
		SlotHead,
		SlotMainHand,
		SlotOffHand,
		SlotPickaxe,
		SlotHatchet,
	},
	{ // stage back
		SlotShirt,
		SlotPants,
		SlotChest,
		SlotArms,
		SlotShoulders,
		SlotLegs,
		SlotWaist,
		SlotHands,
		SlotFeet,
		SlotHead,
		SlotMainHand,
		SlotOffHand,
		SlotPickaxe,
		SlotHatchet,
	},
	{ // stage left
		SlotShirt,
		SlotPants,
		SlotChest,
		SlotArms,
		SlotShoulders,
		SlotLegs,
		SlotWaist,
		SlotHands,
		SlotFeet,
		SlotHead,
		SlotMainHand,
		SlotOffHand,
		SlotPickaxe,
		SlotHatchet,
	},
	{ // stage right
		SlotShirt,
		SlotPants,
		SlotChest,
		SlotArms,
		SlotShoulders,
		SlotLegs,
		SlotWaist,
		SlotHands,
		SlotFeet,
		SlotHead,
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

	materials []*material.Material

	wearer HeroLike // not saved
}

var _ world.Item = (*Equip)(nil)

func init() {
	world.Register("equip", world.Visible((*Equip)(nil)))

	world.RegisterSpawnFunc(material.WrapSpawnFunc(func(mat *material.Material, s string) world.Visible {
		wood, stone, metal := mat.Get()
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
					if !haveWood && len(wood) != 0 {
						continue
					}
					if !haveStone && len(stone) != 0 {
						continue
					}
					if !haveMetal && len(metal) != 0 {
						continue
					}
					return &Equip{
						slot:      EquipSlot(t),
						kind:      uint64(i),
						materials: []*material.Material{mat},
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
	attached := make([]world.ObjectLike, 0, len(e.materials))
	for _, m := range e.materials {
		attached = append(attached, m)
	}
	return 2, map[string]interface{}{
		"s": uint16(e.slot),
		"k": e.kind,
		"c": colors,
		"m": uint(len(e.materials)),
	}, attached
}

func (e *Equip) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		material := &material.Material{}
		world.InitObject(material)
		material.Load(0, data.(map[string]interface{})["m"], nil)

		attached = []world.ObjectLike{material}
		fallthrough
	case 1:
		data.(map[string]interface{})["m"] = uint(1)
		fallthrough
	case 2:
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

		e.materials = make([]*material.Material, dataMap["m"].(uint))
		for i := range e.materials {
			e.materials[i] = attached[i].(*material.Material)
		}
		attached = attached[len(e.materials):]
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (e *Equip) Name() string {
	var name string
	for _, m := range e.materials {
		name += m.Name()
	}
	return name + equippables[e.slot][e.kind].name
}

func (e *Equip) Examine() (string, [][][2]string) {
	_, info := e.VisibleObject.Examine()

	// TODO: material info

	return equippables[e.slot][e.kind].examine, info
}

func (e *Equip) Sprite() string {
	return equippables[e.slot][e.kind].sprite
}

func (e *Equip) Colors() []string {
	defaultColors := equippables[e.slot][e.kind].colors
	colors := make([]string, len(defaultColors))
	copy(colors, defaultColors)
	for i, m := range e.materials {
		for j, c := range colors {
			switch c {
			case woodColors[i]:
				colors[j] = m.WoodColor()
			case stoneColors[i]:
				colors[j] = m.StoneColor()
			case metalColors[i]:
				colors[j] = m.MetalColor()
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

func (e *Equip) Scale() uint {
	if e.wearer == nil {
		return 1
	}
	return e.wearer.Scale()
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
	var volume uint64
	for _, m := range e.materials {
		volume += m.Volume()
	}
	return volume
}

func (e *Equip) Weight() uint64 {
	var weight uint64
	for _, m := range e.materials {
		weight += m.Weight()
	}
	return weight
}

func (e *Equip) AdminOnly() bool {
	return equippables[e.slot][e.kind].adminOnly
}

func (e *Equip) Actions(player world.PlayerLike) []string {
	actions := e.VisibleObject.Actions(player)
	if e.Position() == nil {
		if e.wearer == nil {
			actions = append(actions, "equip")
		} else {
			actions = append(actions, "unequip")
		}
	}
	return actions
}

func (e *Equip) Interact(player world.PlayerLike, action string) {
	p := player.(*Player)
	switch action {
	case "equip":
		if e.Position() != nil || e.wearer != nil {
			return
		}
		p.Equip(e)
	case "unequip":
		if e.Position() != nil || e.wearer == nil {
			return
		}
		p.Unequip(e.slot)
	default:
		e.VisibleObject.Interact(player, action)
	}
}

func (e *Equip) Stat(stat world.Stat) *big.Int {
	s := &big.Int{}
	for _, m := range e.materials {
		s.Add(s, m.Stat(stat))
	}
	return s
}

func (e *Equip) Quality() *big.Int {
	max := big.NewInt(1)
	for _, m := range e.materials {
		q := m.Quality()
		if max.Cmp(q) < 0 {
			max = q
		}
	}
	return max
}
