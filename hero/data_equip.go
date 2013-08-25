package hero

const metalColor = "__metal__"
const stoneColor = "__stone__"
const woodColor = "__wood__"

var equippables = [equipSlotCount][]struct {
	name    string
	examine string
	sprite  string
	colors  []string
	height  uint // default 32
	width   uint // default 32

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
			name:    "crown of the Origin",
			examine: "legend says this crown can only be touched by the Founders.",
			sprite:  "hat_admin_crown",
			colors:  []string{"#ff0", "#0ff", "#888", "#888", "#888", "#888"},
			height:  64,

			animationOverrides: map[string]string{
				"":   "_ac",
				"wa": "wa_ac",
			},
			adminOnly: true,
		},
		{
			name:    "helmet",
			examine: "it's like a metal sock for your head!",
			sprite:  "hat_helmet",
			colors:  []string{metalColor},
			height:  48,
		},
		{
			name:   "Steve's suave swoosh",
			sprite: "hat_steves_suave_swoosh",
			colors: []string{"#c8421e"},
			height: 64,

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
			name:   "plain sneakers",
			sprite: "shoes_basic",
			colors: []string{"#cec"},
			height: 48,
		},
		{
			name:   "boots",
			sprite: "shoes_boots",
			colors: []string{metalColor},
			height: 48,
		},
	},
	SlotShoulders: {
		{
			name:   "pauldrons",
			sprite: "pauldrons_basic",
			colors: []string{metalColor},
			height: 48,
		},
	},
	SlotChest: {
		{
			name:   "chain mail",
			sprite: "chest_chainmail",
			colors: []string{metalColor},
			height: 48,
		},
	},
	SlotArms: {
		{
			name:   "vambraces",
			sprite: "vambraces_basic",
			colors: []string{metalColor},
			height: 48,
		},
	},
	SlotHands: {
		{
			name:   "gauntlets",
			sprite: "gauntlets_basic",
			colors: []string{metalColor},
			height: 48,
		},
	},
	SlotWaist: {
		{
			name:   "tassets",
			sprite: "tassets_basic",
			colors: []string{metalColor},
			height: 48,
		},
	},
	SlotLegs: {
		{
			name:   "greaves",
			sprite: "greaves_basic",
			colors: []string{metalColor},
			height: 48,
		},
	},
	SlotMainHand: {},
	SlotOffHand:  {},
	SlotPickaxe: {
		{
			name:   "pickaxe",
			sprite: "item_tools",
			colors: []string{woodColor, metalColor},
			height: 32,
		},
	},
	SlotHatchet: {
		{
			name:   "hatchet",
			sprite: "item_tools",
			colors: []string{woodColor, "", metalColor},
			height: 32,
		},
	},
}
