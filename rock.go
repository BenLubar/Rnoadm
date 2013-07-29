package main

type RockType uint16

const (
	Granite RockType = iota

	rockTypeCount
)

var rockTypeInfo = [rockTypeCount]struct {
	Name  string
	Color Color
}{
	Granite: {
		Name:  "granite",
		Color: "#948e85",
	},
}

type MetalType uint16

const (
	_ MetalType = iota
	Iron

	metalTypeCount
)

var metalTypeInfo = [metalTypeCount]struct {
	Name  string
	Color Color
}{
	Iron: {
		Name:  "iron",
		Color: "#79493d",
	},
}

type Rock struct {
	Type RockType
	Ore  MetalType
}

func (r *Rock) Name() string {
	return rockTypeInfo[r.Type].Name + " rock"
}

func (r *Rock) Examine() string {
	if r.Ore != 0 {
		return "a " + rockTypeInfo[r.Type].Name + " rock containing " + metalTypeInfo[r.Ore].Name + " ore."
	}
	return "a " + rockTypeInfo[r.Type].Name + " rock."
}

func (r *Rock) Paint() (rune, Color) {
	return '◊', rockTypeInfo[r.Type].Color
}

func (r *Rock) Blocking() bool {
	return true
}

func (r *Rock) InteractOptions() []string {
	return []string{"mine", "quarry", "prospect"}
}
