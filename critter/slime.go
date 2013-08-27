package critter

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"sync"
)

// Slime is a legacy object. When a Slime is in the word, it replaces itself
// with an equivelent Critter.
type Slime struct {
	world.VisibleObject

	slimupation CritterType
	gelTone     string

	mtx sync.Mutex
}

func init() {
	world.Register("slime", world.Visible((*Slime)(nil)))
}

func (s *Slime) Save() (uint, interface{}, []world.ObjectLike) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return 0, map[string]interface{}{
		"s": uint64(s.slimupation),
		"t": s.gelTone,
	}, nil
}

func (s *Slime) Load(version uint, data interface{}, attached []world.ObjectLike) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		s.slimupation = CritterType(dataMap["s"].(uint64))
		s.gelTone = dataMap["t"].(string)

	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}

}

func (s *Slime) Name() string {
	return s.Type().Name()
}

func (s *Slime) Examine() (string, [][][2]string) {
	_, info := s.VisibleObject.Examine()

	return s.Type().Examine(), info
}

func (s *Slime) Type() CritterType {
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
	if pos := s.Position(); pos != nil {
		if pos.Remove(s) {
			pos.Add(&Critter{
				kind:   s.Type(),
				colors: s.Colors(),
			})
		}
	}
}
