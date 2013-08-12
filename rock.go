package main

import (
	"fmt"
	"math"
)

type resourceInfo struct {
	Article    string
	Name       string
	Color      Color
	ExtraColor Color
	Strength   uint64
	lowStr     uint64
	sqrtStr    uint64
	Density    uint64 // centigrams per cubic centimeter (cg/cc)
}

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

var rockTypeInfo = [rockTypeCount]resourceInfo{
	Granite: {
		Article:  "a ",
		Name:     "granite",
		Color:    "#948e85",
		Strength: 50,
		Density:  260, // Source: Wolfram|Alpha - 2013-09-08
	},
	Adminstone: {
		Article:  "an ",
		Name:     "adminstone",
		Color:    "#1e0036",
		Strength: 1 << 62,
		Density:  15000,
	},
	Limestone: {
		Article:  "a ",
		Name:     "limestone",
		Color:    "#bebd8f",
		Strength: 50,
		Density:  271, // Source: Wolfram|Alpha - 2013-09-08
	},
	Stone0: {
		Article:  "a ",
		Name:     "stone0",
		Color:    "#000",
		Strength: 5,
		Density:  150,
	},
	Stone1: {
		Article:  "a ",
		Name:     "stone1",
		Color:    "#111",
		Strength: 20,
		Density:  160,
	},
	Stone2: {
		Article:  "a ",
		Name:     "stone2",
		Color:    "#222",
		Strength: 80,
		Density:  170,
	},
	Stone3: {
		Article:  "a ",
		Name:     "stone3",
		Color:    "#333",
		Strength: 300,
		Density:  180,
	},
	Stone4: {
		Article:  "a ",
		Name:     "stone4",
		Color:    "#444",
		Strength: 1000,
		Density:  190,
	},
	Stone5: {
		Article:  "a ",
		Name:     "stone5",
		Color:    "#555",
		Strength: 5000,
		Density:  200,
	},
	Stone6: {
		Article:  "a ",
		Name:     "stone6",
		Color:    "#666",
		Strength: 20000,
		Density:  210,
	},
	Stone7: {
		Article:  "a ",
		Name:     "stone7",
		Color:    "#777",
		Strength: 80000,
		Density:  220,
	},
	Stone8: {
		Article:  "a ",
		Name:     "stone8",
		Color:    "#888",
		Strength: 300000,
		Density:  230,
	},
	Stone9: {
		Article:  "a ",
		Name:     "stone9",
		Color:    "#999",
		Strength: 1000000,
		Density:  240,
	},
	Stone10: {
		Article:  "a ",
		Name:     "stone10",
		Color:    "#aaa",
		Strength: 5000000,
		Density:  250,
	},
	Stone11: {
		Article:  "a ",
		Name:     "stone11",
		Color:    "#bbb",
		Strength: 20000000,
		Density:  260,
	},
	Stone12: {
		Article:  "a ",
		Name:     "stone12",
		Color:    "#ccc",
		Strength: 80000000,
		Density:  270,
	},
	Stone13: {
		Article:  "a ",
		Name:     "stone13",
		Color:    "#ddd",
		Strength: 300000000,
		Density:  280,
	},
	Stone14: {
		Article:  "a ",
		Name:     "stone14",
		Color:    "#eee",
		Strength: 1000000000,
		Density:  290,
	},
	Stone15: {
		Article:  "a ",
		Name:     "stone15",
		Color:    "#fff",
		Strength: 5000000000,
		Density:  300,
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
	RustyMetal

	metalTypeCount
)

var metalTypeInfo = [metalTypeCount]resourceInfo{
	Iron: {
		Article:  "an ",
		Name:     "iron",
		Color:    "#79493d",
		Strength: 50,
		Density:  787, // Source: Wolfram|Alpha - 2013-09-08
	},
	Unobtainium: {
		Article:  "an ",
		Name:     "unobtainium",
		Color:    "#cd8aff",
		Strength: 1 << 62,
		Density:  2256,
	},
	Copper: {
		Article:  "a ",
		Name:     "copper",
		Color:    "#af633e",
		Strength: 50,
		Density:  896, // Source: Wolfram|Alpha - 2013-09-08
	},
	Metal0: {
		Article:  "a ",
		Name:     "metal0",
		Color:    "#000",
		Strength: 5,
		Density:  850,
	},
	Metal1: {
		Article:  "a ",
		Name:     "metal1",
		Color:    "#111",
		Strength: 20,
		Density:  860,
	},
	Metal2: {
		Article:  "a ",
		Name:     "metal2",
		Color:    "#222",
		Strength: 80,
		Density:  870,
	},
	Metal3: {
		Article:  "a ",
		Name:     "metal3",
		Color:    "#333",
		Strength: 300,
		Density:  880,
	},
	Metal4: {
		Article:  "a ",
		Name:     "metal4",
		Color:    "#444",
		Strength: 1000,
		Density:  890,
	},
	Metal5: {
		Article:  "a ",
		Name:     "metal5",
		Color:    "#555",
		Strength: 5000,
		Density:  900,
	},
	Metal6: {
		Article:  "a ",
		Name:     "metal6",
		Color:    "#666",
		Strength: 20000,
		Density:  910,
	},
	Metal7: {
		Article:  "a ",
		Name:     "metal7",
		Color:    "#777",
		Strength: 80000,
		Density:  920,
	},
	Metal8: {
		Article:  "a ",
		Name:     "metal8",
		Color:    "#888",
		Strength: 300000,
		Density:  930,
	},
	Metal9: {
		Article:  "a ",
		Name:     "metal9",
		Color:    "#999",
		Strength: 1000000,
		Density:  940,
	},
	Metal10: {
		Article:  "a ",
		Name:     "metal10",
		Color:    "#aaa",
		Strength: 5000000,
		Density:  950,
	},
	Metal11: {
		Article:  "a ",
		Name:     "metal11",
		Color:    "#bbb",
		Strength: 20000000,
		Density:  960,
	},
	Metal12: {
		Article:  "a ",
		Name:     "metal12",
		Color:    "#ccc",
		Strength: 80000000,
		Density:  970,
	},
	Metal13: {
		Article:  "a ",
		Name:     "metal13",
		Color:    "#ddd",
		Strength: 300000000,
		Density:  980,
	},
	Metal14: {
		Article:  "a ",
		Name:     "metal14",
		Color:    "#eee",
		Strength: 1000000000,
		Density:  990,
	},
	Metal15: {
		Article:  "a ",
		Name:     "metal15",
		Color:    "#fff",
		Strength: 5000000000,
		Density:  1000,
	},
	RustyMetal: {
		Article:  "a ",
		Name:     "rusty",
		Color:    "#4c271e",
		Strength: 50,
		Density:  100,
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
			moveSchedule := FindPath(zone, tx, ty, x, y, false)
			schedule = &ScheduleSchedule{moveSchedule, schedule}
		}
		player.schedule = schedule
		player.Unlock()
	case 1: // quarry
		player.Lock()
		var schedule Schedule = &MineQuarrySchedule{X: x, Y: y, R: r, Mine: false}
		if tx, ty := player.TileX, player.TileY; (tx-x)*(tx-x)+(ty-y)*(ty-y) > 1 {
			moveSchedule := FindPath(zone, tx, ty, x, y, false)
			schedule = &ScheduleSchedule{moveSchedule, schedule}
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
	Type    RockType
	Quality uint64
}

func (s *Stone) Name() string {
	return rockTypeInfo[s.Type].Name + " stone"
}

func (s *Stone) Examine() string {
	return "some " + rockTypeInfo[s.Type].Name + " stone."
}

func (s *Stone) Blocking() bool {
	return false
}

func (s *Stone) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:   s.Name(),
		Sprite: "item_stone",
		Colors: []Color{rockTypeInfo[s.Type].Color},
		Item:   true,
	}
}

