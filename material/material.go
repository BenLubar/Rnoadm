package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"image/color"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

type material struct {
	world.Object

	wood   *WoodType
	stone  *StoneType
	metal  *MetalType
	volume uint64
}

func init() {
	world.Register("mc", world.ObjectLike((*material)(nil)))
}
func (m *material) Save() (uint, interface{}, []world.ObjectLike) {
	data := map[string]interface{}{"v": m.volume}
	switch {
	case m.wood != nil:
		data["w"] = uint64(*m.wood)
	case m.stone != nil:
		data["s"] = uint64(*m.stone)
	case m.metal != nil:
		data["m"] = uint64(*m.metal)
	}
	return 0, data, nil
}

func (m *material) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		m.volume = dataMap["v"].(uint64)
		if wood, ok := dataMap["w"].(uint64); ok {
			m.wood = (*WoodType)(&wood)
		}
		if stone, ok := dataMap["s"].(uint64); ok {
			m.stone = (*StoneType)(&stone)
		}
		if metal, ok := dataMap["m"].(uint64); ok {
			m.metal = (*MetalType)(&metal)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (m *material) data() *MaterialData {
	switch {
	case m.wood != nil:
		return m.wood.Data()
	case m.stone != nil:
		return m.stone.Data()
	case m.metal != nil:
		return m.metal.Data()
	}
	panic("untyped material")
}

type Material struct {
	world.Object

	components []*material
	quality    big.Int
}

func init() {
	world.Register("mat", world.ObjectLike((*Material)(nil)))
}

func (m *Material) Save() (uint, interface{}, []world.ObjectLike) {
	attached := make([]world.ObjectLike, len(m.components))
	for i, c := range m.components {
		attached[i] = c
	}
	return 1, map[string]interface{}{"q": &m.quality}, attached
}

func (m *Material) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		materials, _ := data.(map[string]interface{})
		if materials == nil {
			materials = make(map[string]interface{}, 1)
			data = materials
		}

		if wood, ok := materials["w"].(uint64); ok {
			vol, ok := materials["vw"].(uint64)
			if !ok {
				vol = 100
			}
			attached = append(attached, world.InitObject(&material{
				wood:   (*WoodType)(&wood),
				volume: vol,
			}))
		}
		if stone, ok := materials["s"].(uint64); ok {
			vol, ok := materials["vs"].(uint64)
			if !ok {
				vol = 100
			}
			attached = append(attached, world.InitObject(&material{
				stone:  (*StoneType)(&stone),
				volume: vol,
			}))
		}
		if metal, ok := materials["m"].(uint64); ok {
			vol, ok := materials["vm"].(uint64)
			if !ok {
				vol = 100
			}
			attached = append(attached, world.InitObject(&material{
				metal:  (*MetalType)(&metal),
				volume: vol,
			}))
		}

		if _, ok := materials["q"].(*big.Int); ok {
			// do nothing
		} else if quality, ok := materials["q"].(uint64); ok {
			materials["q"] = big.NewInt(int64(quality))
		} else {
			materials["q"] = big.NewInt(1 << 62)
		}
		fallthrough
	case 1:
		dataMap := data.(map[string]interface{})
		m.quality = *dataMap["q"].(*big.Int)
		m.components = make([]*material, len(attached))
		for i, c := range attached {
			m.components[i] = c.(*material)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
	m.sortComponents()
}

func (m *Material) Info() [][][2]string {
	var info [][][2]string

	maybe := func(name string, stat world.Stat) {
		var total, volume, tmp big.Int
		for _, c := range m.components {
			tmp.SetUint64(c.volume)
			volume.Add(&volume, &tmp)
			total.Add(&total, tmp.Mul(c.data().Stat(stat), &tmp))
		}
		if volume.Sign() == 0 {
			return
		}
		total.Div(tmp.Mul(&total, &m.quality), &volume)
		switch total.Sign() {
		case 0:
			return
		case 1:
			info = append(info, [][2]string{
				{"+" + Comma(&total), "#4f4"},
				{name, "#ccc"},
			})
		case -1:
			info = append(info, [][2]string{
				{Comma(&total), "#f44"},
				{name, "#ccc"},
			})
		}
	}

	maybe(" power", world.StatPower)
	maybe(" magic", world.StatMagic)
	maybe(" agility", world.StatAgility)
	maybe(" luck", world.StatLuck)
	maybe(" intelligence", world.StatIntelligence)
	maybe(" stamina", world.StatStamina)
	maybe(" integrity", world.StatIntegrity)

	maybe(" melee damage", world.StatMeleeDamage)
	maybe(" magic damage", world.StatMagicDamage)
	maybe(" mana", world.StatMana)
	maybe(" mana regen", world.StatManaRegen)
	maybe(" crit chance", world.StatCritChance)
	maybe(" attack speed", world.StatAttackSpeed)

	maybe(" melee armor", world.StatMeleeArmor)
	maybe(" magic armor", world.StatMagicArmor)
	maybe(" health", world.StatHealth)
	maybe(" health regen", world.StatHealthRegen)
	maybe(" resistance", world.StatResistance)
	maybe(" movement speed", world.StatMovementSpeed)

	maybe(" gathering", world.StatGathering)
	maybe(" structure health", world.StatStructureHealth)

	return info
}

func (m *Material) Weight() uint64 {
	var weight uint64
	for _, c := range m.components {
		weight += c.volume * c.data().Density() / 100
	}
	return weight
}

func (m *Material) Volume() uint64 {
	var volume uint64
	for _, c := range m.components {
		volume += c.volume
	}
	return volume
}

func (m *Material) Name() string {
	var names []string
	wood, stone, metal := m.Get()

	for _, w := range wood {
		names = append(names, w.Data().Name())
	}
	for _, s := range stone {
		names = append(names, s.Data().Name())
	}
	for _, m := range metal {
		names = append(names, m.Data().Name())
	}
	name := strings.Join(names, "-")
	if name == "" {
		return ""
	}
	return name + " "
}

func (m *Material) Get() (wood []WoodType, stone []StoneType, metal []MetalType) {
	for _, c := range m.components {
		switch {
		case c.wood != nil:
			wood = append(wood, *c.wood)
		case c.stone != nil:
			stone = append(stone, *c.stone)
		case c.metal != nil:
			metal = append(metal, *c.metal)
		}
	}
	return wood, stone, metal
}

func (m *Material) Quality() *big.Int {
	return &m.quality
}

func (m *Material) Copy(volume uint64) *Material {
	return m.copy(volume, func(c *material) bool {
		return true
	})
}

func (m *Material) CopyWood(volume uint64) *Material {
	return m.copy(volume, func(c *material) bool {
		return c.wood != nil
	})
}

func (m *Material) CopyStone(volume uint64) *Material {
	return m.copy(volume, func(c *material) bool {
		return c.stone != nil
	})
}

func (m *Material) CopyMetal(volume uint64) *Material {
	return m.copy(volume, func(c *material) bool {
		return c.metal != nil
	})
}

func (m *Material) copy(volume uint64, filter func(*material) bool) *Material {
	copy := &Material{
		quality: *(&big.Int{}).Set(&m.quality),
	}
	world.InitObject(copy)
	var total uint64
	for _, c := range m.components {
		if filter(c) {
			total += c.volume
			copyc := *c
			copy.components = append(copy.components, &copyc)
		}
	}

	if total == volume {
		return copy
	}

	for _, c := range copy.components {
		v := c.volume * volume / total
		total -= c.volume
		volume -= v
		c.volume = v
	}

	return copy
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

		var m Material
		world.InitObject(&m)
		m.quality.SetUint64(1000)
		var setQuality bool
	find:
		for {
			m.sortComponents()
			if v := f(&m, s); v != nil {
				return v
			}
			if !setQuality && len(s) > 3 && s[0] == 'q' {
				var l int
				for l = 0; l < len(s)-2 && s[l] != ' '; l++ {
				}
				if s[l] == ' ' {
					var quality big.Int
					_, ok := quality.SetString(s[1:l], 10)
					if ok {
						m.quality = quality
						setQuality = true
						s = s[l+1:]
						continue
					}
				}
			}
			for i := range woodTypes {
				t := WoodType(i)
				if ok, volume := prefix(t.Data().Name()); ok {
					m.components = append(m.components, &material{
						wood:   &t,
						volume: volume,
					})
					continue find
				}
			}
			for i := range stoneTypes {
				t := StoneType(i)
				if ok, volume := prefix(t.Data().Name()); ok {
					m.components = append(m.components, &material{
						stone:  &t,
						volume: volume,
					})
					continue find
				}
			}
			for i := range metalTypes {
				t := MetalType(i)
				if ok, volume := prefix(t.Data().Name()); ok {
					m.components = append(m.components, &material{
						metal:  &t,
						volume: volume,
					})
					continue find
				}
			}
			return nil
		}
	}
}

