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
		switch r.Intn(2) {
		case 0:
			n.Last1T = NameSurnameConsonant
			switch r.Intn(3) {
			case 0:
				n.Last2T = NameEndConsonant
			case 1:
				n.Last2T = NameEndConsonantMaybeE
			case 2:
				n.Last2T = NameEndVowel
			}
		case 1:
			n.Last1T = NameSurnameVowel
			switch r.Intn(2) {
			case 0:
				n.Last2T = NameEndConsonant
			case 1:
				n.Last2T = NameEndConsonantMaybeE
			}
		}
		n.Last1 = uint16(r.Intn(len(names[n.Last1T])))
		n.Last2 = uint16(r.Intn(len(names[n.Last2T])))
		if n.Last2T == NameEndConsonantMaybeE && r.Intn(2) == 0 {
			n.Last3T = NameUtil
			n.Last3 = nameutilE
		}
	}

	return n
}
