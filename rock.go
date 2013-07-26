package main

import (
	"github.com/nsf/termbox-go"
)

type RockType uint8

const (
	Coal RockType = iota
	Iron
	Granite
	Quartz
	Limestone
	Sandstone
	Obsidian
	Diamond
	Plastic
	Empty
	Vorpal
	Wabe
	Molten
	Sand
	Carbonite
	Helium

	rockTypeCount
)

var rockTypeInfo = [rockTypeCount]struct {
	Name string
}{
	Coal: {
		Name: "coal",
	},
	Iron: {
		Name: "iron",
	},
	Granite: {
		Name: "granite",
	},
	Quartz: {
		Name: "quartz",
	},
	Limestone: {
		Name: "limestone",
	},
	Sandstone: {
		Name: "sandstone",
	},
	Obsidian: {
		Name: "obsidian",
	},
	Diamond: {
		Name: "diamond",
	},
	Plastic: {
		Name: "plastic",
	},
	Empty: {
		Name: "empty",
	},
	Vorpal: {
		Name: "vorpal",
	},
	Wabe: {
		Name: "wabe",
	},
	Molten: {
		Name: "molten",
	},
	Sand: {
		Name: "sand",
	},
	Carbonite: {
		Name: "carbonite",
	},
	Helium: {
		Name: "helium",
	},
}

type Rock struct {
	Type RockType
}

func (r *Rock) Name() string {
	return rockTypeInfo[r.Type].Name + " rock"
}

func (r *Rock) Examine() string {
	return "a rock containing " + rockTypeInfo[r.Type].Name + " ore."
}

func (r *Rock) Paint() (rune, termbox.Attribute) {
	return '\u25B2', termbox.ColorWhite
}

func (r *Rock) Blocking() bool {
	return true
}
