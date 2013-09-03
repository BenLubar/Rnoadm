package material

type StoneType uint64

func (t StoneType) Name() string     { return stoneTypes[t].name }
func (t StoneType) Color() string    { return stoneTypes[t].color }
func (t StoneType) Strength() uint64 { return stoneTypes[t].strength }
func (t StoneType) Density() uint64  { return stoneTypes[t].density }

var stoneTypes = []struct {
	name     string
	color    string
	strength uint64
	density  uint64
}{
	{
		name:     "stone0",
		strength: 5 << 0,
		color:    "#000",
		density:  200,
	},
	{
		name:     "stone1",
		strength: 5 << 1,
		color:    "#111",
		density:  210,
	},
	{
		name:     "stone2",
		strength: 5 << 2,
		color:    "#222",
		density:  220,
	},
	{
		name:     "stone3",
		strength: 5 << 3,
		color:    "#333",
		density:  230,
	},
	{
		name:     "stone4",
		strength: 5 << 4,
		color:    "#444",
		density:  240,
	},
	{
		name:     "stone5",
		strength: 5 << 5,
		color:    "#555",
		density:  250,
	},
	{
		name:     "stone6",
		strength: 5 << 6,
		color:    "#666",
		density:  260,
	},
	{
		name:     "stone7",
		strength: 5 << 7,
		color:    "#777",
		density:  270,
	},
	{
		name:     "stone8",
		strength: 5 << 8,
		color:    "#888",
		density:  280,
	},
	{
		name:     "stone9",
		strength: 5 << 9,
		color:    "#999",
		density:  290,
	},
	{
		name:     "stone10",
		strength: 5 << 10,
		color:    "#aaa",
		density:  300,
	},
	{
		name:     "stone11",
		strength: 5 << 11,
		color:    "#bbb",
		density:  310,
	},
	{
		name:     "stone12",
		strength: 5 << 12,
		color:    "#ccc",
		density:  320,
	},
	{
		name:     "stone13",
		strength: 5 << 13,
		color:    "#ddd",
		density:  330,
	},
	{
		name:     "stone14",
		strength: 5 << 14,
		color:    "#eee",
		density:  340,
	},
	{
		name:     "stone15",
		strength: 5 << 15,
		color:    "#fff",
		density:  350,
	},
}
