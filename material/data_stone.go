package material

import (
	"image/color"
)

type StoneType uint64

func (t StoneType) Data() *MaterialData {
	return &stoneTypes[t]
}

var stoneTypes = []MaterialData{
	{
		name:    "stone0",
		color0:  color.Gray{0x00},
		color1:  color.Gray{0x00},
		density: 260,
	},
	{
		name:    "stone1",
		color0:  color.Gray{0x11},
		color1:  color.Gray{0x11},
		density: 260,
	},
	{
		name:    "stone2",
		color0:  color.Gray{0x22},
		color1:  color.Gray{0x22},
		density: 260,
	},
	{
		name:    "stone3",
		color0:  color.Gray{0x33},
		color1:  color.Gray{0x33},
		density: 260,
	},
	{
		name:    "stone4",
		color0:  color.Gray{0x44},
		color1:  color.Gray{0x44},
		density: 260,
	},
	{
		name:    "stone5",
		color0:  color.Gray{0x55},
		color1:  color.Gray{0x55},
		density: 260,
	},
	{
		name:    "stone6",
		color0:  color.Gray{0x66},
		color1:  color.Gray{0x66},
		density: 260,
	},
	{
		name:    "stone7",
		color0:  color.Gray{0x77},
		color1:  color.Gray{0x77},
		density: 260,
	},
	{
		name:    "stone8",
		color0:  color.Gray{0x88},
		color1:  color.Gray{0x88},
		density: 260,
	},
	{
		name:    "stone9",
		color0:  color.Gray{0x99},
		color1:  color.Gray{0x99},
		density: 260,
	},
	{
		name:    "stone10",
		color0:  color.Gray{0xaa},
		color1:  color.Gray{0xaa},
		density: 260,
	},
	{
		name:    "stone11",
		color0:  color.Gray{0xbb},
		color1:  color.Gray{0xbb},
		density: 260,
	},
	{
		name:    "stone12",
		color0:  color.Gray{0xcc},
		color1:  color.Gray{0xcc},
		density: 260,
	},
	{
		name:    "stone13",
		color0:  color.Gray{0xdd},
		color1:  color.Gray{0xdd},
		density: 260,
	},
	{
		name:    "stone14",
		color0:  color.Gray{0xee},
		color1:  color.Gray{0xee},
		density: 260,
	},
	{
		name:    "stone15",
		color0:  color.Gray{0xff},
		color1:  color.Gray{0xff},
		density: 260,
	},
}
