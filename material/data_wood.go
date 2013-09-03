package material

type WoodType uint64

func (t WoodType) Name() string      { return woodTypes[t].name }
func (t WoodType) BarkColor() string { return woodTypes[t].barkColor }
func (t WoodType) LeafColor() string { return woodTypes[t].leafColor }
func (t WoodType) LeafType() uint    { return woodTypes[t].leafType }
func (t WoodType) Strength() uint64  { return woodTypes[t].strength }
func (t WoodType) Density() uint64   { return woodTypes[t].density }

var woodTypes = []struct {
	name      string
	barkColor string
	leafColor string
	leafType  uint
	strength  uint64
	density   uint64
}{
	{
		name:      "wood0",
		strength:  5 << 0,
		barkColor: "#000",
		density:   180,
	},
	{
		name:      "wood1",
		strength:  5 << 1,
		barkColor: "#111",
		density:   180,
	},
	{
		name:      "wood2",
		strength:  5 << 2,
		barkColor: "#222",
		density:   180,
	},
	{
		name:      "wood3",
		strength:  5 << 3,
		barkColor: "#333",
		density:   180,
	},
	{
		name:      "wood4",
		strength:  5 << 4,
		barkColor: "#444",
		density:   180,
	},
	{
		name:      "wood5",
		strength:  5 << 5,
		barkColor: "#555",
		density:   180,
	},
	{
		name:      "wood6",
		strength:  5 << 6,
		barkColor: "#666",
		density:   180,
	},
	{
		name:      "wood7",
		strength:  5 << 7,
		barkColor: "#777",
		density:   180,
	},
	{
		name:      "wood8",
		strength:  5 << 8,
		barkColor: "#888",
		density:   180,
	},
	{
		name:      "wood9",
		strength:  5 << 9,
		barkColor: "#999",
		density:   180,
	},
	{
		name:      "wood10",
		strength:  5 << 10,
		barkColor: "#aaa",
		density:   180,
	},
	{
		name:      "wood11",
		strength:  5 << 11,
		barkColor: "#bbb",
		density:   180,
	},
	{
		name:      "wood12",
		strength:  5 << 12,
		barkColor: "#ccc",
		density:   180,
	},
	{
		name:      "wood13",
		strength:  5 << 13,
		barkColor: "#ddd",
		density:   180,
	},
	{
		name:      "wood14",
		strength:  5 << 14,
		barkColor: "#eee",
		density:   180,
	},
	{
		name:      "wood15",
		strength:  5 << 15,
		barkColor: "#fff",
		density:   180,
	},
	{
		name:      "wood16",
		strength:  1, // TODO
		barkColor: "#d2b48c",
		leafColor: "#4b5a3f",
		leafType:  1,
		density:   180,
	},
	{
		name:      "wood17",
		strength:  1, // TODO
		barkColor: "#b5aa8b",
		leafColor: "#cf5123",
		leafType:  2,
		density:   180,
	},
}
