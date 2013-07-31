package main

import (
	"math"
)

type RockType uint16

const (
	Granite RockType = iota
	Adminstone
	Limestone

	rockTypeCount
)

var rockTypeInfo = [rockTypeCount]struct {
	Name     string
	Color    Color
	Strength uint64
	sqrtStr  uint64
}{
	Granite: {
		Name:     "granite",
		Color:    "#948e85",
		Strength: 50,
	},
	Adminstone: {
		Name:     "adminstone",
		Color:    "#1e0036",
		Strength: 1 << 62,
	},
	Limestone: {
		Name:     "limestone",
		Color:    "#bebd8f",
		Strength: 45,
	},
}

func init() {
	for t := range rockTypeInfo {
		if rockTypeInfo[t].Strength >= 1<<60 {
			rockTypeInfo[t].sqrtStr = rockTypeInfo[t].Strength - 1
		} else {
			rockTypeInfo[t].sqrtStr = uint64(math.Sqrt(float64(rockTypeInfo[t].Strength)))
		}
	}
}

type MetalType uint16

const (
	_ MetalType = iota
	Iron
	Unobtainium
	Copper

	metalTypeCount
)

var metalTypeInfo = [metalTypeCount]struct {
	Name     string
	Color    Color
	Strength uint64
	sqrtStr  uint64
}{
	Iron: {
		Name:     "iron",
		Color:    "#79493d",
		Strength: 100,
	},
	Unobtainium: {
		Name:     "unobtainium",
		Color:    "#ff00ff",
		Strength: 1 << 62,
	},
	Copper: {
		Name:     "copper",
		Color:    "#af633e",
		Strength: 65,
	},
}

func init() {
	for t := range metalTypeInfo {
		if metalTypeInfo[t].Strength >= 1<<60 {
			metalTypeInfo[t].sqrtStr = metalTypeInfo[t].Strength - 1
		} else {
			metalTypeInfo[t].sqrtStr = uint64(math.Sqrt(float64(metalTypeInfo[t].Strength)))
		}
	}
}

type Rock struct {
	Type RockType
	Ore  MetalType
	Big  bool
}

func (r *Rock) Name() string {
	return rockTypeInfo[r.Type].Name + " rock"
}

func (r *Rock) Examine() string {
	if r.Ore != 0 {
		if r.Big {
			return "a deposit of " + rockTypeInfo[r.Type].Name + " rock containing rich " + metalTypeInfo[r.Ore].Name + " ore."
		}
		return "a deposit of " + rockTypeInfo[r.Type].Name + " rock containing " + metalTypeInfo[r.Ore].Name + " ore."
	}
	return "a deposit of " + rockTypeInfo[r.Type].Name + " rock."
}

func (r *Rock) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "rock_base", rockTypeInfo[r.Type].Color)
	if r.Ore != 0 {
		setcell(x, y, "", "rock_ore_small", metalTypeInfo[r.Ore].Color)
		if r.Big {
			setcell(x, y, "", "rock_ore_big", metalTypeInfo[r.Ore].Color)
		}
	}
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

func (s *Stone) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "item_stone", rockTypeInfo[s.Type].Color)
}

func (s *Stone) Blocking() bool {
	return false
}

func (s *Stone) InteractOptions() []string {
	return nil
}

func (s *Stone) IsItem() {}

func (s *Stone) AdminOnly() bool {
	return rockTypeInfo[s.Type].Strength >= 1<<60
}

type Ore struct {
	Type MetalType
}

func (o *Ore) Name() string {
	return metalTypeInfo[o.Type].Name + " ore"
}

func (o *Ore) Examine() string {
	return "some " + metalTypeInfo[o.Type].Name + " ore."
}

func (o *Ore) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "item_ore", metalTypeInfo[o.Type].Color)
}

func (o *Ore) Blocking() bool {
	return false
}

func (o *Ore) InteractOptions() []string {
	return nil
}

func (o *Ore) IsItem() {}

func (o *Ore) AdminOnly() bool {
	return metalTypeInfo[o.Type].Strength >= 1<<60
}

type Pickaxe struct {
	Head MetalType
	Handle WoodType
}

func (p *Pickaxe) Name() string {
	return metalTypeInfo[p.Head].Name + " pickaxe"
}

func (p *Pickaxe) Examine() string {
	return "a pickaxe made from " + metalTypeInfo[p.Head].Name + " and " + woodTypeInfo[p.Handle].Name + "."
}

func (p *Pickaxe) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "item_tool_handle", woodTypeInfo[p.Handle].Color)
	setcell(x, y, "", "item_tool_pickaxe", metalTypeInfo[p.Head].Color)
}

func (p *Pickaxe) Blocking() bool {
	return false
}

func (p *Pickaxe) InteractOptions() []string {
	return nil
}

func (p *Pickaxe) IsItem() {}

func (p *Pickaxe) AdminOnly() bool {
	return metalTypeInfo[p.Head].Strength >= 1<<60 || woodTypeInfo[p.Handle].Strength >= 1<<60
}
