package material

import (
	"image/color"
)

type WoodType uint64

// Name of the resource
func (t WoodType) Name() string { return woodTypes[t].name }

// Color of the resource (processed form)
func (t WoodType) Color() color.Color { return woodTypes[t].color0 }

// LeafColor of the resource (node)
func (t WoodType) LeafColor() color.Color { return woodTypes[t].color1 }

// LeafType of the resource (node)
func (t WoodType) LeafType() uint8 { return woodTypes[t].leaf }

// Density of the resource (centigrams per cubic centimeter)
func (t WoodType) Density() uint64 { return woodTypes[t].density }

// Durability (resistance to item degradation)
func (t WoodType) Durability() uint64 { return woodTypes[t].durability }

// Power (general stat)
func (t WoodType) Power() uint64 { return woodTypes[t].power }

// Resistance (general stat)
func (t WoodType) Resistance() uint64 { return woodTypes[t].resist }

// Swiftness (general stat)
func (t WoodType) Swiftness() uint64 { return woodTypes[t].swift }

// Spirituality (general stat)
func (t WoodType) Spirituality() uint64 { return woodTypes[t].spirit }

// MeleeDamage (offensive stat)
func (t WoodType) MeleeDamage() uint64 { return woodTypes[t].meleeDmg }

// MagicDamage (offensive stat)
func (t WoodType) MagicDamage() uint64 { return woodTypes[t].magicDmg }

// AttackSpeed (offensive stat)
func (t WoodType) AttackSpeed() uint64 { return woodTypes[t].attackSpd }

// MeleeDefense (defensive stat)
func (t WoodType) MeleeDefense() uint64 { return woodTypes[t].meleeDef }

// MagicDefense (defensive stat)
func (t WoodType) MagicDefense() uint64 { return woodTypes[t].magicDef }

// MovementSpeed (defensive stat)
func (t WoodType) MovementSpeed() uint64 { return woodTypes[t].moveSpd }

// Health (defensive stat)
func (t WoodType) Health() uint64 { return woodTypes[t].health }

// MiningStrength (tool stat)
func (t WoodType) MiningStrength() uint64 { return woodTypes[t].mining }

// ChoppingStrength (tool stat)
func (t WoodType) ChoppingStrength() uint64 { return woodTypes[t].chopping }

// GatheringSpeed (tool stat)
func (t WoodType) GatheringSpeed() uint64 { return woodTypes[t].gatherSpd }

// StructureHealth (structure stat)
func (t WoodType) StructureHealth() uint64 { return woodTypes[t].healthStrc }

