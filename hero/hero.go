package hero

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

type HeroLike interface {
	world.Combat

	Race() Race
	Gender() Gender
	Occupation() Occupation

	Equip(*Equip)
	Unequip(EquipSlot)
	GetEquip(EquipSlot) *Equip
	world.InventoryLike
	canHoldItem(item world.Visible) bool
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

	return 1, map[string]interface{}{
		"name":       h.name.serialize(),
		"race":       uint64(h.race),
		"gender":     uint64(h.gender),
		"occupation": uint64(h.occupation),
		"skin":       h.skinTone,
		"birth":      h.birth,
		"death":      h.death,
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

		if dataMap["skin"] == nil {
			dataMap["skin"] = uint(rand.Intn(len(Race(dataMap["race"].(uint64)).SkinTones())))
		}

		if dataMap["equipped"] == nil {
			r := rand.New(rand.NewSource(rand.Int63()))
			h.equip(&Equip{
				slot:         SlotShirt,
				kind:         0,
				customColors: []string{randomColor(r)},
			})
			h.equip(&Equip{
				slot:         SlotPants,
				kind:         0,
				customColors: []string{randomColor(r)},
			})
			h.equip(&Equip{
				slot: SlotFeet,
				kind: 0,
			})

		}

		var err error
		dataMap["birth"], err = time.Parse(time.RFC3339Nano, dataMap["birth"].(string))
		if err != nil {
			panic(err)
		}
		dataMap["death"], err = time.Parse(time.RFC3339Nano, dataMap["death"].(string))
		if err != nil {
			panic(err)
		}
		fallthrough

	case 1:
		dataMap := data.(map[string]interface{})
		h.CombatObject = *attached[0].(*world.CombatObject)
		h.name.unserialize(dataMap["name"].(map[string]interface{}))
		h.race = Race(dataMap["race"].(uint64))
		h.gender = Gender(dataMap["gender"].(uint64))
		h.occupation = Occupation(dataMap["occupation"].(uint64))
		h.skinTone = dataMap["skin"].(uint)
		if h.equipped == nil {
			equipCount := dataMap["equipped"].(uint)
			for _, o := range attached[1 : equipCount+1] {
				h.equip(o.(*Equip))
			}
			itemCount := dataMap["items"].(uint)
			h.items = make([]world.Visible, itemCount)
			for i, o := range attached[equipCount+1 : itemCount+equipCount+1] {
				h.items[i] = o.(world.Visible)
			}
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

func (h *Hero) Examine() (string, [][][2]string) {
	_, info := h.CombatObject.Examine()

	info = append(info, [][2]string{
		{h.Race().Name(), "#4fc"},
		{" race", "#ccc"},
	}, [][2]string{
		{h.Gender().Name(), "#4fc"},
		{" gender", "#ccc"},
	})

	return "a hero.", info
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
	if h.animationTicks > 0 {
		h.animationTicks--
		if h.animationTicks == 0 {
			h.animation = ""
			if t := h.Position(); t != nil {
				h.mtx.Unlock()
				t.Zone().Update(t, h.Outer())
				return
			}
		}
	}
	h.mtx.Unlock()
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

func (h *Hero) MaxHealth() *big.Int {
	var max big.Int

	h.mtx.Lock()
	for _, e := range h.equipped {
		max.Add(&max, e.Material().Health())
	}
	h.mtx.Unlock()

	max.Mul(&max, world.TuningHealthMultiplier)

	max.Add(&max, h.Race().BaseHealth())

	return &max
}

func (h *Hero) MaxQuality() *big.Int {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	max := &big.Int{}

	for _, e := range h.equipped {
		q := e.Material().Quality()
		if q.Cmp(max) > 0 {
			max = q
		}
	}

	return max
}

func (h *Hero) MeleeDamage() *big.Int {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	var damage big.Int

	for _, e := range h.equipped {
		damage.Add(&damage, e.Material().MeleeDamage())
	}

	return &damage
}

func (h *Hero) MeleeArmor() *big.Int {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	var armor big.Int

	for _, e := range h.equipped {
		armor.Add(&armor, e.Material().MeleeArmor())
	}

	return &armor
}

func (h *Hero) CritChance() *big.Int {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	var crit big.Int

	for _, e := range h.equipped {
		crit.Add(&crit, e.Material().CritChance())
	}

	return &crit
}

func (h *Hero) Resistance() *big.Int {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	var resist big.Int

	for _, e := range h.equipped {
		resist.Add(&resist, e.Material().Resistance())
	}

	return &resist
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

func (h *Hero) Equip(e *Equip) {
	if pos := h.Position(); pos != nil {
		defer pos.Zone().Update(pos, h.Outer())
	}

	h.mtx.Lock()
	defer h.mtx.Unlock()

	h.equip(e)
}

func (h *Hero) equip(e *Equip) {
	if h.equipped == nil {
		h.equipped = make(map[EquipSlot]*Equip)
	}
	if old, ok := h.equipped[e.slot]; ok {
		old.wearer = nil
		h.giveItem(old)
	}
	e.wearer = h.Outer().(HeroLike)
	h.equipped[e.slot] = e
}

func (h *Hero) GetEquip(slot EquipSlot) *Equip {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	return h.equipped[slot]
}

func (h *Hero) Unequip(slot EquipSlot) {
	if pos := h.Position(); pos != nil {
		defer pos.Zone().Update(pos, h.Outer())
	}

	h.mtx.Lock()
	defer h.mtx.Unlock()

	h.unequip(slot)
}

func (h *Hero) unequip(slot EquipSlot) {
	if e, ok := h.equipped[slot]; ok {
		e.wearer = nil
		h.giveItem(e)
		delete(h.equipped, slot)
	}
}

func (h *Hero) notifyInventoryChanged() {
	// do nothing
}

func (h *Hero) GiveItem(item world.Visible) bool {
	if i, ok := item.(world.Item); !ok || i.AdminOnly() {
		if a, ok := h.Outer().(world.AdminLike); !ok || !a.IsAdmin() {
			if m, ok := h.Outer().(world.SendMessageLike); ok {
				m.SendMessage("the " + item.Name() + " deems you unworthy and decides to stay where it is.")
			}
			return false
		}
	}

	h.mtx.Lock()
	if h.Outer().(HeroLike).canHoldItem(item) {
		h.giveItem(item)
	}
	h.mtx.Unlock()

	return true
}

func (h *Hero) canHoldItem(item world.Visible) bool {
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
