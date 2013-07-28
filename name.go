package main

import (
	"math/rand"
)

type Name struct {
	Name             string
	FrontCompoundST  NameSubtype
	FrontCompound    uint16
	RearCompoundST   NameSubtype
	RearCompound     uint16
	Adjective1ST     NameSubtype
	Adjective1       uint16
	Adjective2ST     NameSubtype
	Adjective2       uint16
	HyphenCompoundST NameSubtype
	HyphenCompound   uint16
	TheST            NameSubtype
	The              uint16
	OfST             NameSubtype
	Of               uint16
}

func (n *Name) String() string {
	s := n.Name
	if n.FrontCompound != 0 || n.RearCompound != 0 {
		s += " " + names[n.FrontCompoundST][n.FrontCompound] + names[n.RearCompoundST][n.RearCompound]
	}
	if n.The != 0 {
		s += " the "
		if n.Adjective1 != 0 {
			s += names[n.Adjective1ST][n.Adjective1] + " "
		}
		if n.Adjective2 != 0 {
			s += names[n.Adjective2ST][n.Adjective2] + " "
		}
		if n.HyphenCompound != 0 {
			s += names[n.HyphenCompoundST][n.HyphenCompound] + "-"
		}
		s += names[n.TheST][n.The]
	}
	if n.Of != 0 {
		s += " of " + names[n.OfST][n.Of]
	}
	return s
}

func GenerateName(r *rand.Rand, subtypes ...NameSubtype) *Name {
	n := &Name{
		Name: "TODONAME",
	}
	if len(subtypes) == 0 {
		return n
	}
	if r.Intn(2) == 0 {
		n.FrontCompoundST = subtypes[r.Intn(len(subtypes))]
		n.FrontCompound = uint16(r.Intn(len(names[n.FrontCompoundST])))
	}
	if r.Intn(2) == 0 {
		n.RearCompoundST = subtypes[r.Intn(len(subtypes))]
		n.RearCompound = uint16(r.Intn(len(names[n.RearCompoundST])))
	}
	if r.Intn(2) == 0 {
		n.Adjective1ST = subtypes[r.Intn(len(subtypes))]
		n.Adjective1 = uint16(r.Intn(len(names[n.Adjective1ST])))
	}
	if r.Intn(2) == 0 {
		n.Adjective2ST = subtypes[r.Intn(len(subtypes))]
		n.Adjective2 = uint16(r.Intn(len(names[n.Adjective2ST])))
	}
	if r.Intn(2) == 0 {
		n.HyphenCompoundST = subtypes[r.Intn(len(subtypes))]
		n.HyphenCompound = uint16(r.Intn(len(names[n.HyphenCompoundST])))
	}
	if r.Intn(2) == 0 {
		n.TheST = subtypes[r.Intn(len(subtypes))]
		n.The = uint16(r.Intn(len(names[n.TheST])))
	}
	if r.Intn(2) == 0 {
		n.OfST = subtypes[r.Intn(len(subtypes))]
		n.Of = uint16(r.Intn(len(names[n.OfST])))
	}
	return n
}

type NameSubtype uint16

const (
	NameZone NameSubtype = iota
	NameForest
	NamePlains
	NameHills
	NameLake

	nameSubtypeCount
)

// leave the first one of each subtype blank. add only to the end of each list.
var names = [nameSubtypeCount][]string{
	NameZone: {
		"",
		"area",
		"zone",
		"region",
		"lair",
		"territory",
		"realm",
	},
	NameForest: {
		"",
		"forest",
		"glade",
		"grove",
		"timberland",
		"woodland",
		"weald",
	},
	NamePlains: {
		"",
		"plains",
		"steppe",
		"plateau",
		"prairie",
		"meadow",
		"field",
		"moors",
	},
	NameHills: {
		"",
		"hills",
		"foothills",
		"bluff",
		"ridge",
		"hillocks",
		"knoll",
		"mesa",
		"mound",
	},
	NameLake: {
		"",
		"lake",
		"loch",
		"reservoir",
		"basin",
		"sea",
	},
}
