package material

type WoodType uint64

func (t WoodType) Name() string      { return woodTypes[t].name }
func (t WoodType) BarkColor() string { return woodTypes[t].barkColor }
func (t WoodType) LeafColor() string { return woodTypes[t].leafColor }
func (t WoodType) LeafType() uint    { return woodTypes[t].leafType }
func (t WoodType) Strength() uint64  { return woodTypes[t].strength }

var woodTypes = []struct {
	name      string
	barkColor string
	leafColor string
	leafType  uint
	strength  uint64
}{
	{
		name:      "wood0",
		strength:  5 << 0,
		barkColor: "#000",
	},
	{
		name:      "wood1",
		strength:  5 << 1,
		barkColor: "#111",
	},
	{
		name:      "wood2",
		strength:  5 << 2,
		barkColor: "#222",
	},
	{
		name:      "wood3",
		strength:  5 << 3,
		barkColor: "#333",
	},
	{
		name:      "wood4",
		strength:  5 << 4,
		barkColor: "#444",
	},
	{
		name:      "wood5",
		strength:  5 << 5,
		barkColor: "#555",
	},
	{
		name:      "wood6",
		strength:  5 << 6,
		barkColor: "#666",
	},
	{
		name:      "wood7",
		strength:  5 << 7,
		barkColor: "#777",
	},
	{
		name:      "wood8",
		strength:  5 << 8,
		barkColor: "#888",
	},
	{
		name:      "wood9",
		strength:  5 << 9,
		barkColor: "#999",
	},
	{
		name:      "wood10",
		strength:  5 << 10,
		barkColor: "#aaa",
	},
	{
		name:      "wood11",
		strength:  5 << 11,
		barkColor: "#bbb",
	},
	{
		name:      "wood12",
		strength:  5 << 12,
		barkColor: "#ccc",
	},
	{
		name:      "wood13",
		strength:  5 << 13,
		barkColor: "#ddd",
	},
	{
		name:      "wood14",
		strength:  5 << 14,
		barkColor: "#eee",
	},
	{
		name:      "wood15",
		strength:  5 << 15,
		barkColor: "#fff",
	},
	{
		name:      "wood16",
		strength:  1, // TODO
		barkColor: "#d2b48c",
		leafColor: "#4b5a3f",
		leafType:  1,
	},
	{
		name:      "wood17",
		strength:  1, // TODO
		barkColor: "#b5aa8b",
		leafColor: "#cf5123",
		leafType:  2,
	},
}
