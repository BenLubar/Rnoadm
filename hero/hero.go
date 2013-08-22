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

	world.InventoryLike
	notifyInventoryChanged()
}

type Hero struct {
	world.CombatObject

	name HeroName

	race       Race
	gender     Gender
	occupation Occupation
	skinTone   uint

	equipped map[EquipSlot]*Equip
	items    []world.Visible

	birth time.Time
	death time.Time

	facing         uint   // not saved
	animation      string // not saved
	animationTicks uint   // not saved

	mtx sync.Mutex
}

var _ world.InventoryLike = (*Hero)(nil)

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

func (h *Hero) Save() (uint, interface{}, []world.ObjectLike) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	attached := []world.ObjectLike{&h.CombatObject}
	for _, e := range h.equipped {
		attached = append(attached, e)
	}
	for _, i := range h.items {
		attached = append(attached, i)
	}

	return 0, map[string]interface{}{
		"name":       h.name.serialize(),
		"race":       uint64(h.race),
		"gender":     uint64(h.gender),
		"occupation": uint64(h.occupation),
		"skin":       h.skinTone,
		"birth":      h.birth.Format(time.RFC3339Nano),
		"death":      h.death.Format(time.RFC3339Nano),
		"equipped":   uint(len(h.equipped)),
		"items":      uint(len(h.items)),
	}, attached
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

		equipCount, ok := dataMap["equipped"].(uint)
		if ok {
			h.equipped = make(map[EquipSlot]*Equip, equipCount)
			for _, o := range attached[1 : equipCount+1] {
				e := o.(*Equip)
				h.equipped[e.slot] = e
			}
			itemCount := dataMap["items"].(uint)
			h.items = make([]world.Visible, itemCount)
			for i, o := range attached[equipCount+1 : itemCount+equipCount+1] {
				h.items[i] = o.(world.Visible)
			}
		} else {
			r := rand.New(rand.NewSource(rand.Int63()))
			h.equipped = map[EquipSlot]*Equip{
				SlotShirt: {
					slot:         SlotShirt,
					kind:         0,
					customColors: []string{randomColor(r)},
				},
				SlotPants: {
					slot:         SlotPants,
					kind:         0,
					customColors: []string{randomColor(r)},
				},
				SlotFeet: {
					slot: SlotFeet,
					kind: 0,
				},
			}
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}

	for _, e := range h.equipped {
		e.wearer = h
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
	h.mtx.Lock()
	tone := h.skinTone
	h.mtx.Unlock()

	return []string{h.Race().SkinTones()[tone]}
}

func (h *Hero) Attached() []world.Visible {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	attached := make([]world.Visible, 0, len(h.equipped))
	for _, s := range equipFacing[h.facing/3] {
		if e := h.equipped[s]; e != nil {
			attached = append(attached, e)
		}
	}
	return attached
}

func (h *Hero) Think() {
	h.CombatObject.Think()

	h.mtx.Lock()
	defer h.mtx.Unlock()
	if h.animationTicks > 0 {
		h.animationTicks--
		if h.animationTicks == 0 {
			h.animation = ""
			if t := h.Position(); t != nil {
				go t.Zone().Update(t, h.Outer())
			}
		}
	}
}

func (h *Hero) AnimationType() string {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	return h.animation
}

func (h *Hero) SpritePos() (uint, uint) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	return h.facing, 0
}

func (h *Hero) SpriteSize() (uint, uint) {
	return h.Race().SpriteSize()
}

func (h *Hero) MaxHealth() uint64 {
	return h.Race().BaseHealth()
}

func (h *Hero) NotifyPosition(old, new *world.Tile) {
	if old == nil || new == nil {
		h.mtx.Lock()
		h.animationTicks = 0
		h.animation = ""
		h.facing = 0
		h.mtx.Unlock()
		return
	}
	ox, oy := old.Position()
	nx, ny := new.Position()

	h.mtx.Lock()
	switch {
	case ox-1 == nx && oy == ny:
		h.animationTicks = 3
		h.animation = "wa" // walk (alternating)
		h.facing = 6
	case ox+1 == nx && oy == ny:
		h.animationTicks = 3
		h.animation = "wa" // walk (alternating)
		h.facing = 9
	case ox == nx && oy-1 == ny:
		h.animationTicks = 3
		h.animation = "wa" // walk (alternating)
		h.facing = 3
	case ox == nx && oy+1 == ny:
		h.animationTicks = 3
		h.animation = "wa" // walk (alternating)
		h.facing = 0
	}
	h.mtx.Unlock()

	new.Zone().Update(new, h.Outer())
}

func (h *Hero) Inventory() []world.Visible {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	inventory := make([]world.Visible, len(h.items))
	copy(inventory, h.items)
	return inventory
}

func (h *Hero) notifyInventoryChanged() {
	// do nothing
}

func (h *Hero) GiveItem(item world.Visible) bool {
	if i, ok := item.(world.Item); !ok || i.AdminOnly() {
		return false
	}

	h.mtx.Lock()
	if true {
		h.giveItem(item)
	}
	h.mtx.Unlock()

	return true
}

func (h *Hero) giveItem(item world.Visible) {
	h.items = append(h.items, item)
	h.Outer().(HeroLike).notifyInventoryChanged()
}

func (h *Hero) RemoveItem(item world.Visible) bool {
	found := false
	h.mtx.Lock()
	for i, o := range h.items {
		if o == item {
			h.items = append(h.items[:i], h.items[i+1:]...)
			found = true
			h.Outer().(HeroLike).notifyInventoryChanged()
			break
		}
	}
	h.mtx.Unlock()

	return found
}
