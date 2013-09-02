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

	material *material.Material

	wearer HeroLike // not saved
}

var _ world.Item = (*Equip)(nil)

func init() {
	world.Register("equip", world.Visible((*Equip)(nil)))

	world.RegisterSpawnFunc(material.WrapSpawnFunc(func(material *material.Material, s string) world.Visible {
		wood, stone, metal := material.Get()
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
						slot:     EquipSlot(t),
						kind:     uint64(i),
						material: material,
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
	if e.material == nil {
		e.material = &material.Material{}
	}
	return 1, map[string]interface{}{
		"s": uint16(e.slot),
		"k": e.kind,
		"c": colors,
	}, []world.ObjectLike{e.material}
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

		e.material = attached[0].(*material.Material)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (e *Equip) Name() string {
	return e.material.Name() + equippables[e.slot][e.kind].name
}

func (e *Equip) Examine() (string, [][][2]string) {
	_, info := e.VisibleObject.Examine()

	info = append(info, e.material.Info()...)

	return equippables[e.slot][e.kind].examine, info
}

func (e *Equip) Sprite() string {
	return equippables[e.slot][e.kind].sprite
}

func (e *Equip) Colors() []string {
	defaultColors := equippables[e.slot][e.kind].colors
	colors := make([]string, len(defaultColors))
	copy(colors, defaultColors)
	wood, stone, metal := e.material.Get()
	for i, c := range colors {
		switch c {
		case woodColor:
			if wood == nil {
				colors[i] = ""
			} else {
				colors[i] = wood.BarkColor()
			}
		case stoneColor:
			if stone == nil {
				colors[i] = ""
			} else {
				colors[i] = stone.Color()
			}
		case metalColor:
			if metal == nil {
				colors[i] = ""
			} else {
				colors[i] = metal.Color()
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
	return 1 // TODO
}

func (e *Equip) Weight() uint64 {
	return 0 // TODO
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
		p.mtx.Lock()
		if old, ok := p.equipped[e.slot]; ok {
			p.mtx.Unlock()
			old.Interact(player, "unequip")
			p.mtx.Lock()
		}
		if _, ok := p.equipped[e.slot]; ok {
			p.mtx.Unlock()
			return
		}
		for i, item := range p.items {
			if e == item {
				p.items = append(p.items[:i], p.items[i+1:]...)
				p.equipped[e.slot] = e
				e.wearer = &p.Hero
				p.notifyInventoryChanged()
				p.mtx.Unlock()
				p.Position().Zone().Update(p.Position(), p)
				return
			}
		}
		p.mtx.Unlock()
	case "unequip":
		if e.Position() != nil || e.wearer == nil {
			return
		}
		p.mtx.Lock()
		if p.equipped[e.slot] != e {
			p.mtx.Unlock()
			return
		}
		if !p.canHoldItem(e) {
			p.mtx.Unlock()
			return
		}
		e.wearer = nil
		p.giveItem(e)
		delete(p.equipped, e.slot)
		p.notifyInventoryChanged()
		p.mtx.Unlock()
		p.Position().Zone().Update(p.Position(), p)
	default:
		e.VisibleObject.Interact(player, action)
	}
}
