package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"github.com/dustin/go-humanize"
	"sort"
)

type Ingot struct {
	world.VisibleObject

	materials []*Material
}

func init() {
	world.Register("ingot", world.Visible((*Ingot)(nil)))

	world.RegisterSpawnFunc(WrapSpawnFunc(func(material *Material, s string) world.Visible {
		wood, stone, metal := material.Get()
		if wood != nil || stone != nil || metal == nil {
			return nil
		}
		if s == "ingot" {
			return &Ingot{
				materials: []*Material{material},
			}
		}
		return nil
	}))
}

func (i *Ingot) Save() (uint, interface{}, []world.ObjectLike) {
	attached := []world.ObjectLike{}

	for _, m := range i.materials {
		attached = append(attached, m)
	}

	return 1, map[string]interface{}{
		"m": uint(len(i.materials)),
	}, attached
}

func (i *Ingot) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		if attached[1].(*Material).Name() == "" {
			attached = attached[:1]
		}
		data = map[string]interface{}{
			"m": uint(len(attached)),
		}
		fallthrough
	case 1:
		dataMap := data.(map[string]interface{})
		i.materials = make([]*Material, dataMap["m"].(uint))
		for j := range i.materials {
			i.materials[j] = attached[j].(*Material)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
	i.sortAlloy()
}

func (i *Ingot) Name() string {
	materials := ""
	for _, m := range i.materials {
		if name := m.Name(); name != "" {
			if materials == "" {
				materials = name
			} else {
				materials = materials[:len(materials)-1] + "-" + name
			}
		}
	}
	return materials + "ingot"
}

func (i *Ingot) Examine() (string, [][][2]string) {
	_, info := i.VisibleObject.Examine()

	info = append(info, [][2]string{
		{humanize.Comma(int64(i.Quality())), "#4fc"},
		{" quality", "#ccc"},
	}, [][2]string{
		{humanize.Comma(int64(i.Strength())), "#4fc"},
		{" strength", "#ccc"},
	})

	return "a bar of metal.", info
}

func (i *Ingot) Sprite() string {
	return "item_ingot"
}

func (i *Ingot) Colors() []string {
	return []string{i.materials[0].metal.Color()}
}

func (i *Ingot) Strength() uint64 {
	var strength uint64
	for _, m := range i.materials {
		strength += m.metal.Strength()
	}
	return strength
}

func (i *Ingot) Quality() uint64 {
	var quality uint64
	for _, m := range i.materials {
		quality += m.Quality()
	}
	return quality
}

func (i *Ingot) Volume() uint64 {
	var volume uint64
	for _, m := range i.materials {
		volume += m.Volume()
	}
	return volume
}

func (i *Ingot) Weight() uint64 {
	var weight uint64
	for _, m := range i.materials {
		weight += m.Weight()
	}
	return weight
}

func (i *Ingot) AdminOnly() bool {
	return false
}

func (i *Ingot) sortAlloy() {
	sort.Sort(sortMetals(i.materials))
}

type sortMetals []*Material

func (s sortMetals) Len() int {
	return len(s)
}

func (s sortMetals) Less(i, j int) bool {
	return *s[i].metal > *s[j].metal
}

func (s sortMetals) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
