package material

import (
	"github.com/BenLubar/Rnoadm/world"
	"image/color"
	"math/big"
)

type MetalType uint64

func (t MetalType) Data() *MaterialData {
	return &metalTypes[t]
}

var metalTypes = []MaterialData{
	{
		name:    "metal0",
		color0:  color.Gray{0x00},
		color1:  color.Gray{0x00},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatPower: big.NewInt(1000),
		},
	},
	{
		name:    "metal1",
		color0:  color.Gray{0x11},
		color1:  color.Gray{0x11},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatMagic: big.NewInt(1000),
		},
	},
	{
		name:    "metal2",
		color0:  color.Gray{0x22},
		color1:  color.Gray{0x22},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatAgility: big.NewInt(1000),
		},
	},
	{
		name:    "metal3",
		color0:  color.Gray{0x33},
		color1:  color.Gray{0x33},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatLuck: big.NewInt(1000),
		},
	},
	{
		name:    "metal4",
		color0:  color.Gray{0x44},
		color1:  color.Gray{0x44},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatIntelligence: big.NewInt(1000),
		},
	},
	{
		name:    "metal5",
		color0:  color.Gray{0x55},
		color1:  color.Gray{0x55},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatStamina: big.NewInt(1000),
		},
	},
	{
		name:    "metal6",
		color0:  color.Gray{0x66},
		color1:  color.Gray{0x66},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatMeleeDamage: big.NewInt(700),
			world.StatMagicDamage: big.NewInt(700),
		},
	},
	{
		name:    "metal7",
		color0:  color.Gray{0x77},
		color1:  color.Gray{0x77},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatMeleeArmor: big.NewInt(700),
			world.StatMagicArmor: big.NewInt(700),
		},
	},
	{
		name:    "metal8",
		color0:  color.Gray{0x88},
		color1:  color.Gray{0x88},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatMana:   big.NewInt(700),
			world.StatHealth: big.NewInt(700),
		},
	},
	{
		name:    "metal9",
		color0:  color.Gray{0x99},
		color1:  color.Gray{0x99},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatManaRegen:   big.NewInt(700),
			world.StatHealthRegen: big.NewInt(700),
		},
	},
	{
		name:    "metal10",
		color0:  color.Gray{0xaa},
		color1:  color.Gray{0xaa},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatCritChance: big.NewInt(1000),
		},
	},
	{
		name:    "metal11",
		color0:  color.Gray{0xbb},
		color1:  color.Gray{0xbb},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatResistance: big.NewInt(1000),
		},
	},
	{
		name:    "metal12",
		color0:  color.Gray{0xcc},
		color1:  color.Gray{0xcc},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatCritChance:  big.NewInt(700),
			world.StatAttackSpeed: big.NewInt(700),
		},
	},
	{
		name:    "metal13",
		color0:  color.Gray{0xdd},
		color1:  color.Gray{0xdd},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatResistance:    big.NewInt(700),
			world.StatMovementSpeed: big.NewInt(700),
		},
	},
	{
		name:    "metal14",
		color0:  color.Gray{0xee},
		color1:  color.Gray{0xee},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatGathering: big.NewInt(1000),
		},
	},
	{
		name:    "metal15",
		color0:  color.Gray{0xff},
		color1:  color.Gray{0xff},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatStructureHealth: big.NewInt(1000),
		},
	},
	{
		name:    "metal16",
		color0:  color.Gray{0xff},
		color1:  color.Gray{0xff},
		density: 800,

		stats: map[world.Stat]*big.Int{
			world.StatIntegrity: big.NewInt(1000),
		},
	},
}
