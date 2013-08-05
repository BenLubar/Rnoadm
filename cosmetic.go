package main

type cosmeticInfo struct {
	Name    string
	Examine string

	Height uint8

	Base        string
	BaseColor   Color
	Layer1Color Color
	Layer2Color Color
	Layer3Color Color
	Layer4Color Color

	AdminOnly bool
}

type HatType uint16

const (
	ChickenHat HatType = iota
	BearRidingAUnicycle
	UnusualStrangeVintageHoundDog
	SpanishWarMask
	BigWig

	hatTypeCount
)

var hatTypeInfo = [hatTypeCount]cosmeticInfo{
	ChickenHat: {
		Name:    "novelty foam chicken hat",
		Examine: "when you need to go in style, go in a novelty foam chicken's decapitated head.",

		Base:        "hat_chicken",
		BaseColor:   "#fff",
		Layer1Color: "#ff0",
		Layer2Color: "#f00",

		AdminOnly: true,
	},
	BearRidingAUnicycle: {
		Name:    "bear riding a unicycle",
		Examine: "russia, amirite?",

		Height: 77,

		Base:      "hat_bear_riding_a_unicycle",
		BaseColor: "#fff",

		AdminOnly: true,
	},
	UnusualStrangeVintageHoundDog: {
		Name:    "unusual strange vintage hound dog",
		Examine: "your opponents will be all shook up when they see these sweet shades and coif. kills: 72",

		Base:        "hat_unusual_strange_vintage_hound_dog",
		BaseColor:   "#fff",
		Layer1Color: "#fff",

		AdminOnly: true,
	},
	SpanishWarMask: {
		Name:    "spanish war mask",
		Examine: "many Pedros have worn this mask before you.",

		Base:        "hat_spanish_war_mask",
		BaseColor:   "#900",
		Layer1Color: "#007",

		AdminOnly: true,
	},
	BigWig: {
		Name:    "big wig",
		Examine: "congress?",

		Base:      "hat_big_wig",
		BaseColor: "#630",
	},
}

type ShirtType uint16

const (
	HipHopTeeShirt ShirtType = iota

	shirtTypeCount
)

var shirtTypeInfo = [shirtTypeCount]cosmeticInfo{
	HipHopTeeShirt: {
		Name:    "hip hop tee shirt",
		Examine: "$120. by fruit, feat. the loom.",

		Base:      "shirt_basic",
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

		Base:      "pants_basic",
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

		Base:      "shoes_basic",
		BaseColor: "#eef8f0",
	},
}

type Hat struct {
	networkID
	Type        HatType
	CustomColor [5]Color
}

type Shirt struct {
	networkID
	Type        ShirtType
	CustomColor [5]Color
}

type Pants struct {
	networkID
	Type        PantsType
	CustomColor [5]Color
}

