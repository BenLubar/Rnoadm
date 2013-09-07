package material

import (
	"image/color"
	"math/big"
)

type MetalType uint64

// Name of the resource
func (t MetalType) Name() string { return metalTypes[t].name }

// Color of the resource (processed form)
func (t MetalType) Color() color.Color { return metalTypes[t].color0 }

// OreColor of the resource (unprocessed form)
func (t MetalType) OreColor() color.Color { return metalTypes[t].color1 }

// Density of the resource (centigrams per cubic centimeter)
func (t MetalType) Density() uint64 { return metalTypes[t].density }

// Durability (resistance to item degradation)
func (t MetalType) Durability() *big.Int { return metalTypes[t].durability }

var metalTypes = []struct {
	name       string
	color0     color.Color
	color1     color.Color
	density    uint64
	durability *big.Int
}{
	{
		name:       "metal0",
		color0:     color.Gray{0x00},
		color1:     color.Gray{0x00},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal1",
		color0:     color.Gray{0x11},
		color1:     color.Gray{0x11},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal2",
		color0:     color.Gray{0x22},
		color1:     color.Gray{0x22},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal3",
		color0:     color.Gray{0x33},
		color1:     color.Gray{0x33},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal4",
		color0:     color.Gray{0x44},
		color1:     color.Gray{0x44},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal5",
		color0:     color.Gray{0x55},
		color1:     color.Gray{0x55},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal6",
		color0:     color.Gray{0x66},
		color1:     color.Gray{0x66},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal7",
		color0:     color.Gray{0x77},
		color1:     color.Gray{0x77},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal8",
		color0:     color.Gray{0x88},
		color1:     color.Gray{0x88},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal9",
		color0:     color.Gray{0x99},
		color1:     color.Gray{0x99},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal10",
		color0:     color.Gray{0xaa},
		color1:     color.Gray{0xaa},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal11",
		color0:     color.Gray{0xbb},
		color1:     color.Gray{0xbb},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal12",
		color0:     color.Gray{0xcc},
		color1:     color.Gray{0xcc},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal13",
		color0:     color.Gray{0xdd},
		color1:     color.Gray{0xdd},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal14",
		color0:     color.Gray{0xee},
		color1:     color.Gray{0xee},
		density:    800,
		durability: big.NewInt(1000),
	},
	{
		name:       "metal15",
		color0:     color.Gray{0xff},
		color1:     color.Gray{0xff},
		density:    800,
		durability: big.NewInt(1000),
	},
}
