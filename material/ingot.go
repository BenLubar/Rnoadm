package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"math/big"
)

type Ingot struct {
	world.VisibleObject

	material *Material
}

func init() {
	world.Register("ingot", world.Visible((*Ingot)(nil)))

	world.RegisterSpawnFunc(WrapSpawnFunc(func(material *Material, s string) world.Visible {
		wood, stone, metal := material.Get()
		if len(wood) != 0 || len(stone) != 0 || len(metal) == 0 {
			return nil
		}
		if s == "ingot" {
			return &Ingot{
				material: material,
			}
		}
		return nil
	}))
}

func (i *Ingot) Save() (uint, interface{}, []world.ObjectLike) {
	return 2, uint(0), []world.ObjectLike{i.material}
}

func (i *Ingot) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		if attached[1].(*Material).Name() == "" {
			attached = attached[:1]
		}
		fallthrough
	case 1:
		material := attached[0].(*Material)
		for _, m := range attached[1:] {
			material.components = append(material.components, m.(*Material).components...)
		}
		material.sortComponents()
		attached = attached[:1]
		fallthrough
	case 2:
		i.material = attached[0].(*Material)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (i *Ingot) Name() string {
	return i.material.Name() + "ingot"
}

func (i *Ingot) Examine() (string, [][][2]string) {
	_, info := i.VisibleObject.Examine()

	info = append(info, i.material.Info()...)

	return "a bar of metal.", info
}

func (i *Ingot) Sprite() string {
	return "item_ingot"
}

func (i *Ingot) Colors() []string {
	return []string{i.material.MetalColor()}
}

func (i *Ingot) Quality() *big.Int {
	return i.material.Quality()
}

func (i *Ingot) Volume() uint64 {
	return i.material.Volume()
}

func (i *Ingot) Weight() uint64 {
	return i.material.Weight()
}

func (i *Ingot) AdminOnly() bool {
	return false
}
