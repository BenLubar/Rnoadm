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
	Name  string
	Color termbox.Attribute
}{
	Coal: {
		Name:  "coal",
		Color: termbox.ColorBlue,
	},
	Iron: {
		Name:  "iron",
		Color: termbox.ColorRed,
	},
	Granite: {
		Name:  "granite",
		Color: termbox.ColorCyan,
	},
	Quartz: {
		Name:  "quartz",
		Color: termbox.ColorWhite,
	},
	Limestone: {
		Name:  "limestone",
		Color: termbox.ColorCyan,
	},
	Sandstone: {
		Name:  "sandstone",
		Color: termbox.ColorYellow,
	},
	Obsidian: {
		Name:  "obsidian",
		Color: termbox.ColorBlue,
	},
	Diamond: {
		Name:  "diamond",
		Color: termbox.ColorWhite,
	},
	Plastic: {
		Name:  "plastic",
		Color: termbox.ColorMagenta,
	},
	Empty: {
		Name:  "empty",
		Color: termbox.ColorBlue,
	},
	Vorpal: {
		Name:  "vorpal",
		Color: termbox.ColorMagenta,
	},
	Wabe: {
		Name:  "wabe",
		Color: termbox.ColorGreen,
	},
	Molten: {
		Name:  "molten",
		Color: termbox.ColorRed,
	},
	Sand: {
		Name:  "sand",
		Color: termbox.ColorYellow,
	},
	Carbonite: {
		Name:  "carbonite",
		Color: termbox.ColorBlue,
	},
	Helium: {
		Name:  "helium",
		Color: termbox.ColorWhite,
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
	return '◊', rockTypeInfo[r.Type].Color
}

func (r *Rock) Blocking() bool {
	return true
}
