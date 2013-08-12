package main

type cosmeticInfo struct {
	Article string
	Name    string
	Examine string

	Height uint16
	Sprite string
	Colors []Color

	HealthBonus   uint64
	ArmorBonus    uint64
	DamageBonus   uint64
	AccuracyBonus uint64
	AttackSpeed   uint8

	Volume      uint64
	Density     uint64 // centigrams per cubic centimeter (cg/cc)
	MetalVolume uint64

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
	Weapon

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
	"right hand",
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
	Weapon,
}

const metalColor = "metal"

var cosmetics = [cosmeticTypeCount][]cosmeticInfo{
	Headwear: {
		{
			Article: "the none and only ",
			Name:    "nonexistent hat",
			Examine: "this hat doesn't actually exist. it is a placeholder item.",
		},
		{
			Article: "a ",
			Name:    "novelty foam chicken hat",
			Examine: "when you need to go in style, go in a novelty foam chicken's decapitated head.",

			Sprite: "hat_novelty_foam_chicken",
			Colors: []Color{"#fff", "#ff0", "#f00"},

			Volume:  739,
			Density: 19,

			AdminOnly: true,
		},
		{
			Article: "a ",
			Name:    "bear riding a unicycle",
			Examine: "russia, amirite?",

			Height: 87,

			Sprite: "hat_bear_riding_a_unicycle",
			Colors: []Color{"#ccf", "#ecb", "#fff"},

			Volume:  2000000,
			Density: 15,

			AdminOnly: true,
		},
		{
			Article: "an ",
			Name:    "unusual strange vintage hound dog",
			Examine: "your opponents will be all shook up when they see these sweet shades and coif. kills: 72",

			Sprite: "hat_unusual_strange_vintage_hound_dog",
			Colors: []Color{"#fff", "#fff"},

			Volume:  30,
			Density: 15,

			AdminOnly: true,
		},
		{
			Article: "a ",
			Name:    "spanish war mask",
			Examine: "many Pedros have worn this mask before you.",

			Sprite: "hat_spanish_war_mask",
			Colors: []Color{"#900", "#007"},

			Volume:  5,
			Density: 10,

			AdminOnly: true,
		},
		{
			Article: "a ",
			Name:    "helmet",
			Examine: "55% guaranteed to protect your brain from space aliens.",

			Sprite: "hat_helmet",
			Colors: []Color{metalColor},

			HealthBonus: 500,
			ArmorBonus:  500,

			MetalVolume: 20,
		},
	},
	Shirt: {
		{
			Article: "a ",
			Name:    "plain shirt",
			Examine: "the ultimate in style: a $6 monochromatic t-shirt.",

			Sprite: "shirt_basic",
			Colors: []Color{"#fff"},

			Volume:  10,
			Density: 17,
		},
	},
	Pants: {
		{
			Article: "a pair of ",
			Name:    "plain jeans",
			Examine: "worn out at the knees, but they still fit.",

			Sprite: "pants_basic",
			Colors: []Color{"#758a9d"},

			Volume:  10,
			Density: 21,
		},
	},
	Shoes: {
		{
			Article: "a pair of ",
			Name:    "sneakers",
			Examine: "plain old generic sneakers. not even a brand name.",

			Sprite: "shoes_basic",
			Colors: []Color{"#eef8f0"},

			Volume:  25,
			Density: 70,
		},
		{
			Article: "a pair of ",
			Name:    "boots",
			Examine: "forget steel toes. these are metal everywhere.",

			Sprite: "shoes_boots",
			Colors: []Color{metalColor},

			HealthBonus: 100,
			ArmorBonus:  100,

			MetalVolume: 20,
		},
	},
	Breastplate: {
		{
			Article: "a ",
			Name:    "breastplate",
			Examine: "this should protect me from backstabbers. and frontstabbers.",

			Sprite: "breastplate_basic",
			Colors: []Color{metalColor},

			HealthBonus: 1000,
			ArmorBonus:  1000,

			MetalVolume: 20,
		},
	},
	Pauldrons: {
		{
			Article: "a pair of ",
			Name:    "pauldrons",
			Examine: "protects you from shoulder aliens.",

			Sprite: "pauldrons_basic",
			Colors: []Color{metalColor},

			HealthBonus: 100,
			ArmorBonus:  100,

			MetalVolume: 20,
		},
	},
	Vambraces: {
		{
			Article: "a pair of ",
			Name:    "vambraces",
			Examine: "metal wristbands that make you look like you're from a 6th century motorcycle gang.",

			Sprite: "vambraces_basic",
			Colors: []Color{metalColor},

			HealthBonus: 200,
			ArmorBonus:  200,

			MetalVolume: 20,
		},
	},
	Gauntlets: {
		{
			Article: "a pair of ",
			Name:    "gauntlets",
			Examine: "metal gloves? sounds like a recipe for disaster.",

			Sprite: "gauntlets_basic",
			Colors: []Color{metalColor},

			HealthBonus: 100,
			ArmorBonus:  100,

			MetalVolume: 20,
		},
	},
	Tassets: {
		{
			Article: "",
			Name:    "tassets",
			Examine: "didn't they ban those?",

			Sprite: "tassets_basic",
			Colors: []Color{metalColor},

			HealthBonus: 500,
			ArmorBonus:  500,

			MetalVolume: 20,
		},
	},
	Greaves: {
		{
			Article: "a pair of ",
			Name:    "greaves",
			Examine: "there was a joke here, but then it took an arrow to the knee.",

			Sprite: "greaves_basic",
			Colors: []Color{metalColor},

			HealthBonus: 200,
			ArmorBonus:  200,

			MetalVolume: 20,
		},
	},
	Weapon: {
		{
			Article: "a ",
			Name:    "sword",
			Examine: "the sword that once belonged to the person who made it.",

			Sprite: "weapon_sword",
			Colors: []Color{metalColor, "#ec6"},

			DamageBonus:   1000,
			AccuracyBonus: 1000,
			AttackSpeed:   4,

			MetalVolume: 20,
		},
	},
}

