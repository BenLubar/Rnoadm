package material

import (
	"image/color"
)

type WoodType uint64

func (t WoodType) Data() *MaterialData {
	return &woodTypes[t]
}

var woodTypes = []MaterialData{
	{
		name:    "wood0",
		color0:  color.Gray{0x00},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood1",
		color0:  color.Gray{0x11},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood2",
		color0:  color.Gray{0x22},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood3",
		color0:  color.Gray{0x33},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood4",
		color0:  color.Gray{0x44},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood5",
		color0:  color.Gray{0x55},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood6",
		color0:  color.Gray{0x66},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood7",
		color0:  color.Gray{0x77},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood8",
		color0:  color.Gray{0x88},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood9",
		color0:  color.Gray{0x99},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood10",
		color0:  color.Gray{0xaa},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood11",
		color0:  color.Gray{0xbb},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood12",
		color0:  color.Gray{0xcc},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood13",
		color0:  color.Gray{0xdd},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood14",
		color0:  color.Gray{0xee},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood15",
		color0:  color.Gray{0xff},
		color1:  color.Alpha{0},
		density: 65,
	},
	{
		name:    "wood16",
		color0:  color.RGBA{0xd2, 0xb4, 0x8c, 0xff},
		color1:  color.RGBA{0x4b, 0x5a, 0x3f, 0xff},
		skin:    1,
		density: 65,
	},
	{
		name:    "wood17",
		color0:  color.RGBA{0xb5, 0xaa, 0x8b, 0xff},
		color1:  color.RGBA{0xcf, 0x51, 0x23, 0xff},
		skin:    2,
		density: 65,
	},
}
