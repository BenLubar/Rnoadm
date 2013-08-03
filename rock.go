package main

import (
	"math"
)

type RockType uint16

const (
	Granite RockType = iota
	Adminstone
	Limestone

	rockTypeCount
)

var rockTypeInfo = [rockTypeCount]struct {
	Name     string
	Color    Color
	Strength uint64
	sqrtStr  uint64
}{
	Granite: {
		Name:     "granite",
		Color:    "#948e85",
		Strength: 50,
	},
	Adminstone: {
		Name:     "adminstone",
		Color:    "#1e0036",
		Strength: 1 << 62,
	},
	Limestone: {
		Name:     "limestone",
		Color:    "#bebd8f",
		Strength: 50,
	},
}

func init() {
	for t := range rockTypeInfo {
		if rockTypeInfo[t].Strength >= 1<<60 {
			rockTypeInfo[t].sqrtStr = rockTypeInfo[t].Strength - 1
		} else {
			rockTypeInfo[t].sqrtStr = uint64(math.Sqrt(float64(rockTypeInfo[t].Strength)))
		}
	}
}

type MetalType uint16

const (
	_ MetalType = iota
	Iron
	Unobtainium
	Copper

	metalTypeCount
)

var metalTypeInfo = [metalTypeCount]struct {
	Name     string
	Color    Color
	Strength uint64
	sqrtStr  uint64
}{
	Iron: {
		Name:     "iron",
		Color:    "#79493d",
		Strength: 50,
	},
	Unobtainium: {
		Name:     "unobtainium",
		Color:    "#ff00ff",
		Strength: 1 << 62,
	},
	Copper: {
		Name:     "copper",
		Color:    "#af633e",
		Strength: 50,
	},
}

func init() {
	for t := range metalTypeInfo {
		if metalTypeInfo[t].Strength >= 1<<60 {
			metalTypeInfo[t].sqrtStr = metalTypeInfo[t].Strength - 1
		} else {
			metalTypeInfo[t].sqrtStr = uint64(math.Sqrt(float64(metalTypeInfo[t].Strength)))
		}
	}
}

type Rock struct {
	Type RockType
	Ore  MetalType
	Big  bool
}

func (r *Rock) Name() string {
	return rockTypeInfo[r.Type].Name + " rock"
}

func (r *Rock) Examine() string {
	if r.Ore != 0 {
		if r.Big {
			return "a deposit of " + rockTypeInfo[r.Type].Name + " rock containing rich " + metalTypeInfo[r.Ore].Name + " ore."
		}
		return "a deposit of " + rockTypeInfo[r.Type].Name + " rock containing " + metalTypeInfo[r.Ore].Name + " ore."
	}
	return "a deposit of " + rockTypeInfo[r.Type].Name + " rock."
}

/*func (r *Rock) Paint(x, y int, setcell func(int, int, PaintCell)) {
	setcell(x, y, PaintCell{
		Sprite: "rock_base",
		Color:  rockTypeInfo[r.Type].Color,
		ZIndex: 51,
	})
	if r.Ore != 0 {
		setcell(x, y, PaintCell{
			Sprite: "rock_ore_small",
			Color:  metalTypeInfo[r.Ore].Color,
			ZIndex: 52,
		})
		if r.Big {
			setcell(x, y, PaintCell{
				Sprite: "rock_ore_big",
				Color:  metalTypeInfo[r.Ore].Color,
				ZIndex: 53,
			})
		}
	}
}*/

func (r *Rock) Blocking() bool {
	return true
}

func (r *Rock) InteractOptions() []string {
	return []string{"mine", "quarry"}
}

func (r *Rock) Interact(x uint8, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // mine
		player.Lock()
		var schedule Schedule = &MineQuarrySchedule{X: x, Y: y, R: r, Mine: true}
		if tx, ty := player.TileX, player.TileY; (tx-x)*(tx-x)+(ty-y)*(ty-y) > 1 {
			moveSchedule := MoveSchedule(FindPath(zone, tx, ty, x, y, false))
			schedule = &ScheduleSchedule{&moveSchedule, schedule}
		}
		player.schedule = schedule
		player.Unlock()
	case 1: // quarry
		player.Lock()
		var schedule Schedule = &MineQuarrySchedule{X: x, Y: y, R: r, Mine: false}
		if tx, ty := player.TileX, player.TileY; (tx-x)*(tx-x)+(ty-y)*(ty-y) > 1 {
			moveSchedule := MoveSchedule(FindPath(zone, tx, ty, x, y, false))
			schedule = &ScheduleSchedule{&moveSchedule, schedule}
		}
		player.schedule = schedule
		player.Unlock()
	}
}

