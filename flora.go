package main

type FloraType uint16

const (
	LeafPlant FloraType = iota
	FlowerPlant
	BulbPlant
	MagmaFlowerPlant

	floraTypeCount
)

var floraTypeInfo = [floraTypeCount]struct {
	Name        string
	LeafColor   Color
	BulbColor   Color
	StemColor   Color
	FlowerColor Color

	CompassPetalColor    Color
	BoringPetalColor     Color
	SuspiciousPetalColor Color
}{
	LeafPlant: {
		Name:      "leaf",
		LeafColor: "#7a0",
		StemColor: "#7a0",
	},
	FlowerPlant: {
		Name:                 "flower",
		LeafColor:            "#6f0",
		StemColor:            "#6f0",
		FlowerColor:          "#af6",
		SuspiciousPetalColor: "#0ec",
	},
	BulbPlant: {
		Name:      "bulb",
		LeafColor: "#0fc",
		BulbColor: "#f0f",
	},
	MagmaFlowerPlant: {
		Name:      "magma flower",
		LeafColor: "#311",
		StemColor: "#f00",
		BulbColor: "#522",
	},
}

type Flora struct {
	networkID
	Type FloraType
}

func (f *Flora) Name() string {
	return floraTypeInfo[f.Type].Name + " plant"
}

func (f *Flora) Examine() string {
	return "a " + floraTypeInfo[f.Type].Name + " plant."
}

func (f *Flora) Serialize() *NetworkedObject {
	info := floraTypeInfo[f.Type]
	return &NetworkedObject{
		Sprite: "plant",
		Colors: []Color{info.LeafColor, info.StemColor, info.BulbColor, info.BoringPetalColor, info.CompassPetalColor, info.SuspiciousPetalColor, info.FlowerColor},
	}
}

func (f *Flora) Blocking() bool {
	return false
}

func (f *Flora) InteractOptions() []string {
	return []string{"pick"}
}

func (f *Flora) Interact(x uint8, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // pick
	}
}

func (f *Flora) ZIndex() int {
	return 0
}
