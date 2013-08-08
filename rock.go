package main

import (
	"math"
)

type RockType uint16

const (
	Granite RockType = iota
	Adminstone
	Limestone
	Stone0
	Stone1
	Stone2
	Stone3
	Stone4
	Stone5
	Stone6
	Stone7
	Stone8
	Stone9
	Stone10
	Stone11
	Stone12
	Stone13
	Stone14
	Stone15

	rockTypeCount
)

var rockTypeInfo = [rockTypeCount]struct {
	Name     string
	Color    Color
	Strength uint64
	lowStr   uint64
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
	Stone0: {
		Name:     "stone0",
		Color:    "#000",
		Strength: 5,
	},
	Stone1: {
		Name:     "stone1",
		Color:    "#111",
		Strength: 20,
	},
	Stone2: {
		Name:     "stone2",
		Color:    "#222",
		Strength: 80,
	},
	Stone3: {
		Name:     "stone3",
		Color:    "#333",
		Strength: 300,
	},
	Stone4: {
		Name:     "stone4",
		Color:    "#444",
		Strength: 1000,
	},
	Stone5: {
		Name:     "stone5",
		Color:    "#555",
		Strength: 5000,
	},
	Stone6: {
		Name:     "stone6",
		Color:    "#666",
		Strength: 20000,
	},
	Stone7: {
		Name:     "stone7",
		Color:    "#777",
		Strength: 80000,
	},
	Stone8: {
		Name:     "stone8",
		Color:    "#888",
		Strength: 300000,
	},
	Stone9: {
		Name:     "stone9",
		Color:    "#999",
		Strength: 1000000,
	},
	Stone10: {
		Name:     "stone10",
		Color:    "#aaa",
		Strength: 5000000,
	},
	Stone11: {
		Name:     "stone11",
		Color:    "#bbb",
		Strength: 20000000,
	},
	Stone12: {
		Name:     "stone12",
		Color:    "#ccc",
		Strength: 80000000,
	},
	Stone13: {
		Name:     "stone13",
		Color:    "#ddd",
		Strength: 300000000,
	},
	Stone14: {
		Name:     "stone14",
		Color:    "#eee",
		Strength: 1000000000,
	},
	Stone15: {
		Name:     "stone15",
		Color:    "#fff",
		Strength: 5000000000,
	},
}

func init() {
	for t := range rockTypeInfo {
		rockTypeInfo[t].sqrtStr = uint64(math.Sqrt(float64(rockTypeInfo[t].Strength)))
		if rockTypeInfo[t].Strength >= 1<<60 {
			rockTypeInfo[t].lowStr = rockTypeInfo[t].Strength - 1
		} else {
			rockTypeInfo[t].lowStr = rockTypeInfo[t].sqrtStr
		}
	}
}

type MetalType uint16

const (
	_ MetalType = iota
	Iron
	Unobtainium
	Copper
	Metal0
	Metal1
	Metal2
	Metal3
	Metal4
	Metal5
	Metal6
	Metal7
	Metal8
	Metal9
	Metal10
	Metal11
	Metal12
	Metal13
	Metal14
	Metal15

	metalTypeCount
)

