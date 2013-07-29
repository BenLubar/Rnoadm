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

type Stone struct {
	Type RockType
}

func (s *Stone) Name() string {
	return rockTypeInfo[s.Type].Name + " stone"
}

func (s *Stone) Examine() string {
	return "a " + rockTypeInfo[s.Type].Name + " stone."
}

func (s *Stone) Paint() (rune, Color) {
	return '✧', rockTypeInfo[s.Type].Color
}

func (s *Stone) Blocking() bool {
	return false
}

func (s *Stone) InteractOptions() []string {
	return nil
}

func (s *Stone) IsItem() {}

type Ore struct {
	Type MetalType
}

func (o *Ore) Name() string {
	return metalTypeInfo[o.Type].Name + " ore"
}

func (o *Ore) Examine() string {
	return "some " + metalTypeInfo[o.Type].Name + " ore."
}

func (o *Ore) Paint() (rune, Color) {
	return '❖', metalTypeInfo[o.Type].Color
}

func (o *Ore) Blocking() bool {
	return false
}

func (o *Ore) InteractOptions() []string {
	return nil
}

func (o *Ore) IsItem() {}
