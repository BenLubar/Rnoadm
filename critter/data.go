package critter

import (
	"github.com/BenLubar/Rnoadm/world"
	"math/rand"
)

const hex = "0123456789abcdef"

type CritterType uint64

const (
	SlimeMage CritterType = iota
	SlimeBrute
	StickySlime
	Cow

	critterTypeCount
)

var critterInfo = [critterTypeCount]struct {
	name      string
	examine   string
	maxHealth uint64
	sprite    string
	genColors func() []string
	schedule  func(*Critter, *world.Tile) world.Schedule
}{
	SlimeMage: {
		name:      "slime mage",
		examine:   "a slime with some magical abilities.",
		maxHealth: 100,
		sprite:    "critter_slime",
		genColors: func() []string {
			return []string{string([]byte{
				'#',
				hex[rand.Intn(6)],
				hex[rand.Intn(6)],
				hex[rand.Intn(6)+10],
			})}
		},
		schedule: followHeroSchedule(7),
	},
	SlimeBrute: {
		name:      "slime brute",
		examine:   "a slime that can bench press a whole hero.",
		maxHealth: 250,
		sprite:    "critter_slime",
		genColors: func() []string {
			return []string{string([]byte{
				'#',
				hex[rand.Intn(6)+10],
				hex[rand.Intn(6)],
				hex[rand.Intn(6)],
			})}
		},
		schedule: followHeroSchedule(7),
	},
	StickySlime: {
		name:      "sticky slime",
		examine:   "a slime that's a little stickier than most.",
		maxHealth: 50,
		sprite:    "critter_slime",
		genColors: func() []string {
			return []string{string([]byte{
				'#',
				hex[rand.Intn(6)],
				hex[rand.Intn(6)+10],
				hex[rand.Intn(6)],
			})}
		},
		schedule: followHeroSchedule(15),
	},
	Cow: {
		name:      "cow",
		examine:   "it's a cow, i guess.",
		maxHealth: 10,
		sprite:    "critter_slime",
		genColors: func() []string {
			return []string{string([]byte{
				'#',
				hex[rand.Intn(6)],
				hex[rand.Intn(6)+10],
				hex[rand.Intn(6)],
			})}
		},
		schedule: wanderSchedule(4),
	},
}

func (t CritterType) Name() string             { return critterInfo[t].name }
func (t CritterType) Examine() string          { return critterInfo[t].examine }
func (t CritterType) MaxHealth() uint64        { return critterInfo[t].maxHealth }
func (t CritterType) Sprite() string           { return critterInfo[t].sprite }
func (t CritterType) GenerateColors() []string { return critterInfo[t].genColors() }
func (t CritterType) Schedule(c *Critter, pos *world.Tile) world.Schedule {
	return critterInfo[t].schedule(c, pos)
}
