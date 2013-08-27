package critter

type Slimupation uint64

const (
	SlimupationMage Slimupation = iota
	SlimupationBrute
	SlimupationSticky

	slimupationCount
)

var slimupationInfo = [slimupationCount]struct {
	title  string
	flavor string
	color  string
	radius int
}{
	SlimupationMage: {
		title:  "slime mage",
		flavor: "a slime with some magical abilities.",
		color:  "#00d",
		radius: 7,
	},
	SlimupationBrute: {
		title:  "slime brute",
		flavor: "a slime that can bench press a whole hero.",
		color:  "#d00",
		radius: 7,
	},
	SlimupationSticky: {
		title:  "sticky slime",
		flavor: "a slime that's a little stickier than most.",
		color:  "#0d0",
		radius: 14,
	},
}

func (s Slimupation) Name() string   { return slimupationInfo[s].title }
func (s Slimupation) Flavor() string { return slimupationInfo[s].flavor }
