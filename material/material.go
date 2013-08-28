package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"github.com/dustin/go-humanize"
	"strconv"
)

type Material struct {
	world.Object

	wood    *WoodType
	stone   *StoneType
	metal   *MetalType
	quality uint64
}

func init() {
	world.Register("mat", world.ObjectLike((*Material)(nil)))
}

func (m *Material) Save() (uint, interface{}, []world.ObjectLike) {
	materials := make(map[string]interface{})
	if m.wood != nil {
		materials["w"] = uint64(*m.wood)
	}
	if m.stone != nil {
		materials["s"] = uint64(*m.stone)
	}
	if m.metal != nil {
		materials["m"] = uint64(*m.metal)
	}
	materials["q"] = m.quality
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
		if quality, ok := materials["q"].(uint64); ok {
			m.quality = quality
		} else {
			m.quality = 1 << 62
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (m *Material) Info() [][][2]string {
	var info [][][2]string

	info = append(info, [][2]string{
		{humanize.Comma(int64(m.quality)), "#4fc"},
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

func (m *Material) Quality() uint64 {
	return m.quality
}

func WrapSpawnFunc(f func(*Material, string) world.Visible) func(string) world.Visible {
	return func(s string) world.Visible {
		prefix := func(p string) bool {
			return len(s) > len(p) && s[:len(p)] == p && s[len(p)] == ' '
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
					quality, err := strconv.ParseUint(s[1:l], 0, 64)
					if err == nil {
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
					if prefix(t.Name()) {
						material.wood = &t
						s = s[len(t.Name())+1:]
						continue find
					}
				}
			}
			if material.stone == nil {
				for i := range stoneTypes {
					t := StoneType(i)
					if prefix(t.Name()) {
						material.stone = &t
						s = s[len(t.Name())+1:]
						continue find
					}
				}
			}
			if material.metal == nil {
				for i := range metalTypes {
					t := MetalType(i)
					if prefix(t.Name()) {
						material.metal = &t
						s = s[len(t.Name())+1:]
						continue find
					}
				}
			}
			return lastGood
		}
	}
}