func (r *Rock) ZIndex() int {
	return 0
}

type Stone struct {
	Type RockType
	Uninteractable
}

func (s *Stone) Name() string {
	return rockTypeInfo[s.Type].Name + " stone"
}

func (s *Stone) Examine() string {
	return "a " + rockTypeInfo[s.Type].Name + " stone."
}

/*func (s *Stone) Paint(x, y int, setcell func(int, int, PaintCell)) {
	setcell(x, y, PaintCell{
		Sprite: "item_stone",
		Color:  rockTypeInfo[s.Type].Color,
		ZIndex: 75,
	})
}*/

func (s *Stone) Blocking() bool {
	return false
}

func (s *Stone) IsItem() {}

func (s *Stone) AdminOnly() bool {
	return rockTypeInfo[s.Type].Strength >= 1<<60
}

func (s *Stone) ZIndex() int {
	return 25
}

type Ore struct {
	Type MetalType
	Uninteractable
}

func (o *Ore) Name() string {
	return metalTypeInfo[o.Type].Name + " ore"
}

func (o *Ore) Examine() string {
	return "some " + metalTypeInfo[o.Type].Name + " ore."
}

/*func (o *Ore) Paint(x, y int, setcell func(int, int, PaintCell)) {
	setcell(x, y, PaintCell{
		Sprite: "item_ore",
		Color:  metalTypeInfo[o.Type].Color,
		ZIndex: 75,
	})
}*/

func (o *Ore) Blocking() bool {
	return false
}

func (o *Ore) IsItem() {}

func (o *Ore) AdminOnly() bool {
	return metalTypeInfo[o.Type].Strength >= 1<<60
}

func (o *Ore) ZIndex() int {
	return 25
}

type Pickaxe struct {
	Head   MetalType
	Handle WoodType
}

func (p *Pickaxe) Name() string {
	return metalTypeInfo[p.Head].Name + " pickaxe"
}

func (p *Pickaxe) Examine() string {
	return "a pickaxe made from " + metalTypeInfo[p.Head].Name + " and " + woodTypeInfo[p.Handle].Name + "."
}

/*func (p *Pickaxe) Paint(x, y int, setcell func(int, int, PaintCell)) {
	setcell(x, y, PaintCell{
		Sprite: "item_tool_handle",
		Color:  woodTypeInfo[p.Handle].Color,
		ZIndex: 75,
	})
	setcell(x, y, PaintCell{
		Sprite: "item_tool_pickaxe",
		Color:  metalTypeInfo[p.Head].Color,
		ZIndex: 76,
	})
}

func (p *Pickaxe) PaintWorn(x, y int, setcell func(int, int, PaintCell), frame uint8, offsetX, offsetY int8) {
	setcell(x, y, PaintCell{
		Sprite: "tiny_pick_stick",
		Color:  woodTypeInfo[p.Handle].Color,
		SheetX: frame,
		X:      offsetX,
		Y:      offsetY,
		ZIndex: 506,
	})
	setcell(x, y, PaintCell{
		Sprite: "tiny_pick_head",
		Color:  metalTypeInfo[p.Head].Color,
		SheetX: frame,
		X:      offsetX,
		Y:      offsetY,
		ZIndex: 507,
	})
}*/

func (p *Pickaxe) Blocking() bool {
	return false
}

func (p *Pickaxe) InteractOptions() []string {
	return []string{"add to toolbelt"}
}

func (p *Pickaxe) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // add to toolbelt
		player.Lock()
		player.Equip(p, true)
		player.Unlock()
	}
}

func (p *Pickaxe) IsItem() {}

