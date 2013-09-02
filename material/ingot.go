package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"github.com/dustin/go-humanize"
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

	for _, m := range i.materials {
		info = append(info, m.Info()...)

		info = append(info, [][2]string{
			{humanize.Comma(int64(m.metal.Strength())), "#4fc"},
			{" strength", "#ccc"},
		})
	}

	return "a bar of metal.", info
}

func (i *Ingot) Sprite() string {
	return "item_ingot"
}

func (i *Ingot) Colors() []string {
	return []string{i.materials[0].metal.Color()}
}

func (i *Ingot) Volume() uint64 {
	return 1 // TODO
}

func (i *Ingot) Weight() uint64 {
	return 0 // TODO
}

func (i *Ingot) AdminOnly() bool {
	return false
}
