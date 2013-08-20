package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
)

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

func (t *Tree) Sprite() string {
	return "tree"
}

func (t *Tree) SpriteSize() (uint, uint) {
	return 96, 128
}

func (t *Tree) Colors() []string {
	switch t.kind.LeafType() {
	default:
		fallthrough
	case 0: // no leaves
		return []string{t.kind.BarkColor()}
	case 1: // deciduous
		return []string{t.kind.BarkColor(), t.kind.LeafColor()}
	case 2: // coniferous
		return []string{t.kind.BarkColor(), "", t.kind.LeafColor()}
	}
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
