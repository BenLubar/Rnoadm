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
		Color:     "#b1b2b0",
		LeafColor: "#c0ffee",
		Strength:  1 << 63,
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
		Color:     "#f5ead5",
		LeafColor: "#899a8c",
		Strength:  12,
	},
}

func init() {
	for t := range woodTypeInfo {
		if woodTypeInfo[t].Strength > 1<<63-1 {
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
	setcell(x, y, "", "tree_trunk_l0", woodTypeInfo[t.Type].Color)
	if color := woodTypeInfo[t.Type].LeafColor; color != "" {
		setcell(x, y, "", "tree_leaves_l1", color)
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
	setcell(x, y, "➬", "", woodTypeInfo[l.Type].Color)
}

func (l *Logs) Blocking() bool {
	return false
}

func (l *Logs) InteractOptions() []string {
	return nil
}

func (l *Logs) IsItem() {}

func (l *Logs) AdminOnly() bool {
	return woodTypeInfo[l.Type].Strength > 1<<63-1
}