type sortMaterial []*material

func (s sortMaterial) Len() int {
	return len(s)
}

func (s sortMaterial) Less(i, j int) bool {
	a, b := s[i], s[j]
	switch {
	case a.wood != nil && b.wood == nil:
		return true
	case a.wood == nil && b.wood != nil:
		return false
	case a.wood != nil && b.wood != nil:
		return *a.wood < *b.wood

	case a.stone != nil && b.stone == nil:
		return true
	case a.stone == nil && b.stone != nil:
		return false
	case a.stone != nil && b.stone != nil:
		return *a.stone < *b.stone

	case a.metal != nil && b.metal == nil:
		return true
	case a.metal == nil && b.metal != nil:
		return false
	case a.metal != nil && b.metal != nil:
		return *a.metal < *b.metal
	}
	return false
}

func (s sortMaterial) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (m *Material) sortComponents() {
	sort.Sort(sortMaterial(m.components))
}

func (m *Material) WoodColor() string {
	var R, G, B, A uint64
	var totalVolume uint64
	for _, c := range m.components {
		if c.wood == nil {
			continue
		}
		r, g, b, a := c.wood.Data().Color().RGBA()
		totalVolume += c.volume
		R += uint64(r) * c.volume
		G += uint64(g) * c.volume
		B += uint64(b) * c.volume
		A += uint64(a) * c.volume
	}
	if totalVolume == 0 || A == 0 {
		return ""
	}
	return toCSSColor(color.RGBA64{
		uint16(R / totalVolume),
		uint16(G / totalVolume),
		uint16(B / totalVolume),
		uint16(A / totalVolume),
	})
}