var woodTypes = []struct {
	name       string
	color0     color.Color
	color1     color.Color
	leaf       uint8
	density    uint64
	durability uint64
	power      uint64
	resist     uint64
	swift      uint64
	spirit     uint64
	//element    []string
	meleeDmg   uint64
	magicDmg   uint64
	attackSpd  uint64
	meleeDef   uint64
	magicDef   uint64
	moveSpd    uint64
	health     uint64
	mining     uint64
	chopping   uint64
	gatherSpd  uint64
	healthStrc uint64
}{
	{
		name:       "wood0",
		color0:     color.Gray{0x00},
		color1:     color.Alpha{0},
		density:    300,
		durability: 0,
		power:      1,
		resist:     2,
		swift:      3,
		spirit:     4,
		//element:  {"dogs","kittens"}
		meleeDmg:   5,
		magicDmg:   6,
		attackSpd:  7,
		meleeDef:   8,
		magicDef:   9,
		moveSpd:    10,
		health:     11,
		mining:     12,
		chopping:   13,
		gatherSpd:  14,
		healthStrc: 15,
	},
	{
		name:       "wood1",
		color0:     color.Gray{0x11},
		color1:     color.Alpha{0},
		density:    310,
		durability: 1,
		power:      2,
		resist:     3,
		swift:      4,
		spirit:     5,
		//element:  {"dogs","kittens"}
		meleeDmg:   6,
		magicDmg:   7,
		attackSpd:  8,
		meleeDef:   9,
		magicDef:   10,
		moveSpd:    11,
		health:     12,
		mining:     13,
		chopping:   14,
		gatherSpd:  15,
		healthStrc: 0,
	},
	{
		name:       "wood2",
		color0:     color.Gray{0x22},
		color1:     color.Alpha{0},
		density:    320,
		durability: 2,
		power:      3,
		resist:     4,
		swift:      5,
		spirit:     6,
		//element:  {"dogs","kittens"}
		meleeDmg:   7,
		magicDmg:   8,
		attackSpd:  9,
		meleeDef:   10,
		magicDef:   11,
		moveSpd:    12,
		health:     13,
		mining:     14,
		chopping:   15,
		gatherSpd:  0,
		healthStrc: 1,
	},
	{
		name:       "wood3",
		color0:     color.Gray{0x33},
		color1:     color.Alpha{0},
		density:    330,
		durability: 3,
		power:      4,
		resist:     5,
		swift:      6,
		spirit:     7,
		//element:  {"dogs","kittens"}
		meleeDmg:   8,
		magicDmg:   9,
		attackSpd:  10,
		meleeDef:   11,
		magicDef:   12,
		moveSpd:    13,
		health:     14,
		mining:     15,
		chopping:   0,
		gatherSpd:  1,
		healthStrc: 2,
	},
	{
		name:       "wood4",
		color0:     color.Gray{0x44},
		color1:     color.Alpha{0},
		density:    340,
		durability: 4,
		power:      5,
		resist:     6,
		swift:      7,
		spirit:     8,
		//element:  {"dogs","kittens"}
		meleeDmg:   9,
		magicDmg:   10,
		attackSpd:  11,
		meleeDef:   12,
		magicDef:   13,
		moveSpd:    14,
		health:     15,
		mining:     0,
		chopping:   1,
		gatherSpd:  2,
		healthStrc: 3,
	},
	{
		name:       "wood5",
		color0:     color.Gray{0x55},
		color1:     color.Alpha{0},
		density:    350,
		durability: 5,
		power:      6,
		resist:     7,
		swift:      8,
		spirit:     9,
		//element:  {"dogs","kittens"}
		meleeDmg:   10,
		magicDmg:   11,
		attackSpd:  12,
		meleeDef:   13,
		magicDef:   14,
		moveSpd:    15,
		health:     0,
		mining:     1,
		chopping:   2,
		gatherSpd:  3,
		healthStrc: 4,
	},
	{
		name:       "wood6",
		color0:     color.Gray{0x66},
		color1:     color.Alpha{0},
		density:    360,
		durability: 6,
		power:      7,
		resist:     8,
		swift:      9,
		spirit:     10,
		//element:  {"dogs","kittens"}
		meleeDmg:   11,
		magicDmg:   12,
		attackSpd:  13,
		meleeDef:   14,
		magicDef:   15,
		moveSpd:    0,
		health:     1,
		mining:     2,
		chopping:   3,
		gatherSpd:  4,
		healthStrc: 5,
	},
	{
		name:       "wood7",
		color0:     color.Gray{0x77},
		color1:     color.Alpha{0},
		density:    370,
		durability: 7,
		power:      8,
		resist:     9,
		swift:      10,
		spirit:     11,
		//element:  {"dogs","kittens"}
		meleeDmg:   12,
		magicDmg:   13,
		attackSpd:  14,
		meleeDef:   15,
		magicDef:   0,
		moveSpd:    1,
		health:     2,
		mining:     3,
		chopping:   4,
		gatherSpd:  5,
		healthStrc: 6,
	},
	{
		name:       "wood8",
		color0:     color.Gray{0x88},
		color1:     color.Alpha{0},
		density:    380,
		durability: 8,
		power:      9,
		resist:     10,
		swift:      11,
		spirit:     12,
		//element:  {"dogs","kittens"}
		meleeDmg:   13,
		magicDmg:   14,
		attackSpd:  15,
		meleeDef:   0,
		magicDef:   1,
		moveSpd:    2,
		health:     3,
		mining:     4,
		chopping:   5,
		gatherSpd:  6,
		healthStrc: 7,
	},
	{
		name:       "wood9",
		color0:     color.Gray{0x99},
		color1:     color.Alpha{0},
		density:    390,
		durability: 9,
		power:      10,
		resist:     11,
		swift:      12,
		spirit:     13,
		//element:  {"dogs","kittens"}
		meleeDmg:   14,
		magicDmg:   15,
		attackSpd:  0,
		meleeDef:   1,
		magicDef:   2,
		moveSpd:    3,
		health:     4,
		mining:     5,
		chopping:   6,
		gatherSpd:  7,
		healthStrc: 8,
	},
	{
		name:       "wood10",
		color0:     color.Gray{0xaa},
		color1:     color.Alpha{0},
		density:    400,
		durability: 10,
		power:      11,
		resist:     12,
		swift:      13,
		spirit:     14,
		//element:  {"dogs","kittens"}
		meleeDmg:   15,
		magicDmg:   0,
		attackSpd:  1,
		meleeDef:   2,
		magicDef:   3,
		moveSpd:    4,
		health:     5,
		mining:     6,
		chopping:   7,
		gatherSpd:  8,
		healthStrc: 9,
	},
	{
		name:       "wood11",
		color0:     color.Gray{0xbb},
		color1:     color.Alpha{0},
		density:    410,
		durability: 11,
		power:      12,
		resist:     13,
		swift:      14,
		spirit:     15,
		//element:  {"dogs","kittens"}
		meleeDmg:   0,
		magicDmg:   1,
		attackSpd:  2,
		meleeDef:   3,
		magicDef:   4,
		moveSpd:    5,
		health:     6,
		mining:     7,
		chopping:   8,
		gatherSpd:  9,
		healthStrc: 10,
	},
	{
		name:       "wood12",
		color0:     color.Gray{0xcc},
		color1:     color.Alpha{0},
		density:    420, //smoke weed erryday
		durability: 12,
		power:      13,
		resist:     14,
		swift:      15,
		spirit:     0,
		//element:  {"dogs","kittens"}
		meleeDmg:   1,
		magicDmg:   2,
		attackSpd:  3,
		meleeDef:   4,
		magicDef:   5,
		moveSpd:    6,
		health:     7,
		mining:     8,
		chopping:   9,
		gatherSpd:  10,
		healthStrc: 11,
	},
	{
		name:       "wood13",
		color0:     color.Gray{0xdd},
		color1:     color.Alpha{0},
		density:    430,
		durability: 13,
		power:      14,
		resist:     15,
		swift:      0,
		spirit:     1,
		//element:  {"dogs","kittens"}
		meleeDmg:   2,
		magicDmg:   3,
		attackSpd:  4,
		meleeDef:   5,
		magicDef:   6,
		moveSpd:    7,
		health:     8,
		mining:     9,
		chopping:   10,
		gatherSpd:  11,
		healthStrc: 12,
	},
	{
		name:       "wood14",
		color0:     color.Gray{0xee},
		color1:     color.Alpha{0},
		density:    440,
		durability: 14,
		power:      15,
		resist:     0,
		swift:      1,
		spirit:     2,
		//element:  {"dogs","kittens"}
		meleeDmg:   3,
		magicDmg:   4,
		attackSpd:  5,
		meleeDef:   6,
		magicDef:   7,
		moveSpd:    8,
		health:     9,
		mining:     10,
		chopping:   11,
		gatherSpd:  12,
		healthStrc: 13,
	},
	{
		name:       "wood15",
		color0:     color.Gray{0xff},
		color1:     color.Alpha{0},
		density:    450,
		durability: 15,
		power:      0,
		resist:     1,
		swift:      2,
		spirit:     3,
		//element:  {"dogs","kittens"}
		meleeDmg:   4,
		magicDmg:   5,
		attackSpd:  6,
		meleeDef:   7,
		magicDef:   8,
		moveSpd:    9,
		health:     10,
		mining:     11,
		chopping:   12,
		gatherSpd:  13,
		healthStrc: 14,
	},
	{
		name:       "wood16",
		color0:     color.RGBA{0xd2, 0xb4, 0x8c, 0xff},
		color1:     color.RGBA{0x4b, 0x5a, 0x3f, 0xff},
		leaf:       1,
		density:    300,
		durability: 0,
		power:      1,
		resist:     2,
		swift:      3,
		spirit:     4,
		//element:  {"dogs","kittens"}
		meleeDmg:   5,
		magicDmg:   6,
		attackSpd:  7,
		meleeDef:   8,
		magicDef:   9,
		moveSpd:    10,
		health:     11,
		mining:     12,
		chopping:   13,
		gatherSpd:  14,
		healthStrc: 15,
	},
	{
		name:       "wood17",
		color0:     color.RGBA{0xb5, 0xaa, 0x8b, 0xff},
		color1:     color.RGBA{0xcf, 0x51, 0x23, 0xff},
		leaf:       2,
		density:    310,
		durability: 1,
		power:      2,
		resist:     3,
		swift:      4,
		spirit:     5,
		//element:  {"dogs","kittens"}
		meleeDmg:   6,
		magicDmg:   7,
		attackSpd:  8,
		meleeDef:   9,
		magicDef:   10,
		moveSpd:    11,
		health:     12,
		mining:     13,
		chopping:   14,
		gatherSpd:  15,
		healthStrc: 0,
	},
}
