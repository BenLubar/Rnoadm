package main

import (
	"fmt"
	"math"
)

type WoodType uint8

const (
	Oak WoodType = iota
	Beonetwon
	Maple
	Birch
	Willow
	Juniper
	Wood0
	Wood1
	Wood2
	Wood3
	Wood4
	Wood5
	Wood6
	Wood7
	Wood8
	Wood9
	Wood10
	Wood11
	Wood12
	Wood13
	Wood14
	Wood15

	woodTypeCount
)

var woodTypeInfo = [woodTypeCount]resourceInfo{
	Oak: {
		Article:    "an ",
		Name:       "oak",
		Color:      "#dab583",
		ExtraColor: "#919a2a",
		Strength:   50,
		Density:    65, // Source: Wolfram|Alpha - 2013-09-08
	},
	Beonetwon: {
		Article:    "a ",
		Name:       "beonetwon",
		Color:      "#00b120",
		ExtraColor: "#b120ee",
		Strength:   1 << 62,
		Density:    1,
	},
	DeadTree: {
		Article:  "a ",
		Name:     "rotting",
		Color:    "#5f5143",
		Strength: 50,
		Density:  30,
	},
	Maple: {
		Article:    "a ",
		Name:       "maple",
		Color:      "#ffb963",
		ExtraColor: "#aa5217",
		Strength:   50,
		Density:    60, // Source: Wolfram|Alpha - 2013-09-08
	},
	Birch: {
		Article:    "a ",
		Name:       "birch",
		Color:      "#d0ddd0",
		ExtraColor: "#29995c",
		Strength:   50,
		Density:    64, // Source: Wolfram|Alpha - 2013-09-08
	},
	Willow: {
		Article:    "a ",
		Name:       "willow",
		Color:      "#9e9067",
		ExtraColor: "#4e6b2c",
		Strength:   50,
		Density:    42, // Source: Wolfram|Alpha - 2013-09-08
	},
	Juniper: {
		Article:    "a ",
		Name:       "juniper",
		Color:      "#c2b19a",
		ExtraColor: "#3e4506",
		Strength:   50,
		Density:    39, // Source: Wolfram|Alpha - 2013-09-08
	},
	Wood0: {
		Article:  "a ",
		Name:     "wood0",
		Color:    "#000",
		Strength: 5,
		Density:  55,
	},
	Wood1: {
		Article:  "a ",
		Name:     "wood1",
		Color:    "#111",
		Strength: 20,
		Density:  56,
	},
	Wood2: {
		Article:  "a ",
		Name:     "wood2",
		Color:    "#222",
		Strength: 80,
		Density:  57,
	},
	Wood3: {
		Article:  "a ",
		Name:     "wood3",
		Color:    "#333",
		Strength: 300,
		Density:  58,
	},
	Wood4: {
		Article:  "a ",
		Name:     "wood4",
		Color:    "#444",
		Strength: 1000,
		Density:  59,
	},
	Wood5: {
		Article:  "a ",
		Name:     "wood5",
		Color:    "#555",
		Strength: 5000,
		Density:  60,
	},
	Wood6: {
		Article:  "a ",
		Name:     "wood6",
		Color:    "#666",
		Strength: 20000,
		Density:  61,
	},
	Wood7: {
		Article:  "a ",
		Name:     "wood7",
		Color:    "#777",
		Strength: 80000,
		Density:  62,
	},
	Wood8: {
		Article:  "a ",
		Name:     "wood8",
		Color:    "#888",
		Strength: 300000,
		Density:  63,
	},
	Wood9: {
		Article:  "a ",
		Name:     "wood9",
		Color:    "#999",
		Strength: 1000000,
		Density:  64,
	},
	Wood10: {
		Article:  "a ",
		Name:     "wood10",
		Color:    "#aaa",
		Strength: 5000000,
		Density:  65,
	},
	Wood11: {
		Article:  "a ",
		Name:     "wood11",
		Color:    "#bbb",
		Strength: 20000000,
		Density:  66,
	},
	Wood12: {
		Article:  "a ",
		Name:     "wood12",
		Color:    "#ccc",
		Strength: 80000000,
		Density:  67,
	},
	Wood13: {
		Article:  "a ",
		Name:     "wood13",
		Color:    "#ddd",
		Strength: 300000000,
		Density:  68,
	},
	Wood14: {
		Article:  "a ",
		Name:     "wood14",
		Color:    "#eee",
		Strength: 1000000000,
		Density:  69,
	},
	Wood15: {
		Article:  "a ",
		Name:     "wood15",
		Color:    "#fff",
		Strength: 5000000000,
		Density:  70,
	},
}

