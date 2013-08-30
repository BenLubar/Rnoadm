package hero

import (
	"math/rand"
)

type HeroName struct {
	FirstT NameSubtype
	First  uint64

	Nickname string

	Last1T NameSubtype
	Last1  uint64
	Last2T NameSubtype
	Last2  uint64
	Last3T NameSubtype
	Last3  uint64

	PrefixT NameSubtype
	Prefix  uint64

	SuffixT NameSubtype
	Suffix  uint64
}

func (n *HeroName) serialize() map[string]interface{} {
	return map[string]interface{}{
		"v":   uint(1),
		"ft":  uint64(n.FirstT),
		"f":   n.First,
		"n":   n.Nickname,
		"l1t": uint64(n.Last1T),
		"l1":  n.Last1,
		"l2t": uint64(n.Last2T),
		"l2":  n.Last2,
		"l3t": uint64(n.Last3T),
		"l3":  n.Last3,
		"pt":  uint64(n.PrefixT),
		"p":   n.Prefix,
		"st":  uint64(n.SuffixT),
		"s":   n.Suffix,
	}
}

func (n *HeroName) unserialize(data map[string]interface{}) {
	version := data["v"].(uint)

	switch version {
	case 1:
		n.Prefix = data["p"].(uint64)
		n.PrefixT = NameSubtype(data["pt"].(uint64))
		n.Suffix = data["s"].(uint64)
		n.SuffixT = NameSubtype(data["st"].(uint64))
		fallthrough
	case 0:
		n.First = data["f"].(uint64)
		n.FirstT = NameSubtype(data["ft"].(uint64))
		n.Nickname = data["n"].(string)
		n.Last1 = data["l1"].(uint64)
		n.Last1T = NameSubtype(data["l1t"].(uint64))
		n.Last2 = data["l2"].(uint64)
		n.Last2T = NameSubtype(data["l2t"].(uint64))
		n.Last3 = data["l3"].(uint64)
		n.Last3T = NameSubtype(data["l3t"].(uint64))
	}
}

// Name converts a numeric HeroName into a human-readable string. If Name would
// return an empty string, it instead returns the string "unknown".
func (n *HeroName) Name() string {
	if n == nil {
		return "unknown"
	}

	var buf []byte

	if name := names[n.FirstT][n.First]; name != "" {
		buf = append(append(buf, ' '), name...)
	}

	if n.Nickname != "" {
		buf = append(append(append(buf, " “"...), n.Nickname...), "”"...)
	}

	if name := names[n.Last1T][n.Last1]; name != "" {
		buf = append(append(append(append(buf, ' '), name...), names[n.Last2T][n.Last2]...), names[n.Last3T][n.Last3]...)
	}

	if len(buf) <= 1 {
		return "unknown"
	}
	return names[n.PrefixT][n.Prefix] + string(buf[1:]) + names[n.SuffixT][n.Suffix]
}

// GenerateHumanName randomly generates a name suitable for a human. The first
// name is chosen from a gender-specific list. There is a 3% chance of no
// surname. For the other 97% of names, 25% have a surname beginning with
// a male first name and ending with "son". An additional 5% have a surname
// of a single English word, such as "smith" or "cook". The remaining 70% have
// a surname generated from two syllables.
func GenerateHumanName(r *rand.Rand, gender Gender) HeroName {
	var n HeroName
	switch gender {
	case GenderMale:
		n.FirstT = NameMaleHuman
	case GenderFemale:
		n.FirstT = NameFemaleHuman
	}
	n.First = uint64(r.Intn(len(names[n.FirstT])))

	if r.Intn(100) < 3 {
		return n
	}

	switch c := r.Intn(100); {
	case c < 25:
		n.Last1T = NameMaleHuman
		for {
			n.Last1 = uint64(r.Intn(len(names[n.Last1T])))
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
		n.Last1 = uint64(r.Intn(len(names[n.Last1T])))
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
		n.Last1 = uint64(r.Intn(len(names[n.Last1T])))
		n.Last2 = uint64(r.Intn(len(names[n.Last2T])))
	}

	return n
}
