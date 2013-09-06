package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"image/color"
	"math/big"
	"strconv"
)

type Material struct {
	world.Object

	wood    *WoodType
	stone   *StoneType
	metal   *MetalType
	quality big.Int

	woodVolume  uint64
	stoneVolume uint64
	metalVolume uint64
}

func init() {
	world.Register("mat", world.ObjectLike((*Material)(nil)))
}

func (m *Material) Save() (uint, interface{}, []world.ObjectLike) {
	materials := map[string]interface{}{
		"q":  &m.quality,
		"vw": m.woodVolume,
		"vs": m.stoneVolume,
		"vm": m.metalVolume,
	}
	if m.wood != nil {
		materials["w"] = uint64(*m.wood)
	}
	if m.stone != nil {
		materials["s"] = uint64(*m.stone)
	}
	if m.metal != nil {
		materials["m"] = uint64(*m.metal)
	}
	return 0, materials, nil
}

func (m *Material) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		materials, _ := data.(map[string]interface{})
		if wood, ok := materials["w"].(uint64); ok {
			m.wood = (*WoodType)(&wood)
		}
		if stone, ok := materials["s"].(uint64); ok {
			m.stone = (*StoneType)(&stone)
		}
		if metal, ok := materials["m"].(uint64); ok {
			m.metal = (*MetalType)(&metal)
		}
		if quality, ok := materials["q"].(*big.Int); ok {
			m.quality = *quality
		} else if quality, ok := materials["q"].(uint64); ok {
			m.quality = *big.NewInt(int64(quality))
		} else {
			m.quality = *big.NewInt(1 << 62)
		}
		if _, ok := materials["vw"]; ok {
			m.woodVolume = materials["vw"].(uint64)
			m.stoneVolume = materials["vs"].(uint64)
			m.metalVolume = materials["vm"].(uint64)
		} else {
			if m.wood != nil {
				m.woodVolume = 100
			}
			if m.stone != nil {
				m.stoneVolume = 100
			}
			if m.metal != nil {
				m.metalVolume = 100
			}
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (m *Material) Info() [][][2]string {
	var info [][][2]string

	info = append(info, [][2]string{
		{Comma(&m.quality), "#4fc"},
		{" quality", "#ccc"},
	})

	return info
}

func (m *Material) Wood() (WoodType, bool) {
	if m.wood == nil {
		return 0, false
	}
	return *m.wood, true
}

func (m *Material) Stone() (StoneType, bool) {
	if m.stone == nil {
		return 0, false
	}
	return *m.stone, true
}

func (m *Material) Metal() (MetalType, bool) {
	if m.metal == nil {
		return 0, false
	}
	return *m.metal, true
}

func (m *Material) Weight() uint64 {
	var weight uint64
	if m.wood != nil {
		weight += m.woodVolume * m.wood.Density() / 100
	}
	if m.stone != nil {
		weight += m.stoneVolume * m.stone.Density() / 100
	}
	if m.metal != nil {
		weight += m.metalVolume * m.metal.Density() / 100
	}
	return weight
}

func (m *Material) Volume() uint64 {
	return m.woodVolume + m.stoneVolume + m.metalVolume
}

func (m *Material) Name() string {
	wood, stone, metal := m.Get()
	if wood == nil {
		if stone == nil {
			if metal == nil {
				return ""
			}
			return metal.Name() + " "
		}
		if metal == nil {
			return stone.Name() + " "
		}
		return metal.Name() + " and " + stone.Name() + " "
	}
	if stone == nil {
		if metal == nil {
			return wood.Name() + " "
		}
		return wood.Name() + " and " + metal.Name() + " "
	}
	if metal == nil {
		return wood.Name() + " and " + stone.Name() + " "
	}
	return wood.Name() + ", " + metal.Name() + ", and " + stone.Name() + " "
}

func (m *Material) Get() (*WoodType, *StoneType, *MetalType) {
	return m.wood, m.stone, m.metal
}

func (m *Material) Quality() *big.Int {
	return &m.quality
}

func WrapSpawnFunc(f func(*Material, string) world.Visible) func(string) world.Visible {
	return func(s string) world.Visible {
		prefix := func(p string) (bool, uint64) {
			if len(s) > len(p) && s[:len(p)] == p {
				if s[len(p)] == ' ' {
					s = s[len(p)+1:]
					return true, 100
				}
				if s[len(p)] == ':' {
					l := len(p) + 1
					for ; l < len(s)-2 && s[l] != ' '; l++ {
					}
					if s[l] == ' ' {
						volume, err := strconv.ParseUint(s[len(p)+1:l], 0, 64)
						if err == nil {
							s = s[l+1:]
							return true, volume
						}
					}
				}
			}
			return false, 0
		}

		var lastGood world.Visible
		var setQuality bool
		var material Material
		world.InitObject(&material)
	find:
		for {
			if v := f(&material, s); v != nil {
				lastGood = v
			}
			if !setQuality && len(s) > 3 && s[0] == 'q' {
				var l int
				for l = 0; l < len(s)-2 && s[l] != ' '; l++ {
				}
				if s[l] == ' ' {
					var quality big.Int
					_, ok := quality.SetString(s[1:l], 10)
					if ok {
						material.quality = quality
						setQuality = true
						s = s[l+1:]
						continue
					}
				}
			}
			if material.wood == nil {
				for i := range woodTypes {
					t := WoodType(i)
					if ok, volume := prefix(t.Name()); ok {
						material.wood = &t
						material.woodVolume = volume
						continue find
					}
				}
			}
			if material.stone == nil {
				for i := range stoneTypes {
					t := StoneType(i)
					if ok, volume := prefix(t.Name()); ok {
						material.stone = &t
						material.stoneVolume = volume
						continue find
					}
				}
			}
			if material.metal == nil {
				for i := range metalTypes {
					t := MetalType(i)
					if ok, volume := prefix(t.Name()); ok {
						material.metal = &t
						material.metalVolume = volume
						continue find
					}
				}
			}
			return lastGood
		}
	}
}

func toCSSColor(c color.Color) string {
	const m = 1<<16 - 1
	r, g, b, a := c.RGBA()
	if a == m {
		return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
	}
	return fmt.Sprintf("rgba(%d,%d,%d,%f)", (r*m/a)>>8, (g*m/a)>>8, (b*m/a)>>8, float64(a)/m)
}
