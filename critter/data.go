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
}{
	SlimupationMage: {
		title:  "mage",
		flavor: "with magical properties.",
	},
	SlimupationBrute: {
		title:  "brute",
		flavor: "that can bench press a whole hero.",
	},
	SlimupationSticky: {
		title:  "citizen",
		flavor: "that's a little stickier than most.",
	},
}

func (s Slimupation) Name() string   { return slimupationInfo[s].title }
func (s Slimupation) Flavor() string { return slimupationInfo[s].flavor }