func (p *Pickaxe) AdminOnly() bool {
	return metalTypeInfo[p.Head].Strength >= 1<<60 || woodTypeInfo[p.Handle].Strength >= 1<<60
}

func (p *Pickaxe) ZIndex() int {
	return 25
}

type MineQuarrySchedule struct {
	Delayed bool
	Mine    bool
	X, Y    uint8
	R       *Rock
}

func (s *MineQuarrySchedule) Act(z *Zone, x uint8, y uint8, h *Hero, p *Player) bool {
	if !s.Delayed {
		s.Delayed = true
		h.scheduleDelay = 10
		if p != nil {
			if s.Mine {
				p.SendMessage("you attempt to mine the " + s.R.Name() + ".")
			} else {
				p.SendMessage("you attempt to quarry the " + s.R.Name() + ".")
			}
		}
		return true
	}
	if (s.X-x)*(s.X-x)+(s.Y-y)*(s.Y-y) > 1 {
		if p != nil {
			p.SendMessage("that is too far away!")
		}
		return false
	}

	if s.Mine && s.R.Ore == 0 {
		if p != nil {
			p.SendMessage("there is no ore in this rock.")
		}
		return false
	}

	h.Lock()
	h.Delay++
	pickaxe := h.Toolbelt.Pickaxe
	h.Unlock()
	if pickaxe == nil {
		if p != nil {
			p.SendMessage("you do not have a pickaxe in your toolbelt.")
		}
		return false
	}

	pickaxeMax := metalTypeInfo[pickaxe.Head].Strength + woodTypeInfo[pickaxe.Handle].Strength
	pickaxeMin := metalTypeInfo[pickaxe.Head].sqrtStr + woodTypeInfo[pickaxe.Handle].sqrtStr

	var rockMax, rockMin uint64
	if s.Mine {
		rockMax = metalTypeInfo[s.R.Ore].Strength
		rockMin = metalTypeInfo[s.R.Ore].sqrtStr
	} else {
		rockMax = rockTypeInfo[s.R.Type].Strength
		rockMin = rockTypeInfo[s.R.Type].sqrtStr
	}

	z.Lock()
	r := z.Rand()
	pickaxeScore := uint64(r.Int63n(int64(pickaxeMax-pickaxeMin+1))) + pickaxeMin
	rockScore := uint64(r.Int63n(int64(rockMax-rockMin+1))) + rockMin
	if p != nil {
		switch {
		case pickaxeScore < rockScore/5:
			p.SendMessage("your " + pickaxe.Name() + " doesn't even make a dent in the " + s.R.Name() + ".")
		case pickaxeScore < rockScore*2/3:
			p.SendMessage("your " + pickaxe.Name() + " slightly dents the " + s.R.Name() + ", but nothing interesting happens.")
		case pickaxeScore < rockScore:
			p.SendMessage("your " + pickaxe.Name() + " almost smashes the " + s.R.Name() + " to bits. you carefully replace the rock and prepare for another attempt.")
		case pickaxeScore < rockScore*3/4:
			p.SendMessage("your " + pickaxe.Name() + " just barely makes it through the " + s.R.Name() + ".")
		case pickaxeScore < rockScore*2:
			p.SendMessage("your " + pickaxe.Name() + " smashes the " + s.R.Name() + " with little difficulty.")
		case pickaxeScore > rockScore*1000:
			p.SendMessage("your " + pickaxe.Name() + " smashes the " + s.R.Name() + " like an atomic bomb on an egg.")
		default:
			p.SendMessage("your " + pickaxe.Name() + " smashes the " + s.R.Name() + " like a hammer on an egg.")
		}
	}
	if rockScore <= pickaxeScore {
		if z.Tile(s.X, s.Y).Remove(s.R) {
			z.Unlock()
			z.Repaint()
			h.Lock()
			if s.Mine {
				h.GiveItem(&Ore{Type: s.R.Ore})
			} else {
				h.GiveItem(&Stone{Type: s.R.Type})
			}
			h.Unlock()
			if p != nil {
				p.Repaint()
			}
			return false
		}
	}
	z.Unlock()

	return false
}

func (s *MineQuarrySchedule) NextMove(x, y uint8) (uint8, uint8) {
	return x, y
}
