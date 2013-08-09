package main

import (
	"math"
)

type WoodType uint8

const (
	Oak WoodType = iota
	Beonetwon
	DeadTree
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

var woodTypeInfo = [woodTypeCount]struct {
	Name      string
	Color     Color
	LeafColor Color
	Strength  uint64
	lowStr    uint64
	sqrtStr   uint64
}{
	Oak: {
		Name:      "oak",
		Color:     "#dab583",
		LeafColor: "#919a2a",
		Strength:  50,
	},
	Beonetwon: {
		Name:      "beonetwon",
		Color:     "#00b120",
		LeafColor: "#b120ee",
		Strength:  1 << 62,
	},
	DeadTree: {
		Name:     "rotting",
		Color:    "#5f5143",
		Strength: 50,
	},
	Maple: {
		Name:      "maple",
		Color:     "#ffb963",
		LeafColor: "#aa5217",
		Strength:  50,
	},
	Birch: {
		Name:      "birch",
		Color:     "#d0ddd0",
		LeafColor: "#29995c",
		Strength:  50,
	},
	Willow: {
		Name:      "willow",
		Color:     "#9e9067",
		LeafColor: "#4e6b2c",
		Strength:  50,
	},
	Juniper: {
		Name:      "juniper",
		Color:     "#c2b19a",
		LeafColor: "#3e4506",
		Strength:  50,
	},
	Wood0: {
		Name:     "wood0",
		Color:    "#000",
		Strength: 5,
	},
	Wood1: {
		Name:     "wood1",
		Color:    "#111",
		Strength: 20,
	},
	Wood2: {
		Name:     "wood2",
		Color:    "#222",
		Strength: 80,
	},
	Wood3: {
		Name:     "wood3",
		Color:    "#333",
		Strength: 300,
	},
	Wood4: {
		Name:     "wood4",
		Color:    "#444",
		Strength: 1000,
	},
	Wood5: {
		Name:     "wood5",
		Color:    "#555",
		Strength: 5000,
	},
	Wood6: {
		Name:     "wood6",
		Color:    "#666",
		Strength: 20000,
	},
	Wood7: {
		Name:     "wood7",
		Color:    "#777",
		Strength: 80000,
	},
	Wood8: {
		Name:     "wood8",
		Color:    "#888",
		Strength: 300000,
	},
	Wood9: {
		Name:     "wood9",
		Color:    "#999",
		Strength: 1000000,
	},
	Wood10: {
		Name:     "wood10",
		Color:    "#aaa",
		Strength: 5000000,
	},
	Wood11: {
		Name:     "wood11",
		Color:    "#bbb",
		Strength: 20000000,
	},
	Wood12: {
		Name:     "wood12",
		Color:    "#ccc",
		Strength: 80000000,
	},
	Wood13: {
		Name:     "wood13",
		Color:    "#ddd",
		Strength: 300000000,
	},
	Wood14: {
		Name:     "wood14",
		Color:    "#eee",
		Strength: 1000000000,
	},
	Wood15: {
		Name:     "wood15",
		Color:    "#fff",
		Strength: 5000000000,
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
	colors := []Color{woodTypeInfo[t.Type].Color}
	if leaf := woodTypeInfo[t.Type].LeafColor; leaf != "" {
		colors = append(colors, leaf)
	}
	return &NetworkedObject{
		Name:    t.Name(),
		Options: []string{"chop down"},
		Sprite:  "tree",
		Colors:  colors,
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
	return "a hatchet made from " + metalTypeInfo[h.Head].Name + " and " + woodTypeInfo[h.Handle].Name + "."
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
		player.Lock()
		player.Equip(h, true)
		player.Unlock()
	}
}

func (h *Hatchet) IsItem() {}

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
