package hero

var equippables = [equippableSlotCount][]struct {
	name   string
	sprite string
	colors []string
}{
	SlotHead: {
		{
			name:   "novelty foam chicken hat",
			sprite: "hat_novelty_foam_chicken",
			colors: []string{"#fff", "#fd0", "#f00"},
		},
	},
	SlotShirt: {
		{
			name:   "plain shirt",
			sprite: "shirt_basic",
			colors: []string{"#986"},
		},
	},
	SlotPants: {
		{
			name:   "plain pants",
			sprite: "pants_basic",
			colors: []string{"#754"},
		},
	},
	SlotFeet: {
		{
			name:   "plain shoes",
			sprite: "shoes_basic",
			colors: []string{"#cec"},
		},
	},
}
