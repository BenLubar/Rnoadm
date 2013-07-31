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

	woodTypeCount
)

var woodTypeInfo = [woodTypeCount]struct {
	Name      string
	Color     Color
	LeafColor Color
	Strength  uint64
	sqrtStr   uint64
}{
	Oak: {
		Name:      "oak",
		Color:     "#dab583",
		LeafColor: "#919a2a",
		Strength:  10,
	},
	Beonetwon: {
		Name:      "beonetwon",
		Color:     "#00b120",
		LeafColor: "#b120ee",
		Strength:  1 << 62,
	},
	DeadTree: {
		Name:     "dead",
		Color:    "#975",
		Strength: 5,
	},
	Maple: {
		Name:      "maple",
		Color:     "#ffb963",
		LeafColor: "#aa5217",
		Strength:  15,
	},
	Birch: {
		Name:      "birch",
		Color:     "#d0ddd0",
		LeafColor: "#39ca7c",
		Strength:  12,
	},
}

func init() {
	for t := range woodTypeInfo {
		if woodTypeInfo[t].Strength >= 1<<60 {
			woodTypeInfo[t].sqrtStr = woodTypeInfo[t].Strength - 1
		} else {
			woodTypeInfo[t].sqrtStr = uint64(math.Sqrt(float64(woodTypeInfo[t].Strength)))
		}
	}
}

type Tree struct {
	Type WoodType
}

func (t *Tree) Name() string {
	return woodTypeInfo[t.Type].Name + " tree"
}

func (t *Tree) Examine() string {
	return "a tall " + woodTypeInfo[t.Type].Name + " tree."
}

func (t *Tree) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "tree_trunk", woodTypeInfo[t.Type].Color)
	if color := woodTypeInfo[t.Type].LeafColor; color != "" {
		setcell(x, y, "", "tree_leaves", color)
	}
}

func (t *Tree) Blocking() bool {
	return true
}

func (t *Tree) InteractOptions() []string {
	return []string{"chop down"}
}

type Logs struct {
	Type WoodType
}

func (l *Logs) Name() string {
	return woodTypeInfo[l.Type].Name + " logs"
}

func (l *Logs) Examine() string {
	return "some " + woodTypeInfo[l.Type].Name + " logs."
}

func (l *Logs) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "item_logs", woodTypeInfo[l.Type].Color)
}

func (l *Logs) Blocking() bool {
	return false
}

func (l *Logs) InteractOptions() []string {
	return nil
}

func (l *Logs) IsItem() {}

func (l *Logs) AdminOnly() bool {
	return woodTypeInfo[l.Type].Strength >= 1<<60
}

type Hatchet struct {
	Head   MetalType
	Handle WoodType
}

func (h *Hatchet) Name() string {
	return metalTypeInfo[h.Head].Name + " hatchet"
}

func (h *Hatchet) Examine() string {
	return "a hatchet made from " + metalTypeInfo[h.Head].Name + " and " + woodTypeInfo[h.Handle].Name + "."
}

func (h *Hatchet) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "item_tool_handle", woodTypeInfo[h.Handle].Color)
	setcell(x, y, "", "item_tool_hatchet", metalTypeInfo[h.Head].Color)
}

func (h *Hatchet) Blocking() bool {
	return false
}

func (h *Hatchet) InteractOptions() []string {
	return nil
}

func (h *Hatchet) IsItem() {}

func (h *Hatchet) AdminOnly() bool {
	return metalTypeInfo[h.Head].Strength >= 1<<60 || woodTypeInfo[h.Handle].Strength >= 1<<60
}
