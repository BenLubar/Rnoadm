package main

type cosmeticInfo struct {
	Name    string
	Examine string

	Base      string
	BaseColor Color

	Layer1      string
	Layer1Color Color

	Layer2      string
	Layer2Color Color

	Layer3      string
	Layer3Color Color

	Layer4      string
	Layer4Color Color

	AdminOnly bool
}

type HatType uint16

const (
	ChickenHat HatType = iota
	BearRidingAUnicycle
	UnusualStrangeVintageHoundDog

	hatTypeCount
)

var hatTypeInfo = [hatTypeCount]cosmeticInfo{
	ChickenHat: {
		Name:    "novelty foam chicken hat",
		Examine: "when you need to go in style, go in a novelty foam chicken's decapitated head.",

		Base:      "hat_chicken_base",
		BaseColor: "#fff",

		Layer1:      "hat_chicken_beak",
		Layer1Color: "#ff0",

		Layer2:      "hat_chicken_comb",
		Layer2Color: "#f00",

		AdminOnly: true,
	},
	BearRidingAUnicycle: {
		Name:    "bear riding a unicycle",
		Examine: "russia, amirite?",

		Base:      "hat_bear_riding_a_unicycle",
		BaseColor: "#fff",

		AdminOnly: true,
	},
	UnusualStrangeVintageHoundDog: {
		Name:    "unusual strange vintage hound dog",
		Examine: "your opponents will be all shook up when they see these sweet shades and coif. kills: 72",

		Base:      "hat_unusual_strange_vintage_hound_dog_shades",
		BaseColor: "#fff",

		Layer1:      "hat_unusual_strange_vintage_hound_dog_pomp",
		Layer1Color: "#fff",

		AdminOnly: true,
	},
	SpanishWarMask: {
		Name:	"spanish war mask",
		Examine: "many Pedros have worn this mask before you.",
		
		Base: "hat_spanish_war_mask_bottom",
		BaseColor: #00f,
		
		Layer1: "hat_spanish_war_mask_top",
		Base Color: #f00,
		
		AdminOnly: true,
	},
}

type ShirtType uint16

const (
	PlainWhiteTee ShirtType = iota

	shirtTypeCount
)

var shirtTypeInfo = [shirtTypeCount]cosmeticInfo{
	PlainWhiteTee: {
		Name:    "plain white tee",
		Examine: "$120. by fruit, feat. the loom.",

		Base:      "player_shirt",
		BaseColor: "#fff",
	},
}

type PantsType uint16

const (
	OffBrandJeans PantsType = iota

	pantsTypeCount
)

var pantsTypeInfo = [pantsTypeCount]cosmeticInfo{
	OffBrandJeans: {
		Name:    "off-brand jeans",
		Examine: "these have seen some use.",

		Base:      "player_pants",
		BaseColor: "#758a9d",
	},
}

type ShoeType uint16

const (
	WhiteSneakers ShoeType = iota

	shoeTypeCount
)

var shoeTypeInfo = [shoeTypeCount]cosmeticInfo{
	WhiteSneakers: {
		Name:    "white sneakers",
		Examine: "your favorite pair.",

		Base:      "player_shoes",
		BaseColor: "#eef8f0",
	},
}

type Hat struct {
	Type        HatType
	CustomColor [5]Color
}

type Shirt struct {
	Type        ShirtType
	CustomColor [5]Color
}

type Pants struct {
	Type        PantsType
	CustomColor [5]Color
}

type Shoes struct {
	Type        ShoeType
	CustomColor [5]Color
}

func (h *Hat) Name() string {
	return hatTypeInfo[h.Type].Name
}

func (s *Shirt) Name() string {
	return shirtTypeInfo[s.Type].Name
}

func (p *Pants) Name() string {
	return pantsTypeInfo[p.Type].Name
}

func (s *Shoes) Name() string {
	return shoeTypeInfo[s.Type].Name
}

func (h *Hat) Examine() string {
	return hatTypeInfo[h.Type].Examine
}

func (s *Shirt) Examine() string {
	return shirtTypeInfo[s.Type].Examine
}

func (p *Pants) Examine() string {
	return pantsTypeInfo[p.Type].Examine
}

func (s *Shoes) Examine() string {
	return shoeTypeInfo[s.Type].Examine
}

func (h *Hat) Blocking() bool {
	return false
}

func (s *Shirt) Blocking() bool {
	return false
}

func (p *Pants) Blocking() bool {
	return false
}

func (s *Shoes) Blocking() bool {
	return false
}

func paintCosmetic(x, y int, info cosmeticInfo, custom [5]Color, setcell func(int, int, string, string, Color)) {
	color := info.BaseColor
	if custom[0] != "" {
		color = custom[0]
	}
	setcell(x, y, "", info.Base, color)
	if info.Layer1 != "" {
		color = info.Layer1Color
		if custom[1] != "" {
			color = custom[1]
		}
		setcell(x, y, "", info.Layer1, color)
	}
	if info.Layer2 != "" {
		color = info.Layer2Color
		if custom[2] != "" {
			color = custom[2]
		}
		setcell(x, y, "", info.Layer2, color)
	}
	if info.Layer3 != "" {
		color = info.Layer3Color
		if custom[3] != "" {
			color = custom[3]
		}
		setcell(x, y, "", info.Layer3, color)
	}
	if info.Layer4 != "" {
		color = info.Layer4Color
		if custom[4] != "" {
			color = custom[4]
		}
		setcell(x, y, "", info.Layer4, color)
	}
}

func (h *Hat) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	info := hatTypeInfo[h.Type]
	custom := h.CustomColor
	paintCosmetic(x, y, info, custom, setcell)
}

func (s *Shirt) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	info := shirtTypeInfo[s.Type]
	custom := s.CustomColor
	paintCosmetic(x, y, info, custom, setcell)
}

func (p *Pants) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	info := pantsTypeInfo[p.Type]
	custom := p.CustomColor
	paintCosmetic(x, y, info, custom, setcell)
}

func (s *Shoes) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	info := shoeTypeInfo[s.Type]
	custom := s.CustomColor
	paintCosmetic(x, y, info, custom, setcell)
}

func (h *Hat) InteractOptions() []string {
	return nil
}

func (s *Shirt) InteractOptions() []string {
	return nil
}

func (p *Pants) InteractOptions() []string {
	return nil
}

func (s *Shoes) InteractOptions() []string {
	return nil
}

func (h *Hat) IsItem() {}

func (s *Shirt) IsItem() {}

func (p *Pants) IsItem() {}

func (s *Shoes) IsItem() {}

func (h *Hat) AdminOnly() bool {
	return hatTypeInfo[h.Type].AdminOnly
}

func (s *Shirt) AdminOnly() bool {
	return shirtTypeInfo[s.Type].AdminOnly
}

func (p *Pants) AdminOnly() bool {
	return pantsTypeInfo[p.Type].AdminOnly
}

func (s *Shoes) AdminOnly() bool {
	return shoeTypeInfo[s.Type].AdminOnly
}
