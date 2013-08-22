package material

type StoneType uint64

func (t StoneType) Name() string     { return stoneTypes[t].name }
func (t StoneType) Color() string    { return stoneTypes[t].color }
func (t StoneType) Strength() uint64 { return stoneTypes[t].strength }

var stoneTypes = []struct {
	name     string
	color    string
	strength uint64
}{
	{
		name:     "stone0",
		strength: 5 << 0,
		color:    "#000",
	},
	{
		name:     "stone1",
		strength: 5 << 1,
		color:    "#111",
	},
	{
		name:     "stone2",
		strength: 5 << 2,
		color:    "#222",
	},
	{
		name:     "stone3",
		strength: 5 << 3,
		color:    "#333",
	},
	{
		name:     "stone4",
		strength: 5 << 4,
		color:    "#444",
	},
	{
		name:     "stone5",
		strength: 5 << 5,
		color:    "#555",
	},
	{
		name:     "stone6",
		strength: 5 << 6,
		color:    "#666",
	},
	{
		name:     "stone7",
		strength: 5 << 7,
		color:    "#777",
	},
	{
		name:     "stone8",
		strength: 5 << 8,
		color:    "#888",
	},
	{
		name:     "stone9",
		strength: 5 << 9,
		color:    "#999",
	},
	{
		name:     "stone10",
		strength: 5 << 10,
		color:    "#aaa",
	},
	{
		name:     "stone11",
		strength: 5 << 11,
		color:    "#bbb",
	},
	{
		name:     "stone12",
		strength: 5 << 12,
		color:    "#ccc",
	},
	{
		name:     "stone13",
		strength: 5 << 13,
		color:    "#ddd",
	},
	{
		name:     "stone14",
		strength: 5 << 14,
		color:    "#eee",
	},
	{
		name:     "stone15",
		strength: 5 << 15,
		color:    "#fff",
	},
}
