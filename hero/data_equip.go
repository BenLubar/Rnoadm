package hero

var equippables = [equipSlotCount][]struct {
	name   string
	sprite string
	colors []string
	height uint // default 32
	width  uint // default 32
}{
	SlotHead: {
		{
			name:   "novelty foam chicken hat",
			sprite: "hat_novelty_foam_chicken",
			colors: []string{"#fff", "#fd0", "#f00"},
			height: 64,
		},
		{
			name:   "bear riding a unicycle",
			sprite: "hat_bear_riding_a_unicycle",
			colors: []string{"#ccf", "#ecb", "#fff"},
			height: 87,
		},
	},
	SlotShirt: {
		{
			name:   "plain shirt",
			sprite: "shirt_basic",
			colors: []string{"#986"},
			height: 48,
		},
	},
	SlotPants: {
		{
			name:   "plain pants",
			sprite: "pants_basic",
			colors: []string{"#754"},
			height: 48,
		},
	},
	SlotFeet: {
		{
			name:   "plain shoes",
			sprite: "shoes_basic",
			colors: []string{"#cec"},
			height: 48,
		},
	},
}
