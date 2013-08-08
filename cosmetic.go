package main

type cosmeticInfo struct {
	Article string
	Name    string
	Examine string

	Height uint16
	Sprite string
	Colors []Color

	HealthBonus uint64

	AdminOnly bool
}

type CosmeticType uint16

const (
	Headwear CosmeticType = iota
	Shirt
	Pants
	Shoes
	Breastplate
	Pauldrons
	Vambraces
	Gauntlets
	Tassets
	Greaves

	cosmeticTypeCount
)

var CosmeticSlotName [cosmeticTypeCount]string = [...]string{
	"head",
	"top",
	"legs",
	"feet",
	"chest",
	"shoulders",
	"wrists",
	"hands",
	"waist",
	"ankles",
}

var CosmeticSlotOrder [cosmeticTypeCount]CosmeticType = [...]CosmeticType{
	Pants,
	Shirt,
	Greaves,
	Vambraces,
	Breastplate,
	Tassets,
	Pauldrons,
	Shoes,
	Gauntlets,
	Headwear,
}

const metalColor = "metal"

var cosmetics = [cosmeticTypeCount][]cosmeticInfo{
	Headwear: {
		{
			Article: "the ",
			Name:    "nonexistent hat",
			Examine: "this hat doesn't actually exist.",
		},
		{
			Article: "a ",
			Name:    "novelty foam chicken hat",
			Examine: "when you need to go in style, go in a novelty foam chicken's decapitated head.",

			Sprite: "hat_novelty_foam_chicken",
			Colors: []Color{"#fff", "#ff0", "#f00"},

			AdminOnly: true,
		},
		{
			Article: "a ",
			Name:    "bear riding a unicycle",
			Examine: "russia, amirite?",

			Height: 87,

			Sprite: "hat_bear_riding_a_unicycle",
			Colors: []Color{"#ccf", "#ecb", "#fff"},

			AdminOnly: true,
		},
		{
			Article: "an ",
			Name:    "unusual strange vintage hound dog",
			Examine: "your opponents will be all shook up when they see these sweet shades and coif. kills: 72",

			Sprite: "hat_unusual_strange_vintage_hound_dog",
			Colors: []Color{"#fff", "#fff"},

			AdminOnly: true,
		},
		{
			Article: "a ",
			Name:    "spanish war mask",
			Examine: "many Pedros have worn this mask before you.",

			Sprite: "hat_spanish_war_mask",
			Colors: []Color{"#900", "#007"},

			AdminOnly: true,
		},
		{
			Article: "a ",
			Name:    "helmet",
			Examine: "this will protect me from space aliens!",

			Sprite: "hat_helmet",
			Colors: []Color{metalColor},

			HealthBonus: 5,
		},
	},
	Shirt: {
		{
			Article: "a ",
			Name:    "hip hop tee shirt",
			Examine: "$120. by fruit, feat. the loom.",

			Sprite: "shirt_basic",
			Colors: []Color{"#fff"},
		},
	},
	Pants: {
		{
			Article: "a pair of ",
			Name:    "off-brand jeans",
			Examine: "these have seen some use.",

			Sprite: "pants_basic",
			Colors: []Color{"#758a9d"},
		},
	},
	Shoes: {
		{
			Article: "a pair of ",
			Name:    "sneakers",
			Examine: "your favorite pair.",

			Sprite: "shoes_basic",
			Colors: []Color{"#eef8f0"},
		},
		{
			Article: "a pair of ",
			Name:    "boots",
			Examine: "forget steel toes. these are metal everywhere.",

			Sprite: "shoes_boots",
			Colors: []Color{metalColor},

			HealthBonus: 1,
		},
	},
	Breastplate: {
		{
			Article: "a ",
			Name:    "breastplate",
			Examine: "this should protect me from backstabbers. and frontstabbers.",

			Sprite: "breastplate_basic",
			Colors: []Color{metalColor},

			HealthBonus: 10,
		},
	},
	Pauldrons: {
		{
			Article: "a pair of ",
			Name:    "pauldrons",
			Examine: "posh.",

			Sprite: "pauldrons_basic",
			Colors: []Color{metalColor},

			HealthBonus: 1,
		},
	},
	Vambraces: {
		{
			Article: "a pair of ",
			Name:    "vambraces",
			Examine: "metal wristbands.",

			Sprite: "vambraces_basic",
			Colors: []Color{metalColor},

			HealthBonus: 2,
		},
	},
	Gauntlets: {
		{
			Article: "a pair of ",
			Name:    "gauntlets",
			Examine: "metal gloves? sounds like a recipe for disaster.",

			Sprite: "gauntlets_basic",
			Colors: []Color{metalColor},

			HealthBonus: 1,
		},
	},
	Tassets: {
		{
			Article: "",
			Name:    "tassets",
			Examine: "didn't they ban those?",

			Sprite: "tassets_basic",
			Colors: []Color{metalColor},

			HealthBonus: 5,
		},
	},
	Greaves: {
		{
			Article: "a pair of ",
			Name:    "greaves",
			Examine: "there was a joke here, but then it took an arrow to the knee.",

			Sprite: "greaves_basic",
			Colors: []Color{metalColor},

			HealthBonus: 2,
		},
	},
}

type Cosmetic struct {
	networkID
	Type   CosmeticType
	ID     uint64
	Custom []Color
	Metal  MetalType
}

func (c *Cosmetic) Article() string {
	return cosmetics[c.Type][c.ID].Article
}

func (c *Cosmetic) Name() string {
	name := cosmetics[c.Type][c.ID].Name
	if c.Metal != 0 {
		name = metalTypeInfo[c.Metal].Name + " " + name
	}
	return name
}

func (c *Cosmetic) Examine() string {
	return cosmetics[c.Type][c.ID].Examine
}

func (c *Cosmetic) Serialize() *NetworkedObject {
	info := cosmetics[c.Type][c.ID]
	obj := &NetworkedObject{
		Name:   info.Name,
		Sprite: info.Sprite,
		Height: info.Height,
		Colors: make([]Color, len(info.Colors)),
	}
	copy(obj.Colors, info.Colors)
	for i := range obj.Colors {
		if obj.Colors[i] == metalColor {
			obj.Colors[i] = metalTypeInfo[c.Metal].Color
			if obj.Colors[i] == "" {
				obj.Colors[i] = "#f0f"
			}
		}
	}
	for i := range obj.Colors {
		if len(c.Custom) <= i {
			break
		}
		if c.Custom[i] != "" {
			obj.Colors[i] = c.Custom[i]
		}
	}
	return obj
}

func (c *Cosmetic) Blocking() bool {
	return false
}

func (c *Cosmetic) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // wear
		player.Lock()
		player.Equip(c, true)
		player.Unlock()
	}
}

func (c *Cosmetic) AdminOnly() bool {
	return cosmetics[c.Type][c.ID].AdminOnly || metalTypeInfo[c.Metal].Strength >= 1<<60
}

func (c *Cosmetic) ZIndex() int {
	return 25
}

func (c *Cosmetic) Exists() bool {
	return (c.Type != 0 || c.ID != 0) && c.Type < cosmeticTypeCount && c.ID < uint64(len(cosmetics[c.Type]))
}

func (c *Cosmetic) HealthBonus() uint64 {
	base := cosmetics[c.Type][c.ID].HealthBonus
	bonus := base * 100
	if c.Metal != 0 {
		if metalTypeInfo[c.Metal].Strength > 100000 {
			bonus += metalTypeInfo[c.Metal].Strength / 1000 * base
		} else {
			bonus += metalTypeInfo[c.Metal].Strength * base / 1000
		}
	}
	return bonus
}
