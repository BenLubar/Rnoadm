package main

type Name struct {
	Name           string
	FrontCompound  uint16
	RearCompound   uint16
	Adjective1     uint16
	Adjective2     uint16
	HyphenCompound uint16
	The            uint16
	Of             uint16
}

func (n Name) String() string {
	s := n.Name
	if n.FrontCompound != 0 || n.RearCompound != 0 {
		s += " " + names[n.FrontCompound] + names[n.RearCompound]
	}
	if n.The != 0 {
		s += " the "
		if n.Adjective1 != 0 {
			s += names[n.Adjective1] + " "
		}
		if n.Adjective2 != 0 {
			s += names[n.Adjective2] + " "
		}
		if n.HyphenCompound != 0 {
			s += names[n.HyphenCompound] + "-"
		}
		s += names[n.The]
	}
	if n.Of != 0 {
		s += " of " + names[n.Of]
	}
	return s
}

var names = []string{
	"", // leave this one blank. add only to the end of the list.
}
