package material

import (
	"github.com/BenLubar/Rnoadm/world"
	"image/color"
	"math/big"
)

var (
	defaultDurability = big.NewInt(1000)
	defaultStat       = big.NewInt(250)
	zero              = big.NewInt(0)
)

type MaterialData struct {
	name       string
	color0     color.Color
	color1     color.Color
	skin       uint8
	density    uint64   // centigrams per cubic centimeter
	durability *big.Int // resistance to item degradation

	stats map[world.Stat]*big.Int
}

func (m *MaterialData) Name() string {
	return m.name
}

func (m *MaterialData) Color() color.Color {
	return m.color0
}

func (m *MaterialData) ExtraColor() color.Color {
	return m.color1
}

func (m *MaterialData) Skin() uint8 {
	return m.skin
}

// Density in centigrams per cubic centimeter.
func (m *MaterialData) Density() uint64 {
	return m.density
}

// Durability increases the amount of wear an item can take before degrading.
func (m *MaterialData) Durability() *big.Int {
	if m.durability != nil {
		return m.durability
	}
	return defaultDurability
}

func (m *MaterialData) Stat(s world.Stat) *big.Int {
	if stat, ok := m.stats[s]; ok {
		return stat
	}
	if s < 1<<8 {
		return defaultStat
	}
	return zero
}
