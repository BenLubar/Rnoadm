package main

import (
	"math/rand"
)

type HeroName struct {
	FirstT NameSubtype
	First  uint16

	Nickname string

	Last1T NameSubtype
	Last1  uint16
	Last2T NameSubtype
	Last2  uint16
	Last3T NameSubtype
	Last3  uint16
}

func (n *HeroName) Name() string {
	if n == nil {
		return "unknown"
	}

	var buf []byte

	if name := names[n.FirstT][n.First]; name != "" {
		buf = append(append(buf, ' '), name...)
	}

	if n.Nickname != "" {
		buf = append(append(append(buf, ' ', '"'), n.Nickname...), '"')
	}

	if name := names[n.Last1T][n.Last1]; name != "" {
		buf = append(append(append(append(buf, ' '), name...), names[n.Last2T][n.Last2]...), names[n.Last3T][n.Last3]...)
	}

	if len(buf) <= 1 {
		return "unknown"
	}
	return string(buf[1:])
}

func GenerateHumanName(r *rand.Rand, gender Gender) *HeroName {
	n := &HeroName{}
	switch gender {
	case Male:
		n.FirstT = NameMaleHuman
	case Female:
		n.FirstT = NameFemaleHuman
	}
	n.First = uint16(r.Intn(len(names[n.FirstT])))

	if r.Intn(100) < 3 {
		return n
	}

	switch c := r.Intn(100); {
	case c < 25:
		n.Last1T = NameMaleHuman
		for {
			n.Last1 = uint16(r.Intn(len(names[n.Last1T])))
			name := names[n.Last1T][n.Last1]
			if len(name) >= 2 && name[len(name)-2:] == "ss" {
				continue
			}
			if len(name) >= 3 && name[len(name)-3:] == "son" {
				continue
			}
			break
		}
		n.Last2T = NameUtil
		n.Last2 = nameutilSon
	case c < 30:
		switch r.Intn(4) {
		case 0, 1:
			n.Last1T = NameSurnameConsonant
		case 2:
			n.Last2T = NameUtil
			n.Last2 = nameutilR
			fallthrough
		case 3:
			n.Last1T = NameSurnameVowel
		}
		n.Last1 = uint16(r.Intn(len(names[n.Last1T])))
	default:
		switch r.Intn(3) {
		case 0:
			n.Last1T = NameSurnameConsonant
			switch r.Intn(3) {
			case 0, 1:
				n.Last2T = NameEndConsonant
			case 2:
				n.Last2T = NameEndVowel
			}
		case 1:
			n.Last1T = NameSurnameVowel
			n.Last2T = NameEndConsonant
		case 2:
			n.Last1T = NameFrontSyllableConsonant
			switch r.Intn(3) {
			case 0, 1:
				n.Last2T = NameEndConsonant
			case 2:
				n.Last2T = NameEndVowel
			}
		}
		n.Last1 = uint16(r.Intn(len(names[n.Last1T])))
		n.Last2 = uint16(r.Intn(len(names[n.Last2T])))
	}

	return n
}

type ZoneName struct {
	The        bool
	Possessive bool
	Of         bool
	Hero       *HeroName

	DescriptorT NameSubtype
	Descriptor  uint16

	BiomeT NameSubtype
	Biome  uint16
}

func (n *ZoneName) Name() string {
	if n == nil {
		return "unknown zone"
	}

	var buf []byte

	if n.The {
		buf = append(buf, " The"...)
	}

	if n.Possessive {
		buf = append(append(append(buf, ' '), n.Hero.Name()...), "'s"...)
	}

	if name := names[n.DescriptorT][n.Descriptor]; name != "" {
		buf = append(append(buf, ' '), name...)
	}

	if !n.Possessive && !n.Of && n.Hero != nil {
		buf = append(append(buf, ' '), n.Hero.Name()...)
	}

	if name := names[n.BiomeT][n.Biome]; name != "" {
		buf = append(append(buf, ' '), name...)
	}

	if n.Of {
		buf = append(append(buf, " of "...), n.Hero.Name()...)
	}

	if len(buf) <= 1 {
		return "unknown zone"
	}
	return string(buf[1:])
}
