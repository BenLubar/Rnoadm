package main

import (
	"math"
)

type WoodType uint8

const (
	Oak WoodType = iota
	Beonetwon

	woodTypeCount
)

var woodTypeInfo = [woodTypeCount]struct {
	Name     string
	Color    Color
	Strength uint64
	sqrtStr  uint64
}{
	Oak: {
		Name:     "oak",
		Color:    "#dab583",
		Strength: 10,
	},
	Beonetwon: {
		Name:     "beonetwon",
		Color:    "#b1b2b0",
		Strength: 1 << 63,
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

func (t *Tree) Paint() (rune, Color) {
	return '♣', "green"
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

func (l *Logs) Paint() (rune, Color) {
	return '➬', woodTypeInfo[l.Type].Color
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
