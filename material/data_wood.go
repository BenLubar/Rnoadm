package material

import (
	"image/color"
	"math/big"
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
func (t WoodType) Durability() *big.Int { return woodTypes[t].durability }

var woodTypes = []struct {
	name       string
	color0     color.Color
	color1     color.Color
	leaf       uint8
	density    uint64
	durability *big.Int
}{
	{
		name:       "wood0",
		color0:     color.Gray{0x00},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood1",
		color0:     color.Gray{0x11},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood2",
		color0:     color.Gray{0x22},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood3",
		color0:     color.Gray{0x33},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood4",
		color0:     color.Gray{0x44},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood5",
		color0:     color.Gray{0x55},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood6",
		color0:     color.Gray{0x66},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood7",
		color0:     color.Gray{0x77},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood8",
		color0:     color.Gray{0x88},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood9",
		color0:     color.Gray{0x99},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood10",
		color0:     color.Gray{0xaa},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood11",
		color0:     color.Gray{0xbb},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood12",
		color0:     color.Gray{0xcc},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood13",
		color0:     color.Gray{0xdd},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood14",
		color0:     color.Gray{0xee},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood15",
		color0:     color.Gray{0xff},
		color1:     color.Alpha{0},
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood16",
		color0:     color.RGBA{0xd2, 0xb4, 0x8c, 0xff},
		color1:     color.RGBA{0x4b, 0x5a, 0x3f, 0xff},
		leaf:       1,
		density:    65,
		durability: big.NewInt(1000),
	},
	{
		name:       "wood17",
		color0:     color.RGBA{0xb5, 0xaa, 0x8b, 0xff},
		color1:     color.RGBA{0xcf, 0x51, 0x23, 0xff},
		leaf:       2,
		density:    65,
		durability: big.NewInt(1000),
	},
}
