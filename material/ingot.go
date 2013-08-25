package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"github.com/dustin/go-humanize"
)

type Ingot struct {
	world.VisibleObject

	material *Material
	alloy    *Material
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
				material: material,
				alloy:    &Material{},
			}
		}
		return nil
	}))
}

func (i *Ingot) Save() (uint, interface{}, []world.ObjectLike) {
	return 0, uint(0), []world.ObjectLike{i.material, i.alloy}
}

func (i *Ingot) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		i.material = attached[0].(*Material)
		i.alloy = attached[1].(*Material)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (i *Ingot) Name() string {
	return i.material.Name() + i.alloy.Name() + "ingot"
}

func (i *Ingot) Examine() (string, [][][2]string) {
	_, info := i.VisibleObject.Examine()

	info = append(info, i.material.Info()...)

	info = append(info, [][2]string{
		{humanize.Comma(int64(i.material.metal.Strength())), "#4fc"},
		{" strength", "#ccc"},
	})

	return "a bar of metal.", info
}

func (i *Ingot) Sprite() string {
	return "item_ingot"
}

func (i *Ingot) Colors() []string {
	return []string{i.material.metal.Color()}
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