type Shoes struct {
	networkID
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

func serializeCosmetic(info cosmeticInfo, custom [5]Color) *NetworkedObject {
	colors := []Color{info.BaseColor, info.Layer1Color, info.Layer2Color, info.Layer3Color, info.Layer4Color}

	for i := len(colors) - 1; i > 0; i-- {
		if colors[i] == "" {
			colors = colors[:i]
		}
	}

	for i := range colors {
		if custom[i] != "" {
			colors[i] = custom[i]
		}
	}

	return &NetworkedObject{
		Sprite: info.Base,
		Colors: colors,
	}
}

func (h *Hat) Serialize() *NetworkedObject {
	return serializeCosmetic(hatTypeInfo[h.Type], h.CustomColor)
}
func (s *Shirt) Serialize() *NetworkedObject {
	return serializeCosmetic(shirtTypeInfo[s.Type], s.CustomColor)
}
func (p *Pants) Serialize() *NetworkedObject {
	return serializeCosmetic(pantsTypeInfo[p.Type], p.CustomColor)
}
func (s *Shoes) Serialize() *NetworkedObject {
	return serializeCosmetic(shoeTypeInfo[s.Type], s.CustomColor)
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

/*func paintCosmetic(x, y int, info cosmeticInfo, custom [5]Color, setcell func(int, int, PaintCell), worn bool, frame uint8, offsetX, offsetY int8) {
	var height uint8
	if worn {
		height = info.Height
	}
	color := info.BaseColor
	if custom[0] != "" {
		color = custom[0]
	}
	setcell(x, y, PaintCell{
		Sprite: info.Base,
		Color:  color,
		Height: height,
		SheetX: frame,
		X:      offsetX,
		Y:      offsetY,
		ZIndex: 501,
	})
	if info.Layer1Color != "" {
		color = info.Layer1Color
		if custom[1] != "" {
			color = custom[1]
		}
		setcell(x, y, PaintCell{
			Sprite: info.Base,
			Color:  color,
			Height: height,
			SheetX: frame,
			SheetY: 1,
			X:      offsetX,
			Y:      offsetY,
			ZIndex: 502,
		})
	}
	if info.Layer2Color != "" {
		color = info.Layer2Color
		if custom[2] != "" {
			color = custom[2]
		}
		setcell(x, y, PaintCell{
			Sprite: info.Base,
			Color:  color,
			Height: height,
			SheetX: frame,
			SheetY: 2,
			X:      offsetX,
			Y:      offsetY,
			ZIndex: 503,
		})
	}
	if info.Layer3Color != "" {
		color = info.Layer3Color
		if custom[3] != "" {
			color = custom[3]
		}
		setcell(x, y, PaintCell{
			Sprite: info.Base,
			Color:  color,
			Height: height,
			SheetX: frame,
			SheetY: 3,
			X:      offsetX,
			Y:      offsetY,
			ZIndex: 504,
		})
	}
	if info.Layer4Color != "" {
		color = info.Layer4Color
		if custom[4] != "" {
			color = custom[4]
		}
		setcell(x, y, PaintCell{
			Sprite: info.Base,
			Color:  color,
			Height: height,
			SheetX: frame,
			SheetY: 4,
			X:      offsetX,
			Y:      offsetY,
			ZIndex: 505,
		})
	}
}

func (h *Hat) Paint(x, y int, setcell func(int, int, PaintCell)) {
	info := hatTypeInfo[h.Type]
	custom := h.CustomColor
	paintCosmetic(x, y, info, custom, setcell, false, 0, 0, 0)
}

func (s *Shirt) Paint(x, y int, setcell func(int, int, PaintCell)) {
	info := shirtTypeInfo[s.Type]
	custom := s.CustomColor
	paintCosmetic(x, y, info, custom, setcell, false, 0, 0, 0)
}

func (p *Pants) Paint(x, y int, setcell func(int, int, PaintCell)) {
	info := pantsTypeInfo[p.Type]
	custom := p.CustomColor
	paintCosmetic(x, y, info, custom, setcell, false, 0, 0, 0)
}

func (s *Shoes) Paint(x, y int, setcell func(int, int, PaintCell)) {
	info := shoeTypeInfo[s.Type]
	custom := s.CustomColor
	paintCosmetic(x, y, info, custom, setcell, false, 0, 0, 0)
}

func (h *Hat) PaintWorn(x, y int, setcell func(int, int, PaintCell), frame uint8, offsetX, offsetY int8) {
	info := hatTypeInfo[h.Type]
	custom := h.CustomColor
	paintCosmetic(x, y, info, custom, setcell, true, frame, offsetX, offsetY)
}

func (s *Shirt) PaintWorn(x, y int, setcell func(int, int, PaintCell), frame uint8, offsetX, offsetY int8) {
	info := shirtTypeInfo[s.Type]
	custom := s.CustomColor
	paintCosmetic(x, y, info, custom, setcell, true, frame, offsetX, offsetY)
}

func (p *Pants) PaintWorn(x, y int, setcell func(int, int, PaintCell), frame uint8, offsetX, offsetY int8) {
	info := pantsTypeInfo[p.Type]
	custom := p.CustomColor
	paintCosmetic(x, y, info, custom, setcell, true, frame, offsetX, offsetY)
}

func (s *Shoes) PaintWorn(x, y int, setcell func(int, int, PaintCell), frame uint8, offsetX, offsetY int8) {
	info := shoeTypeInfo[s.Type]
	custom := s.CustomColor
	paintCosmetic(x, y, info, custom, setcell, true, frame, offsetX, offsetY)
}*/

func (h *Hat) InteractOptions() []string {
	return []string{"wear"}
}

func (s *Shirt) InteractOptions() []string {
	return []string{"wear"}
}

func (p *Pants) InteractOptions() []string {
	return []string{"wear"}
}

func (s *Shoes) InteractOptions() []string {
	return []string{"wear"}
}

func (h *Hat) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // wear
		player.Lock()
		player.Equip(h, true)
		player.Unlock()
	}
}

func (s *Shirt) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // wear
		player.Lock()
		player.Equip(s, true)
		player.Unlock()
	}
}

func (p *Pants) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // wear
		player.Lock()
		player.Equip(p, true)
		player.Unlock()
	}
}

func (s *Shoes) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // wear
		player.Lock()
		player.Equip(s, true)
		player.Unlock()
	}
}

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

func (h *Hat) ZIndex() int {
	return 25
}

func (s *Shirt) ZIndex() int {
	return 25
}

func (p *Pants) ZIndex() int {
	return 25
}

func (s *Shoes) ZIndex() int {
	return 25
}
