package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
)

type WoodType uint64

func (t WoodType) Name() string     { return woodTypes[t].name }
func (t WoodType) Strength() uint64 { return woodTypes[t].strength }

var woodTypes = []struct {
	name     string
	strength uint64
}{
	{
		name:     "wood0",
		strength: 5 << 0,
	},
	{
		name:     "wood0",
		strength: 5 << 0,
	},
	{
		name:     "wood1",
		strength: 5 << 1,
	},
	{
		name:     "wood2",
		strength: 5 << 2,
	},
	{
		name:     "wood3",
		strength: 5 << 3,
	},
	{
		name:     "wood4",
		strength: 5 << 4,
	},
	{
		name:     "wood5",
		strength: 5 << 5,
	},
	{
		name:     "wood6",
		strength: 5 << 6,
	},
	{
		name:     "wood7",
		strength: 5 << 7,
	},
	{
		name:     "wood8",
		strength: 5 << 8,
	},
	{
		name:     "wood9",
		strength: 5 << 9,
	},
	{
		name:     "wood10",
		strength: 5 << 10,
	},
	{
		name:     "wood11",
		strength: 5 << 11,
	},
	{
		name:     "wood12",
		strength: 5 << 12,
	},
	{
		name:     "wood13",
		strength: 5 << 13,
	},
	{
		name:     "wood14",
		strength: 5 << 14,
	},
	{
		name:     "wood15",
		strength: 5 << 15,
	},
}

func init() {
	world.RegisterSpawnFunc(func(s string) world.Visible {
		if len(s) > len(" tree") && s[len(s)-len(" tree"):] == " tree" {
			for i, t := range woodTypes {
				if len(s) == len(t.name)+len(" tree") && s[:len(t.name)] == t.name {
					return world.InitObject(&Tree{kind: WoodType(i)}).(world.Visible)
				}
			}
		} else if len(s) > len(" logs") && s[len(s)-len(" logs"):] == " logs" {
			for i, t := range woodTypes {
				if len(s) == len(t.name)+len(" logs") && s[:len(t.name)] == t.name {
					return world.InitObject(&Logs{kind: WoodType(i)}).(world.Visible)
				}
			}
		}
		return nil
	})
}

type Tree struct {
	Node

	kind WoodType
}

func init() {
	world.Register("tree", NodeLike((*Tree)(nil)))
}

func (t *Tree) Save() (uint, interface{}, []world.ObjectLike) {
	return 0, uint64(t.kind), []world.ObjectLike{&t.Node}
}

func (t *Tree) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		t.Node = *attached[0].(*Node)
		t.kind = WoodType(data.(uint64))
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (t *Tree) Strength() uint64 {
	return t.kind.Strength()
}

func (t *Tree) Name() string {
	return t.kind.Name() + " tree"
}

type Logs struct {
	world.VisibleObject

	kind WoodType
}

func init() {
	world.Register("logs", world.Visible((*Logs)(nil)))
}

func (l *Logs) Save() (uint, interface{}, []world.ObjectLike) {
	return 0, uint64(l.kind), nil
}

func (l *Logs) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		l.kind = WoodType(data.(uint64))
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (l *Logs) Name() string {
	return l.kind.Name() + " logs"
}
