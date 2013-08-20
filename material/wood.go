package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
)

type WoodType uint64

func (t WoodType) Name() string      { return woodTypes[t].name }
func (t WoodType) BarkColor() string { return woodTypes[t].barkColor }
func (t WoodType) LeafColor() string { return woodTypes[t].leafColor }
func (t WoodType) LeafType() uint    { return woodTypes[t].leafType }
func (t WoodType) Strength() uint64  { return woodTypes[t].strength }

var woodTypes = []struct {
	name      string
	barkColor string
	leafColor string
	leafType  uint
	strength  uint64
}{
	{
		name:      "wood0",
		strength:  5 << 0,
		barkColor: "#000",
	},
	{
		name:      "wood1",
		strength:  5 << 1,
		barkColor: "#111",
	},
	{
		name:      "wood2",
		strength:  5 << 2,
		barkColor: "#222",
	},
	{
		name:      "wood3",
		strength:  5 << 3,
		barkColor: "#333",
	},
	{
		name:      "wood4",
		strength:  5 << 4,
		barkColor: "#444",
	},
	{
		name:      "wood5",
		strength:  5 << 5,
		barkColor: "#555",
	},
	{
		name:      "wood6",
		strength:  5 << 6,
		barkColor: "#666",
	},
	{
		name:      "wood7",
		strength:  5 << 7,
		barkColor: "#777",
	},
	{
		name:      "wood8",
		strength:  5 << 8,
		barkColor: "#888",
	},
	{
		name:      "wood9",
		strength:  5 << 9,
		barkColor: "#999",
	},
	{
		name:      "wood10",
		strength:  5 << 10,
		barkColor: "#aaa",
	},
	{
		name:      "wood11",
		strength:  5 << 11,
		barkColor: "#bbb",
	},
	{
		name:      "wood12",
		strength:  5 << 12,
		barkColor: "#ccc",
	},
	{
		name:      "wood13",
		strength:  5 << 13,
		barkColor: "#ddd",
	},
	{
		name:      "wood14",
		strength:  5 << 14,
		barkColor: "#eee",
	},
	{
		name:      "wood15",
		strength:  5 << 15,
		barkColor: "#fff",
	},
	{
		name:      "wood16",
		strength:  1, // TODO
		barkColor: "#d2b48c",
		leafColor: "#4b5a3f",
		leafType:  1,
	},
	{
		name:      "wood17",
		strength:  1, // TODO
		barkColor: "#b5aa8b",
		leafColor: "#cf5123",
		leafType:  2,
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