func (s *Stone) Volume() uint64 {
	return 25
}

func (s *Stone) Weight() uint64 {
	return s.Volume() * rockTypeInfo[s.Type].Density / 100
}

func (s *Stone) AdminOnly() bool {
	return rockTypeInfo[s.Type].Strength >= 1<<60
}

func (s *Stone) ZIndex() int {
	return 25
}

type Ore struct {
	networkID
	Type    MetalType
	Quality uint64
}

func (o *Ore) Name() string {
	return metalTypeInfo[o.Type].Name + " ore"
}

func (o *Ore) Examine() string {
	return "some " + metalTypeInfo[o.Type].Name + " ore."
}

func (o *Ore) Blocking() bool {
	return false
}

func (o *Ore) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:   o.Name(),
		Sprite: "item_ore",
		Colors: []Color{metalTypeInfo[o.Type].Color},
		Item:   true,
	}
}

func (o *Ore) Volume() uint64 {
	return 25
}

func (o *Ore) Weight() uint64 {
	return o.Volume() * metalTypeInfo[o.Type].Density / 100
}

func (o *Ore) AdminOnly() bool {
	return metalTypeInfo[o.Type].Strength >= 1<<60
}

func (o *Ore) ZIndex() int {
	return 25
}

type Pickaxe struct {
	networkID
	Head    MetalType
	Handle  WoodType
	Quality uint64
}

func (p *Pickaxe) Name() string {
	return metalTypeInfo[p.Head].Name + " pickaxe"
}

func (p *Pickaxe) Examine() string {
	return fmt.Sprintf("a pickaxe made from %s metal and %s wood.\nscore: %d - %d", metalTypeInfo[p.Head].Name, woodTypeInfo[p.Handle].Name, metalTypeInfo[p.Head].lowStr+woodTypeInfo[p.Handle].lowStr, metalTypeInfo[p.Head].Strength+woodTypeInfo[p.Handle].Strength)
}

func (p *Pickaxe) Blocking() bool {
	return false
}

func (p *Pickaxe) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:    p.Name(),
		Sprite:  "item_tools",
		Colors:  []Color{woodTypeInfo[p.Handle].Color, metalTypeInfo[p.Head].Color},
		Options: []string{"add to toolbelt"},
		Item:    true,
	}
}

func (p *Pickaxe) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // add to toolbelt
		player.Equip(p, true)
	}
}

func (p *Pickaxe) Volume() uint64 {
	return 20 + 20
}

func (p *Pickaxe) Weight() uint64 {
	return (20*metalTypeInfo[p.Head].Density + 20*woodTypeInfo[p.Handle].Density) / 100
}

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
			h.Lock()
			var success bool
			if s.Mine {
				success = h.GiveItem(&Ore{Type: s.R.Ore})
			} else {
				success = h.GiveItem(&Stone{Type: s.R.Type})
			}
			h.Unlock()
			if success {
				SendZoneTileChange(z.X, z.Y, TileChange{
					ID:      s.R.NetworkID(),
					Removed: true,
				})
			} else {
				z.Lock()
				z.Tile(s.X, s.Y).Add(s.R)
				z.Unlock()
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
