package hero

var equippables = [equipSlotCount][]struct {
	name   string
	sprite string
	colors []string
	height uint // default 32
	width  uint // default 32

	animationOverrides map[string]string
	adminOnly          bool
}{
	SlotHead: {
		{
			name:   "novelty foam chicken hat",
			sprite: "hat_novelty_foam_chicken",
			colors: []string{"#fff", "#fd0", "#f00"},
			height: 64,

			adminOnly: true,
		},
		{
			name:   "bear riding a unicycle",
			sprite: "hat_bear_riding_a_unicycle",
			colors: []string{"#ccf", "#ecb", "#fff"},
			height: 87,

			adminOnly: true,
		},
		{
			name:   "unusual strange vintage hound dog",
			sprite: "hat_unusual_strange_vintage_hound_dog",
			colors: []string{"#fff", "#fff"},
			height: 32,

			adminOnly: true,
		},
		{
			name:   "crown of the Origin",
			sprite: "hat_admin_crown",
			colors: []string{"#ff0", "#0ff", "#f00", "#f0f", "#0f0", "#00f"},
			height: 64,

			animationOverrides: map[string]string{
				"": "_ac",
			},
			adminOnly: true,
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
