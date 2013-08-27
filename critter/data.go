package critter

import (
	"github.com/BenLubar/Rnoadm/hero"
	"github.com/BenLubar/Rnoadm/world"
	"math/rand"
)

const hex = "0123456789abcdef"

type CritterType uint64

const (
	SlimeMage CritterType = iota
	SlimeBrute
	StickySlime

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
}

func (t CritterType) Name() string             { return critterInfo[t].name }
func (t CritterType) Examine() string          { return critterInfo[t].examine }
func (t CritterType) MaxHealth() uint64        { return critterInfo[t].maxHealth }
func (t CritterType) Sprite() string           { return critterInfo[t].sprite }
func (t CritterType) GenerateColors() []string { return critterInfo[t].genColors() }
func (t CritterType) Schedule(c *Critter, pos *world.Tile) world.Schedule {
	return critterInfo[t].schedule(c, pos)
}

func followHeroSchedule(radius int) func(*Critter, *world.Tile) world.Schedule {
	return func(c *Critter, pos *world.Tile) world.Schedule {
		lastX, lastY := pos.Position()
		for i := 0; i < 1000; i++ {
			x8, y8 := pos.Position()
			x, y := int(x8), int(y8)
			x += rand.Intn(radius*2) - radius
			y += rand.Intn(radius*2) - radius
			if x < 0 || x > 255 || y < 0 || y > 255 {
				continue
			}
			x8, y8 = uint8(x), uint8(y)
			t := pos.Zone().Tile(x8, y8)
			foundHero := false
			for _, o := range t.Objects() {
				if _, ok := o.(hero.HeroLike); ok {
					foundHero = true
					break
				}
			}
			if foundHero {
				return world.NewWalkSchedule(x8, y8, true)
			}
			lastX, lastY = x8, y8
		}
		return &world.ScheduleSchedule{
			Schedules: []world.Schedule{
				world.NewWalkSchedule(lastX, lastY, false),
				&world.DelaySchedule{Delay: 5},
			},
		}
	}
}
