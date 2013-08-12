package main

import (
	"math/rand"
)

type CritterType uint16

const (
	Slime CritterType = iota

	critterTypeCount
)

var critterInfo = [critterTypeCount]struct {
	Name    string
	Examine string

	Sprite string
	Colors []Color

	Health uint64
}{
	Slime: {
		Name:    "slime",
		Examine: "slimy.",

		Sprite: "critter_slime",
		Colors: []Color{"#0f0"},

		Health: 5000,
	},
}

type Critter struct {
	networkID
	Type CritterType

	Delay uint

	combat *Hero
}

func (c *Critter) Name() string {
	return critterInfo[c.Type].Name
}

func (c *Critter) Examine() string {
	return critterInfo[c.Type].Examine
}

func (c *Critter) Blocking() bool {
	return false
}

func (c *Critter) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:    c.Name(),
		Sprite:  critterInfo[c.Type].Sprite,
		Colors:  critterInfo[c.Type].Colors,
		Moves:   true,
		Options: []string{"attack"},
	}
}

func (c *Critter) Think(z *Zone, x, y uint8) {
	if c.Delay > 0 {
		c.Delay--
		return
	}

	if c.combat == nil {
		tx, ty := x+uint8(rand.Intn(3)-1), y+uint8(rand.Intn(3)-1)
		if tx != x && ty != y {
			return
		}
		z.Lock()
		if t1, t2 := z.Tile(x, y), z.Tile(tx, ty); t2 != nil && !t2.Blocked() && t1.Remove(c) {
			t2.Add(c)
			SendZoneTileChange(z.X, z.Y, TileChange{
				ID: c.NetworkID(),
				X:  tx,
				Y:  ty,
			})
		}
		z.Unlock()
		c.Delay = 8
		return
	}
}