type Cosmetic struct {
	networkID
	Type   CosmeticType
	ID     uint64
	Custom []Color

	Metal MetalType

	Quality uint64
}

func (c *Cosmetic) Article() string {
	article := cosmetics[c.Type][c.ID].Article
	if c.Metal != 0 && (article == "a " || article == "an ") {
		return metalTypeInfo[c.Metal].Article
	}
	return article
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
		Name:    c.Name(),
		Sprite:  info.Sprite,
		Height:  info.Height,
		Colors:  make([]Color, len(info.Colors)),
		Options: []string{"wear"},
		Item:    true,
	}
	copy(obj.Colors, info.Colors)
	for i := range obj.Colors {
		if obj.Colors[i] == metalColor {
			obj.Colors[i] = metalTypeInfo[c.Metal].Color
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
		player.Equip(c, true)
	}
}

func (c *Cosmetic) Volume() uint64 {
	return cosmetics[c.Type][c.ID].Volume + cosmetics[c.Type][c.ID].MetalVolume
}

func (c *Cosmetic) Weight() uint64 {
	return (cosmetics[c.Type][c.ID].Volume*cosmetics[c.Type][c.ID].Density + cosmetics[c.Type][c.ID].MetalVolume*metalTypeInfo[c.Metal].Density) / 100
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
	bonus := base + c.Quality
	bonus += metalTypeInfo[c.Metal].sqrtStr * base / 10
	return bonus
}

func (c *Cosmetic) ArmorBonus() uint64 {
	base := cosmetics[c.Type][c.ID].ArmorBonus
	bonus := base + c.Quality
	bonus += metalTypeInfo[c.Metal].sqrtStr * base / 10
	return bonus
}

func (c *Cosmetic) DamageBonus() uint64 {
	base := cosmetics[c.Type][c.ID].DamageBonus
	bonus := base + c.Quality
	bonus += metalTypeInfo[c.Metal].sqrtStr * base / 10
	return bonus
}

func (c *Cosmetic) AccuracyBonus() uint64 {
	base := cosmetics[c.Type][c.ID].AccuracyBonus
	bonus := base + c.Quality
	bonus += metalTypeInfo[c.Metal].sqrtStr * base / 10
	return bonus
}

func (c *Cosmetic) AttackSpeed() uint8 {
	return cosmetics[c.Type][c.ID].AttackSpeed
}
