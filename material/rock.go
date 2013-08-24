package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
)

func init() {
	world.RegisterSpawnFunc(WrapSpawnFunc(func(material *Material, s string) world.Visible {
		wood, stone, metal := material.Get()
		if wood != nil {
			return nil
		}
		if s == "stone" && stone != nil && metal == nil {
			return &Stone{
				kind: *stone,
			}
		}
		if s == "ore" && stone == nil && metal != nil {
			return &Ore{
				kind: *metal,
			}
		}
		if stone == nil {
			return nil
		}
		if s == "rock" {
			return &Rock{
				material: material,
			}
		}
		if s == "deposit" && metal != nil {
			return &Rock{
				material: material,
				rich:     true,
			}
		}
		return nil
	}))
}

type Rock struct {
	Node

	material *Material
	rich     bool
}

func init() {
	world.Register("rock", NodeLike((*Rock)(nil)))
}

func (r *Rock) Save() (uint, interface{}, []world.ObjectLike) {
	return 0, r.rich, []world.ObjectLike{&r.Node, r.material}
}

func (r *Rock) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		r.Node = *attached[0].(*Node)
		r.material = attached[1].(*Material)
		r.rich = data.(bool)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (r *Rock) Strength() uint64 {
	_, stone, metal := r.material.Get()
	strength := stone.Strength()
	if metal != nil {
		strength += metal.Strength() / 2
	}
	return strength
}

func (r *Rock) Name() string {
	return r.material.Name() + "rock"
}

func (r *Rock) Sprite() string {
	return "rock"
}

func (r *Rock) SpriteSize() (uint, uint) {
	return 32, 32
}

func (r *Rock) Colors() []string {
	_, stone, metal := r.material.Get()
	if r.rich {
		return []string{stone.Color(), metal.OreColor(), metal.OreColor()}
	}
	if metal != nil {
		return []string{stone.Color(), metal.OreColor()}
	}
	return []string{stone.Color()}
}

type Stone struct {
	world.VisibleObject

	kind StoneType
}

func init() {
	world.Register("stone", world.Visible((*Stone)(nil)))
}

func (s *Stone) Save() (uint, interface{}, []world.ObjectLike) {
	return 0, uint64(s.kind), nil
}

func (s *Stone) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		s.kind = StoneType(data.(uint64))
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (s *Stone) Name() string {
	return s.kind.Name() + " stone"
}

func (s *Stone) Sprite() string {
	return "item_stone"
}

func (s *Stone) Colors() []string {
	return []string{s.kind.Color()}
}

func (s *Stone) Volume() uint64 {
	return 1 // TODO
}

func (s *Stone) Weight() uint64 {
	return 0 // TODO
}

func (s *Stone) AdminOnly() bool {
	return false
}

type Ore struct {
	world.VisibleObject

	kind MetalType
}

func init() {
	world.Register("ore", world.Visible((*Ore)(nil)))
}

func (o *Ore) Save() (uint, interface{}, []world.ObjectLike) {
	return 0, uint64(o.kind), nil
}

func (o *Ore) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		o.kind = MetalType(data.(uint64))
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (o *Ore) Name() string {
	return o.kind.Name() + " ore"
}

func (o *Ore) Sprite() string {
	return "item_ore"
}

func (o *Ore) Colors() []string {
	return []string{o.kind.OreColor()}
}

func (o *Ore) Volume() uint64 {
	return 1 // TODO
}

func (o *Ore) Weight() uint64 {
	return 0 // TODO
}

func (o *Ore) AdminOnly() bool {
	return false
}
