package critter

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/hero"
	"github.com/BenLubar/Rnoadm/world"
	"math/rand"
	"sync"
)

type Slime struct {
	world.CombatObject

	slimupation Slimupation
	gelTone     string

	facing         uint   // not saved
	animation      string // not saved
	animationTicks uint   // not saved

	mtx sync.Mutex
}

func init() {
	world.Register("slime", world.Combat((*Slime)(nil)))

	world.RegisterSpawnFunc(func(s string) world.Visible {
		for o := Slimupation(0); o < slimupationCount; o++ {
			if s == o.Name() {
				const hex = "0123456789abcdef"
				color := string([]byte{'#', hex[rand.Intn(6)], hex[rand.Intn(6)+9], hex[rand.Intn(6)]})
				return &Slime{
					slimupation: o,
					gelTone:     color,
				}
			}
		}
		return nil
	})
}

func (s *Slime) Save() (uint, interface{}, []world.ObjectLike) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	attached := []world.ObjectLike{&s.CombatObject}

	return 0, map[string]interface{}{
		"s": uint64(s.slimupation),
		"t": s.gelTone,
	}, attached
}

func (s *Slime) Load(version uint, data interface{}, attached []world.ObjectLike) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		s.CombatObject = *attached[0].(*world.CombatObject)
		s.slimupation = Slimupation(dataMap["s"].(uint64))
		s.gelTone = dataMap["t"].(string)

	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}

}

func (s *Slime) Name() string {
	return s.Slimupation().Name()
}

func (s *Slime) Examine() (string, [][][2]string) {
	_, info := s.CombatObject.Examine()

	return s.Slimupation().Flavor(), info
}

func (s *Slime) Slimupation() Slimupation {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.slimupation
}

func (s *Slime) Sprite() string {
	return "critter_slime"
}

func (s *Slime) Colors() []string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return []string{s.gelTone}
}

func (s *Slime) Think() {
	s.CombatObject.Think()

	if pos := s.Position(); pos != nil && !s.HasSchedule() {
		radius := s.Slimupation().Radius()
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
				s.SetSchedule(world.NewWalkSchedule(x8, y8, true))
				break
			}
			lastX, lastY = x8, y8
		}
		if !s.HasSchedule() {
			s.SetSchedule(&world.ScheduleSchedule{
				Schedules: []world.Schedule{
					world.NewWalkSchedule(lastX, lastY, false),
					&world.DelaySchedule{Delay: 5},
				},
			})
		}
	}

	s.mtx.Lock()
	if s.animationTicks > 0 {
		s.animationTicks--
		if s.animationTicks == 0 {
			s.animation = ""
			if t := s.Position(); t != nil {
				s.mtx.Unlock()
				t.Zone().Update(t, s.Outer())
				return
			}
		}
	}
	s.mtx.Unlock()
}

func (s *Slime) AnimationType() string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.animation
}

func (s *Slime) SpritePos() (uint, uint) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.facing, 0
}

func (s *Slime) MaxHealth() uint64 {
	return 50
}

func (s *Slime) NotifyPosition(old, new *world.Tile) {
	if old == nil || new == nil {
		s.mtx.Lock()
		s.animationTicks = 0
		s.animation = ""
		s.facing = 0
		s.mtx.Unlock()
		return
	}
	ox, oy := old.Position()
	nx, ny := new.Position()

	s.mtx.Lock()
	switch {
	case ox-1 == nx && oy == ny:
		s.animationTicks = 3
		s.animation = "wa" // walk (alternating)
		s.facing = 6
	case ox+1 == nx && oy == ny:
		s.animationTicks = 3
		s.animation = "wa" // walk (alternating)
		s.facing = 9
	case ox == nx && oy-1 == ny:
		s.animationTicks = 3
		s.animation = "wa" // walk (alternating)
		s.facing = 3
	case ox == nx && oy+1 == ny:
		s.animationTicks = 3
		s.animation = "wa" // walk (alternating)
		s.facing = 0
	}
	s.mtx.Unlock()

	new.Zone().Update(new, s.Outer())
}
