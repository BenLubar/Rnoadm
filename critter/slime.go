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
	world.Register("slime", world.Combat((*Slime)(nil)))

	world.RegisterSpawnFunc(func(s string) world.Visible {
		for o := Slimupation(0); o < slimupationCount; o++ {
			if s == "slime "+o.Name() {
				const hex = "0123456789abcdef"
				color := string([]byte{'#', hex[rand.Intn(4)], hex[rand.Intn(4)+11], hex[rand.Intn(4)]})
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
	return "slime " + s.Slimupation().Name()
}

func (s *Slime) Examine() (string, [][][2]string) {
	_, info := s.CombatObject.Examine()

	return "a slime " + s.Slimupation().Flavor(), info
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
