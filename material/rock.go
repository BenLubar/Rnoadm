package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"math/big"
)

func init() {
	world.RegisterSpawnFunc(WrapSpawnFunc(func(material *Material, s string) world.Visible {
		wood, stone, metal := material.Get()
		if len(wood) != 0 {
			return nil
		}
		if s == "stone" && len(stone) != 0 && len(metal) == 0 {
			return &Stone{
				material: material,
			}
		}
		if s == "ore" && len(stone) == 0 && len(metal) != 0 {
			return &Ore{
				material: material,
			}
		}
		if len(stone) == 0 {
			return nil
		}
		if s == "rock" {
			return &Rock{
				material: material,
			}
		}
		if s == "deposit" && len(metal) != 0 {
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

func (r *Rock) Quality() *big.Int {
	return r.material.Quality()
}

func (r *Rock) Name() string {
	if r.rich {
		return r.material.Name() + "deposit"
	}
	return r.material.Name() + "rock"
}

func (r *Rock) Examine() (string, [][][2]string) {
	_, info := r.Node.Examine()

	info = append(info, r.material.Info()...)

	return "a rock.", info
}

func (r *Rock) Sprite() string {
	return "rock"
}

func (r *Rock) SpriteSize() (uint, uint) {
	return 32, 32
}

func (r *Rock) Colors() []string {
	stone, ore := r.material.StoneColor(), r.material.OreColor()
	if r.rich {
		return []string{stone, ore, ore}
	}
	return []string{stone, ore}
}

func (r *Rock) Actions(player world.PlayerLike) []string {
	actions := r.Node.Actions(player)

	actions = append([]string{"quarry"}, actions...)
	if r.material.MetalColor() != "" {
		actions = append([]string{"mine"}, actions...)
	}

	return actions
}

type Stone struct {
	world.VisibleObject

	material *Material
}

func init() {
	world.Register("stone", world.Visible((*Stone)(nil)))
}

func (s *Stone) Save() (uint, interface{}, []world.ObjectLike) {
	return 1, uint(0), []world.ObjectLike{s.material}
}

func (s *Stone) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		material := &Material{components: []*material{&material{}}}
		world.InitObject(material)
		world.InitObject(material.components[0])
		kind := StoneType(data.(uint64))
		material.components[0].stone = &kind
		material.components[0].volume = 100
		material.quality = *big.NewInt(1 << 62)
		attached = append(attached, material)
		fallthrough
	case 1:
		s.material = attached[0].(*Material)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (s *Stone) Material() *Material {
	return s.material
}

func (s *Stone) Name() string {
	return s.material.Name() + "stone"
}

func (s *Stone) Examine() (string, [][][2]string) {
	_, info := s.VisibleObject.Examine()

	info = append(info, s.material.Info()...)

	return "some stones.", info
}

func (s *Stone) Sprite() string {
	return "item_stone"
}

func (s *Stone) Colors() []string {
	return []string{s.material.StoneColor()}
}

func (s *Stone) Volume() uint64 {
	return s.material.Volume()
}

func (s *Stone) Weight() uint64 {
	return s.material.Weight()
}

func (s *Stone) AdminOnly() bool {
	return false
}

type Ore struct {
	world.VisibleObject

	material *Material
}

func init() {
	world.Register("ore", world.Visible((*Ore)(nil)))
}

func (o *Ore) Save() (uint, interface{}, []world.ObjectLike) {
	return 1, uint(0), []world.ObjectLike{o.material}
}

func (o *Ore) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		material := &Material{components: []*material{&material{}}}
		world.InitObject(material)
		world.InitObject(material.components[0])
		kind := MetalType(data.(uint64))
		material.components[0].metal = &kind
		material.components[0].volume = 100
		material.quality = *big.NewInt(1 << 62)
		attached = append(attached, material)
		fallthrough
	case 1:
		o.material = attached[0].(*Material)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (o *Ore) Material() *Material {
	return o.material
}

func (o *Ore) Name() string {
	return o.material.Name() + "ore"
}

func (o *Ore) Examine() (string, [][][2]string) {
	_, info := o.VisibleObject.Examine()

	info = append(info, o.material.Info()...)

	return "some unrefined ore.", info
}

func (o *Ore) Sprite() string {
	return "item_ore"
}

func (o *Ore) Colors() []string {
	return []string{o.material.OreColor()}
}

func (o *Ore) Volume() uint64 {
	return o.material.Volume()
}

func (o *Ore) Weight() uint64 {
	return o.material.Weight()
}

func (o *Ore) AdminOnly() bool {
	return false
}