func (m *Material) LeafColor() string {
	var R, G, B, A uint64
	var totalVolume uint64
	for _, c := range m.components {
		if c.wood == nil {
			continue
		}
		r, g, b, a := c.wood.Data().ExtraColor().RGBA()
		totalVolume += c.volume
		R += uint64(r) * c.volume
		G += uint64(g) * c.volume
		B += uint64(b) * c.volume
		A += uint64(a) * c.volume
	}
	if totalVolume == 0 || A == 0 {
		return ""
	}
	return toCSSColor(color.RGBA64{
		uint16(R / totalVolume),
		uint16(G / totalVolume),
		uint16(B / totalVolume),
		uint16(A / totalVolume),
	})
}

func (m *Material) StoneColor() string {
	var R, G, B, A uint64
	var totalVolume uint64
	for _, c := range m.components {
		if c.stone == nil {
			continue
		}
		r, g, b, a := c.stone.Data().Color().RGBA()
		totalVolume += c.volume
		R += uint64(r) * c.volume
		G += uint64(g) * c.volume
		B += uint64(b) * c.volume
		A += uint64(a) * c.volume
	}
	if totalVolume == 0 || A == 0 {
		return ""
	}
	return toCSSColor(color.RGBA64{
		uint16(R / totalVolume),
		uint16(G / totalVolume),
		uint16(B / totalVolume),
		uint16(A / totalVolume),
	})
}

func (m *Material) MetalColor() string {
	var R, G, B, A uint64
	var totalVolume uint64
	for _, c := range m.components {
		if c.metal == nil {
			continue
		}
		r, g, b, a := c.metal.Data().Color().RGBA()
		totalVolume += c.volume
		R += uint64(r) * c.volume
		G += uint64(g) * c.volume
		B += uint64(b) * c.volume
		A += uint64(a) * c.volume
	}
	if totalVolume == 0 || A == 0 {
		return ""
	}
	return toCSSColor(color.RGBA64{
		uint16(R / totalVolume),
		uint16(G / totalVolume),
		uint16(B / totalVolume),
		uint16(A / totalVolume),
	})
}

func (m *Material) OreColor() string {
	var R, G, B, A uint64
	var totalVolume uint64
	for _, c := range m.components {
		if c.metal == nil {
			continue
		}
		r, g, b, a := c.metal.Data().ExtraColor().RGBA()
		totalVolume += c.volume
		R += uint64(r) * c.volume
		G += uint64(g) * c.volume
		B += uint64(b) * c.volume
		A += uint64(a) * c.volume
	}
	if totalVolume == 0 || A == 0 {
		return ""
	}
	return toCSSColor(color.RGBA64{
		uint16(R / totalVolume),
		uint16(G / totalVolume),
		uint16(B / totalVolume),
		uint16(A / totalVolume),
	})
}

func toCSSColor(c color.Color) string {
	const m = 1<<16 - 1
	r, g, b, a := c.RGBA()
	if a == m {
		return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
	}
	return fmt.Sprintf("rgba(%d,%d,%d,%f)", (r*m/a)>>8, (g*m/a)>>8, (b*m/a)>>8, float64(a)/m)
}

var metaStats = map[world.Stat]world.Stat{
	world.StatMeleeDamage: world.StatPower,
	world.StatMagicDamage: world.StatMagic,
	world.StatMana:        world.StatIntelligence,
	world.StatManaRegen:   world.StatIntelligence,
	world.StatCritChance:  world.StatLuck,
	world.StatAttackSpeed: world.StatAgility,

	world.StatMeleeArmor:    world.StatPower,
	world.StatMagicArmor:    world.StatMagic,
	world.StatHealth:        world.StatStamina,
	world.StatHealthRegen:   world.StatStamina,
	world.StatResistance:    world.StatLuck,
	world.StatMovementSpeed: world.StatAgility,

	world.StatGathering:       world.StatIntegrity,
	world.StatStructureHealth: world.StatIntegrity,
}

func (m *Material) stat(stat, meta *world.Stat) *big.Int {
	var total, volume, tmp big.Int
	for _, c := range m.components {
		tmp.SetUint64(c.volume)
		volume.Add(&volume, &tmp)
		total.Add(&total, tmp.Mul(c.data().Stat(*stat), &tmp))
		if meta != nil {
			total.Add(&total, tmp.Div(tmp.Mul(c.data().Stat(*meta), tmp.SetUint64(c.volume)), world.TuningMetaStatDivisor))
		}
	}
	if volume.Sign() == 0 {
		return &volume
	}
	total.Div(tmp.Mul(&total, &m.quality), &volume)
	return &total
}

func (m *Material) Stat(stat world.Stat) *big.Int {
	var meta *world.Stat
	if m, ok := metaStats[stat]; ok {
		meta = &m
	}
	return m.stat(&stat, meta)
}