var metalTypeInfo = [metalTypeCount]struct {
	Name     string
	Color    Color
	Strength uint64
	lowStr   uint64
	sqrtStr  uint64
}{
	Iron: {
		Name:     "iron",
		Color:    "#79493d",
		Strength: 50,
	},
	Unobtainium: {
		Name:     "unobtainium",
		Color:    "#cd8aff",
		Strength: 1 << 62,
	},
	Copper: {
		Name:     "copper",
		Color:    "#af633e",
		Strength: 50,
	},
	Metal0: {
		Name:     "metal0",
		Color:    "#000",
		Strength: 5,
	},
	Metal1: {
		Name:     "metal1",
		Color:    "#111",
		Strength: 20,
	},
	Metal2: {
		Name:     "metal2",
		Color:    "#222",
		Strength: 80,
	},
	Metal3: {
		Name:     "metal3",
		Color:    "#333",
		Strength: 300,
	},
	Metal4: {
		Name:     "metal4",
		Color:    "#444",
		Strength: 1000,
	},
	Metal5: {
		Name:     "metal5",
		Color:    "#555",
		Strength: 5000,
	},
	Metal6: {
		Name:     "metal6",
		Color:    "#666",
		Strength: 20000,
	},
	Metal7: {
		Name:     "metal7",
		Color:    "#777",
		Strength: 80000,
	},
	Metal8: {
		Name:     "metal8",
		Color:    "#888",
		Strength: 300000,
	},
	Metal9: {
		Name:     "metal9",
		Color:    "#999",
		Strength: 1000000,
	},
	Metal10: {
		Name:     "metal10",
		Color:    "#aaa",
		Strength: 5000000,
	},
	Metal11: {
		Name:     "metal11",
		Color:    "#bbb",
		Strength: 20000000,
	},
	Metal12: {
		Name:     "metal12",
		Color:    "#ccc",
		Strength: 80000000,
	},
	Metal13: {
		Name:     "metal13",
		Color:    "#ddd",
		Strength: 300000000,
	},
	Metal14: {
		Name:     "metal14",
		Color:    "#eee",
		Strength: 1000000000,
	},
	Metal15: {
		Name:     "metal15",
		Color:    "#fff",
		Strength: 5000000000,
	},
}

func init() {
	for t := range metalTypeInfo {
		metalTypeInfo[t].sqrtStr = uint64(math.Sqrt(float64(metalTypeInfo[t].Strength)))
		if metalTypeInfo[t].Strength >= 1<<60 {
			metalTypeInfo[t].lowStr = metalTypeInfo[t].Strength - 1
		} else {
			metalTypeInfo[t].lowStr = metalTypeInfo[t].sqrtStr
		}
	}
}

type Rock struct {
	networkID
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

func (r *Rock) Serialize() *NetworkedObject {
	colors := []Color{rockTypeInfo[r.Type].Color}
	if r.Ore != 0 {
		colors = append(colors, metalTypeInfo[r.Ore].Color)
		if r.Big {
			colors = append(colors, metalTypeInfo[r.Ore].Color)
		}
	}
	return &NetworkedObject{
		Name:    r.Name(),
		Options: []string{"mine", "quarry"},
		Sprite:  "rock",
		Colors:  colors,
	}
}

func (r *Rock) Blocking() bool {
	return true
}

func (r *Rock) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
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
	networkID
	Type RockType
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
	networkID
	Type MetalType
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
	networkID
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
	pickaxeMin := metalTypeInfo[pickaxe.Head].lowStr + woodTypeInfo[pickaxe.Handle].lowStr

	var rockMax, rockMin uint64
	if s.Mine {
		rockMax = metalTypeInfo[s.R.Ore].Strength
		rockMin = metalTypeInfo[s.R.Ore].lowStr
	} else {
		rockMax = rockTypeInfo[s.R.Type].Strength
		rockMin = rockTypeInfo[s.R.Type].lowStr
	}

	z.Lock()
	r := z.Rand()
	pickaxeScore := uint64(r.Int63n(int64(pickaxeMax-pickaxeMin+1))) + pickaxeMin
	rockScore := uint64(r.Int63n(int64(rockMax-rockMin+1))) + rockMin

	if pickaxeScore < rockScore && r.Int63n(int64(rockScore-pickaxeScore)) == 0 {
		pickaxeScore = rockScore
	}

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
			SendZoneTileChange(z.X, z.Y, TileChange{
				ID:      s.R.NetworkID(),
				Removed: true,
			})
			h.Lock()
			if s.Mine {
				h.GiveItem(&Ore{Type: s.R.Ore})
			} else {
				h.GiveItem(&Stone{Type: s.R.Type})
			}
			h.Unlock()
			return false
		}
	}
	z.Unlock()

	return false
}

func (s *MineQuarrySchedule) NextMove(x, y uint8) (uint8, uint8) {
	return x, y
}
