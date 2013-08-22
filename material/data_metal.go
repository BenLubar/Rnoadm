package material

type MetalType uint64

func (t MetalType) Name() string     { return metalTypes[t].name }
func (t MetalType) Color() string    { return metalTypes[t].color }
func (t MetalType) Strength() uint64 { return metalTypes[t].strength }

var metalTypes = []struct {
	name     string
	color    string
	strength uint64
}{
	{
		name:     "metal0",
		strength: 5 << 0,
		color:    "#000",
	},
	{
		name:     "metal1",
		strength: 5 << 1,
		color:    "#111",
	},
	{
		name:     "metal2",
		strength: 5 << 2,
		color:    "#222",
	},
	{
		name:     "metal3",
		strength: 5 << 3,
		color:    "#333",
	},
	{
		name:     "metal4",
		strength: 5 << 4,
		color:    "#444",
	},
	{
		name:     "metal5",
		strength: 5 << 5,
		color:    "#555",
	},
	{
		name:     "metal6",
		strength: 5 << 6,
		color:    "#666",
	},
	{
		name:     "metal7",
		strength: 5 << 7,
		color:    "#777",
	},
	{
		name:     "metal8",
		strength: 5 << 8,
		color:    "#888",
	},
	{
		name:     "metal9",
		strength: 5 << 9,
		color:    "#999",
	},
	{
		name:     "metal10",
		strength: 5 << 10,
		color:    "#aaa",
	},
	{
		name:     "metal11",
		strength: 5 << 11,
		color:    "#bbb",
	},
	{
		name:     "metal12",
		strength: 5 << 12,
		color:    "#ccc",
	},
	{
		name:     "metal13",
		strength: 5 << 13,
		color:    "#ddd",
	},
	{
		name:     "metal14",
		strength: 5 << 14,
		color:    "#eee",
	},
	{
		name:     "metal15",
		strength: 5 << 15,
		color:    "#fff",
	},
}
