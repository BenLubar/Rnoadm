package critter

import (
	"fmt"
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
	world.Register("hero", HeroLike((*Hero)(nil)))

	world.RegisterSpawnFunc(func(s string) world.Visible {
		for i := range raceInfo {
			r := Race(i)
			if r.Name() == s {
				return GenerateHeroRace(rand.New(rand.NewSource(rand.Int63())), r)
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
		"slimupation": uint64(s.slimupation),
		"tone":        s.gelTone,
	}, attached
}

func (s *Slime) Load(version uint, data interface{}, attached []world.ObjectLike) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		s.CombatObject = *attached[0].(*world.CombatObject)
		s.slimupation = Occupation(dataMap["slimupation"].(uint64))
		var ok bool
		s.gelTone, ok = dataMap["tone"].(string)
		if !ok {
			s.gelTone = "#0f0"
		}

	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}

}

func (s *Slime) Name() string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return "Slime " + s.slimupation.Name()
}

func (s *Slime) Examine() string {
	return "a slime " + s.slimupation.ExamineFlavor()
}

func (h *Hero) Slimupation() Slimupation {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.slimupation
}

func (s *Slime) Sprite() string {
	return "critter_slime"
}

func (h *Hero) Think() {
	s.CombatObject.Think()

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

func (h *Hero) MaxHealth() uint64 {
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
