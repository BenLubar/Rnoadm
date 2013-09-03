package material

type MetalType uint64

func (t MetalType) Name() string     { return metalTypes[t].name }
func (t MetalType) Color() string    { return metalTypes[t].color }
func (t MetalType) OreColor() string { return metalTypes[t].color }
func (t MetalType) Strength() uint64 { return metalTypes[t].strength }
func (t MetalType) Density() uint64  { return metalTypes[t].density }

var metalTypes = []struct {
	name     string
	color    string
	strength uint64
	density  uint64
}{
	{
		name:     "metal0",
		strength: 5 << 0,
		color:    "#000",
		density:  300,
	},
	{
		name:     "metal1",
		strength: 5 << 1,
		color:    "#111",
		density:  310,
	},
	{
		name:     "metal2",
		strength: 5 << 2,
		color:    "#222",
		density:  320,
	},
	{
		name:     "metal3",
		strength: 5 << 3,
		color:    "#333",
		density:  330,
	},
	{
		name:     "metal4",
		strength: 5 << 4,
		color:    "#444",
		density:  340,
	},
	{
		name:     "metal5",
		strength: 5 << 5,
		color:    "#555",
		density:  350,
	},
	{
		name:     "metal6",
		strength: 5 << 6,
		color:    "#666",
		density:  360,
	},
	{
		name:     "metal7",
		strength: 5 << 7,
		color:    "#777",
		density:  370,
	},
	{
		name:     "metal8",
		strength: 5 << 8,
		color:    "#888",
		density:  380,
	},
	{
		name:     "metal9",
		strength: 5 << 9,
		color:    "#999",
		density:  390,
	},
	{
		name:     "metal10",
		strength: 5 << 10,
		color:    "#aaa",
		density:  400,
	},
	{
		name:     "metal11",
		strength: 5 << 11,
		color:    "#bbb",
		density:  410,
	},
	{
		name:     "metal12",
		strength: 5 << 12,
		color:    "#ccc",
		density:  420,
	},
	{
		name:     "metal13",
		strength: 5 << 13,
		color:    "#ddd",
		density:  430,
	},
	{
		name:     "metal14",
		strength: 5 << 14,
		color:    "#eee",
		density:  440,
	},
	{
		name:     "metal15",
		strength: 5 << 15,
		color:    "#fff",
		density:  450,
	},
}
