package hero

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"math/rand"
	"sync"
	"time"
)

type HeroLike interface {
	world.Combat

	Race() Race
	Gender() Gender
	Occupation() Occupation
}

type Hero struct {
	world.CombatObject

	name HeroName

	race       Race
	gender     Gender
	occupation Occupation
	skinTone   uint

	birth time.Time
	death time.Time

	mtx sync.Mutex
}

func init() {
	world.Register("hero", HeroLike((*Hero)(nil)))
}

func (h *Hero) Save() (uint, interface{}, []world.ObjectLike) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	return 0, map[string]interface{}{
		"name":       h.name.serialize(),
		"race":       uint64(h.race),
		"gender":     uint64(h.gender),
		"occupation": uint64(h.occupation),
		"skin":       h.skinTone,
		"birth":      h.birth.Format(time.RFC3339Nano),
		"death":      h.death.Format(time.RFC3339Nano),
	}, []world.ObjectLike{&h.CombatObject}
}

func (h *Hero) Load(version uint, data interface{}, attached []world.ObjectLike) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		h.CombatObject = *attached[0].(*world.CombatObject)
		h.name.unserialize(dataMap["name"].(map[string]interface{}))
		h.race = Race(dataMap["race"].(uint64))
		h.gender = Gender(dataMap["gender"].(uint64))
		h.occupation = Occupation(dataMap["occupation"].(uint64))
		var ok bool
		h.skinTone, ok = dataMap["skin"].(uint)
		if !ok {
			h.skinTone = uint(rand.Intn(len(h.race.SkinTones())))
		}
		var err error
		h.birth, err = time.Parse(time.RFC3339Nano, dataMap["birth"].(string))
		if err != nil {
			panic(err)
		}
		h.death, err = time.Parse(time.RFC3339Nano, dataMap["death"].(string))
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (h *Hero) Name() string {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	return h.name.Name()
}

func (h *Hero) Race() Race {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	return h.race
}

func (h *Hero) Gender() Gender {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	return h.gender
}

func (h *Hero) Occupation() Occupation {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	return h.occupation
}

func (h *Hero) Sprite() string {
	return h.Race().Sprite()
}

func (h *Hero) Colors() []string {
	return []string{h.Race().SkinTones()[h.skinTone]}
}

func (h *Hero) MaxHealth() uint64 {
	return h.Race().BaseHealth()
}
