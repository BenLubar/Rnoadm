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
	radius int
}{
	SlimupationMage: {
		title:  "slime mage",
		flavor: "a slime with some magical abilities.",
		radius: 7,
	},
	SlimupationBrute: {
		title:  "slime brute",
		flavor: "a slime that can bench press a whole hero.",
		radius: 7,
	},
	SlimupationSticky: {
		title:  "sticky slime",
		flavor: "a slime that's a little stickier than most.",
		radius: 14,
	},
}

func (s Slimupation) Name() string   { return slimupationInfo[s].title }
func (s Slimupation) Flavor() string { return slimupationInfo[s].flavor }
func (s Slimupation) Radius() int    { return slimupationInfo[s].radius }
