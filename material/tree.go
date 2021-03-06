package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"math/big"
)

func init() {
	world.RegisterSpawnFunc(WrapSpawnFunc(func(material *Material, s string) world.Visible {
		if wood, stone, metal := material.Get(); len(wood) == 0 || len(stone) != 0 || len(metal) != 0 {
			return nil
		}
		switch s {
		case "tree":
			return &Tree{material: material}
		case "logs":
			return &Logs{material: material}
		}
		return nil
	}))
}

type Tree struct {
	Node

	material *Material
}

func init() {
	world.Register("tree", NodeLike((*Tree)(nil)))
}

func (t *Tree) Save() (uint, interface{}, []world.ObjectLike) {
	return 1, uint(0), []world.ObjectLike{&t.Node, t.material}
}

func (t *Tree) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		material := &Material{components: []*material{&material{}}}
		world.InitObject(material)
		world.InitObject(material.components[0])
		kind := WoodType(data.(uint64))
		material.components[0].wood = &kind
		material.components[0].volume = 100
		material.quality = *big.NewInt(1 << 62)
		attached = append(attached, material)
		fallthrough
	case 1:
		t.Node = *attached[0].(*Node)
		t.material = attached[1].(*Material)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (t *Tree) Quality() *big.Int {
	return t.material.Quality()
}

func (t *Tree) Name() string {
	return t.material.Name() + "tree"
}

func (t *Tree) Examine() (string, [][][2]string) {
	_, info := t.Node.Examine()

	info = append(info, t.material.Info()...)

	return "a tree.", info
}

func (t *Tree) Sprite() string {
	return "tree"
}

func (t *Tree) SpriteSize() (uint, uint) {
	return 96, 128
}

func (t *Tree) Colors() []string {
	switch t.material.components[0].wood.Data().Skin() {
	default:
		fallthrough
	case 0: // no leaves
		return []string{t.material.WoodColor()}
	case 1: // deciduous
		return []string{t.material.WoodColor(), t.material.LeafColor()}
	case 2: // coniferous
		return []string{t.material.WoodColor(), "", t.material.LeafColor()}
	}
}

func (t *Tree) Actions(player world.PlayerLike) []string {
	actions := t.Node.Actions(player)

	actions = append([]string{"chop"}, actions...)

	return actions
}

func (t *Tree) Interact(player world.PlayerLike, action string) {
	switch action {
	default:
		t.Node.Interact(player, action)
	case "chop":
		pos := t.Position()
		if pos == nil {
			return
		}
		x, y := pos.Position()
		player.SetSchedule(&world.ScheduleSchedule{
			Schedules: []world.Schedule{
				world.NewWalkSchedule(x, y, true, 0),
				&GatherSchedule{
					Tool:    player,
					Target_: t,
					Item: func(volume uint64) world.Visible {
						return world.InitObject(&Logs{
							material: t.material.CopyWood(volume),
						}).(world.Visible)
					},
				},
			},
		})
	}
}

type Logs struct {
	world.VisibleObject

	material *Material
}

func init() {
	world.Register("logs", world.Visible((*Logs)(nil)))
}

func (l *Logs) Save() (uint, interface{}, []world.ObjectLike) {
	return 0, uint(0), []world.ObjectLike{l.material}
}

func (l *Logs) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		material := &Material{components: []*material{&material{}}}
		world.InitObject(material)
		world.InitObject(material.components[0])
		kind := WoodType(data.(uint64))
		material.components[0].wood = &kind
		material.components[0].volume = 100
		material.quality = *big.NewInt(1 << 62)
		attached = append(attached, material)
		fallthrough
	case 1:
		l.material = attached[0].(*Material)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (l *Logs) Material() *Material {
	return l.material
}

func (l *Logs) Name() string {
	return l.material.Name() + "logs"
}

func (l *Logs) Examine() (string, [][][2]string) {
	_, info := l.VisibleObject.Examine()

	info = append(info, l.material.Info()...)

	return "some logs.", info
}

func (l *Logs) Sprite() string {
	return "item_logs"
}

func (l *Logs) Colors() []string {
	return []string{l.material.WoodColor()}
}

func (l *Logs) Volume() uint64 {
	return l.material.Volume()
}

func (l *Logs) Weight() uint64 {
	return l.material.Weight()
}

func (l *Logs) AdminOnly() bool {
	return false
}
