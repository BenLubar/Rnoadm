package main

type FloraType uint16

const (
	LeafPlant FloraType = iota
	FlowerPlant
	BulbPlant

	floraTypeCount
)

var floraTypeInfo = [floraTypeCount]struct {
	Name        string
	LeafColor   Color
	BulbColor   Color
	StemColor   Color
	FlowerColor Color
	PetalColor  Color
}{
	LeafPlant: {
		Name:      "leaf",
		LeafColor: "#4f0",
	},
	FlowerPlant: {
		Name:        "flower",
		LeafColor:   "#6f0",
		StemColor:   "#6f0",
		FlowerColor: "#af6",
		PetalColor:  "#0ec",
	},
	BulbPlant: {
		Name:      "bulb",
		LeafColor: "#047",
		BulbColor: "#f0f",
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
		setcell(x, y, "", "plant_base_l0", color)
	}
	if color := floraTypeInfo[f.Type].BulbColor; color != "" {
		setcell(x, y, "", "plant_bulb_l1", color)
	}
	if color := floraTypeInfo[f.Type].StemColor; color != "" {
		setcell(x, y, "", "plant_flowerstem_l1", color)
	}
	if color := floraTypeInfo[f.Type].PetalColor; color != "" {
		setcell(x, y, "", "plant_flowerpetals_l2", color)
	}
	if color := floraTypeInfo[f.Type].FlowerColor; color != "" {
		setcell(x, y, "", "plant_flowercenter_l2", color)
	}
}

func (f *Flora) Blocking() bool {
	return false
}

func (f *Flora) InteractOptions() []string {
	return []string{"pick"}
}