func init() {
	for t := range woodTypeInfo {
		woodTypeInfo[t].sqrtStr = uint64(math.Sqrt(float64(woodTypeInfo[t].Strength)))
		if woodTypeInfo[t].Strength >= 1<<60 {
			woodTypeInfo[t].lowStr = woodTypeInfo[t].Strength - 1
		} else {
			woodTypeInfo[t].lowStr = woodTypeInfo[t].sqrtStr
		}
	}
}

type Tree struct {
	networkID
	Type WoodType
}

func (t *Tree) Name() string {
	return woodTypeInfo[t.Type].Name + " tree"
}

func (t *Tree) Examine() string {
	return "a tall " + woodTypeInfo[t.Type].Name + " tree."
}

func (t *Tree) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:    t.Name(),
		Options: []string{"chop down"},
		Sprite:  "tree",
		Colors:  []Color{woodTypeInfo[t.Type].Color, woodTypeInfo[t.Type].ExtraColor},
	}
}

func (t *Tree) Blocking() bool {
	return true
}

func (t *Tree) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // chop down
		player.Lock()
		var schedule Schedule = &ChopTreeSchedule{X: x, Y: y, T: t}
		if tx, ty := player.TileX, player.TileY; (tx-x)*(tx-x)+(ty-y)*(ty-y) > 1 {
			moveSchedule := MoveSchedule(FindPath(zone, tx, ty, x, y, false))
			schedule = &ScheduleSchedule{&moveSchedule, schedule}
		}
		player.schedule = schedule
		player.Unlock()
	}
}

func (t *Tree) ZIndex() int {
	return 0
}

type Logs struct {
	networkID
	Type WoodType
}

func (l *Logs) Name() string {
	return woodTypeInfo[l.Type].Name + " logs"
}

func (l *Logs) Examine() string {
	return "some " + woodTypeInfo[l.Type].Name + " logs."
}

func (l *Logs) Blocking() bool {
	return false
}

func (l *Logs) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:   l.Name(),
		Sprite: "item_logs",
		Colors: []Color{woodTypeInfo[l.Type].Color},
		Item:   true,
	}
}

func (l *Logs) Volume() uint64 {
	return 25
}

func (l *Logs) Weight() uint64 {
	return l.Volume() * woodTypeInfo[l.Type].Density / 100
}

func (l *Logs) AdminOnly() bool {
	return woodTypeInfo[l.Type].Strength >= 1<<60
}

func (l *Logs) ZIndex() int {
	return 25
}

type Hatchet struct {
	networkID
	Head   MetalType
	Handle WoodType
}

func (h *Hatchet) Name() string {
	return metalTypeInfo[h.Head].Name + " hatchet"
}

func (h *Hatchet) Examine() string {
	return fmt.Sprintf("a hatchet made from %s and %s.\nscore: %d - %d", metalTypeInfo[h.Head].Name, woodTypeInfo[h.Handle].Name, metalTypeInfo[h.Head].lowStr+woodTypeInfo[h.Handle].lowStr, metalTypeInfo[h.Head].Strength+woodTypeInfo[h.Handle].Strength)
}

func (h *Hatchet) Blocking() bool {
	return false
}

func (h *Hatchet) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:    h.Name(),
		Sprite:  "item_tools",
		Colors:  []Color{woodTypeInfo[h.Handle].Color, "", metalTypeInfo[h.Head].Color},
		Options: []string{"add to toolbelt"},
		Item:    true,
	}
}

func (h *Hatchet) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // add to toolbelt
		player.Equip(h, true)
	}
}

