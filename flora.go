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
	Type FloraType
}

func (f *Flora) Name() string {
	return floraTypeInfo[f.Type].Name + " plant"
}

func (f *Flora) Examine() string {
	return "a " + floraTypeInfo[f.Type].Name + " plant."
}

func (f *Flora) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	if color := floraTypeInfo[f.Type].LeafColor; color != "" {
		setcell(x, y, "", "item_plant_leaves", color)
	}
	if color := floraTypeInfo[f.Type].StemColor; color != "" {
		setcell(x, y, "", "item_plant_stem", color)
	}
	if color := floraTypeInfo[f.Type].BulbColor; color != "" {
		setcell(x, y, "", "item_plant_bulb", color)
	}
	if color := floraTypeInfo[f.Type].BoringPetalColor; color != "" {
		setcell(x, y, "", "item_plant_flower_boring", color)
	}
	if color := floraTypeInfo[f.Type].CompassPetalColor; color != "" {
		setcell(x, y, "", "item_plant_flower_compass", color)
	}
	if color := floraTypeInfo[f.Type].SuspiciousPetalColor; color != "" {
		setcell(x, y, "", "item_plant_flower_suspicious", color)
	}
	if color := floraTypeInfo[f.Type].FlowerColor; color != "" {
		setcell(x, y, "", "item_plant_flower_center", color)
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
