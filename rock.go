package main

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
	Color Color
}{
	Coal: {
		Name:  "coal",
		Color: "blue",
	},
	Iron: {
		Name:  "iron",
		Color: "red",
	},
	Granite: {
		Name:  "granite",
		Color: "cyan",
	},
	Quartz: {
		Name:  "quartz",
		Color: "white",
	},
	Limestone: {
		Name:  "limestone",
		Color: "cyan",
	},
	Sandstone: {
		Name:  "sandstone",
		Color: "yellow",
	},
	Obsidian: {
		Name:  "obsidian",
		Color: "blue",
	},
	Diamond: {
		Name:  "diamond",
		Color: "white",
	},
	Plastic: {
		Name:  "plastic",
		Color: "magenta",
	},
	Empty: {
		Name:  "empty",
		Color: "blue",
	},
	Vorpal: {
		Name:  "vorpal",
		Color: "magenta",
	},
	Wabe: {
		Name:  "wabe",
		Color: "green",
	},
	Molten: {
		Name:  "molten",
		Color: "red",
	},
	Sand: {
		Name:  "sand",
		Color: "yellow",
	},
	Carbonite: {
		Name:  "carbonite",
		Color: "blue",
	},
	Helium: {
		Name:  "helium",
		Color: "white",
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

func (r *Rock) Paint() (rune, Color) {
	return '◊', rockTypeInfo[r.Type].Color
}

func (r *Rock) Blocking() bool {
	return true
}