func (h *Hatchet) Volume() uint64 {
	return 20 + 20
}

func (h *Hatchet) Weight() uint64 {
	return (20*metalTypeInfo[h.Head].Density + 20*woodTypeInfo[h.Handle].Density) / 100
}

func (h *Hatchet) AdminOnly() bool {
	return metalTypeInfo[h.Head].Strength >= 1<<60 || woodTypeInfo[h.Handle].Strength >= 1<<60
}

func (h *Hatchet) ZIndex() int {
	return 25
}

type ChopTreeSchedule struct {
	Delayed bool
	X, Y    uint8
	T       *Tree
}

func (s *ChopTreeSchedule) Act(z *Zone, x uint8, y uint8, h *Hero, p *Player) bool {
	if !s.Delayed {
		s.Delayed = true
		h.scheduleDelay = 10
		if p != nil {
			p.SendMessage("you attempt to cut the " + s.T.Name() + " down.")
		}
		return true
	}
	if (s.X-x)*(s.X-x)+(s.Y-y)*(s.Y-y) > 1 {
		if p != nil {
			p.SendMessage("that is too far away!")
		}
		return false
	}

	h.Lock()
	h.Delay++
	hatchet := h.Toolbelt.Hatchet
	h.Unlock()
	if hatchet == nil {
		if p != nil {
			p.SendMessage("you do not have a hatchet on your toolbelt.")
		}
		return false
	}

	hatchetMax := metalTypeInfo[hatchet.Head].Strength + woodTypeInfo[hatchet.Handle].Strength
	hatchetMin := metalTypeInfo[hatchet.Head].lowStr + woodTypeInfo[hatchet.Handle].lowStr

	treeMax := woodTypeInfo[s.T.Type].Strength
	treeMin := woodTypeInfo[s.T.Type].lowStr

	z.Lock()
	r := z.Rand()
	hatchetScore := uint64(r.Int63n(int64(hatchetMax-hatchetMin+1))) + hatchetMin
	treeScore := uint64(r.Int63n(int64(treeMax-treeMin+1))) + treeMin

	if hatchetScore < treeScore && r.Int63n(int64(treeScore-hatchetScore)) == 0 {
		hatchetScore = treeScore
	}

	if p != nil {
		switch {
		case hatchetScore < treeScore/5:
			p.SendMessage("your " + hatchet.Name() + " doesn't even make a dent in the " + s.T.Name() + ".")
		case hatchetScore < treeScore*2/3:
			p.SendMessage("your " + hatchet.Name() + " slightly dents the " + s.T.Name() + ", but nothing interesting happens.")
		case hatchetScore < treeScore:
			p.SendMessage("your " + hatchet.Name() + " almost chops the " + s.T.Name() + " to the ground. you carefully replace the tree and prepare for another attempt.")
		case hatchetScore < treeScore*3/4:
			p.SendMessage("your " + hatchet.Name() + " just barely makes it through the " + s.T.Name() + ".")
		case hatchetScore < treeScore*2:
			p.SendMessage("your " + hatchet.Name() + " fells the " + s.T.Name() + " with little difficulty.")
		case hatchetScore > treeScore*1000:
			p.SendMessage("your " + hatchet.Name() + " slices through the " + s.T.Name() + " like a chainsaw through butter.")
		default:
			p.SendMessage("your " + hatchet.Name() + " slices through the " + s.T.Name() + " like a knife through butter.")
		}
	}
	if treeScore <= hatchetScore {
		if z.Tile(s.X, s.Y).Remove(s.T) {
			z.Unlock()
			h.Lock()
			success := h.GiveItem(&Logs{Type: s.T.Type})
			h.Unlock()
			if success {
				SendZoneTileChange(z.X, z.Y, TileChange{
					ID:      s.T.NetworkID(),
					Removed: true,
				})
			} else {
				z.Lock()
				z.Tile(s.X, s.Y).Add(s.T)
				z.Unlock()
			}
			return false
		}
	}
	z.Unlock()

	return false
}

func (s *ChopTreeSchedule) NextMove(x, y uint8) (uint8, uint8) {
	return x, y
}
