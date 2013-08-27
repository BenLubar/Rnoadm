package critter

import (
	"github.com/BenLubar/Rnoadm/hero"
	"github.com/BenLubar/Rnoadm/world"
	"math/rand"
)

func noSchedule(*Critter, *world.Tile) world.Schedule {
	return nil
}

func followHeroSchedule(radius int) func(*Critter, *world.Tile) world.Schedule {
	if radius <= 0 {
		return noSchedule
	}

	return func(c *Critter, pos *world.Tile) world.Schedule {
		lastX, lastY := pos.Position()
		for _, i := range rand.Perm(radius * radius * 4) {
			x8, y8 := pos.Position()
			x, y := int(x8), int(y8)
			x += i%(radius*2) - radius
			y += i/(radius*2) - radius
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
				return world.NewWalkSchedule(x8, y8, true, 2)
			}
			lastX, lastY = x8, y8
		}
		return &world.ScheduleSchedule{
			Schedules: []world.Schedule{
				world.NewWalkSchedule(lastX, lastY, false, 4),
				&world.DelaySchedule{Delay: 5},
			},
		}
	}
}

func wanderSchedule(distance int) func(*Critter, *world.Tile) world.Schedule {
	if distance <= 0 {
		return noSchedule
	}

	return func(c *Critter, pos *world.Tile) world.Schedule {
		x8, y8 := pos.Position()
		x, y := int(x8), int(y8)
		x += rand.Intn(distance*2) - distance
		y += rand.Intn(distance*2) - distance
		if x < 0 || x > 255 || y < 0 || y > 255 {
			return nil
		}
		return &world.ScheduleSchedule{
			Schedules: []world.Schedule{
				&world.DelaySchedule{Delay: 5},
				world.NewWalkSchedule(uint8(x), uint8(y), false, 6),
			},
		}
	}
}
