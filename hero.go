package main

import (
	"compress/gzip"
	"encoding/binary"
	"encoding/gob"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Gender uint8

const (
	Male Gender = iota
	Female

	genderCount
)

var genderInfo = [genderCount]struct {
	Name             string
	PronounSubject   string
	PronounObject    string
	PronounPosessive string
	NounOffspring    string
}{
	Male: {
		Name:             "male",
		PronounSubject:   "he",
		PronounObject:    "him",
		PronounPosessive: "his",
		NounOffspring:    "son",
	},
	Female: {
		Name:             "female",
		PronounSubject:   "she",
		PronounObject:    "her",
		PronounPosessive: "her",
		NounOffspring:    "daughter",
	},
}

type Race uint16

const (
	Human Race = iota

	raceCount
)

var raceInfo = [raceCount]struct {
	Name         string
	Genders      []Gender
	Occupations  []Occupation
	SkinTones    []Color
	BaseHealth   uint64
	BaseArmor    uint64
	BaseDamage   uint64
	BaseAccuracy uint64
}{
	Human: {
		Name:         "human",
		Genders:      []Gender{Male, Female},
		Occupations:  []Occupation{Citizen},
		SkinTones:    []Color{"#ffe3cc", "#ffdbbd", "#ffcda3", "#f7e9dc", "#edd0b7", "#e8d1be", "#e5c298", "#e3c3a8", "#c9a281", "#c2a38a", "#ba9c82", "#ad8f76", "#a17a5a", "#876d58", "#6e5948", "#635041", "#4f3f33"},
		BaseHealth:   10000,
		BaseArmor:    100,
		BaseDamage:   100,
		BaseAccuracy: 100,
	},
}

type Occupation uint16

const (
	Adventurer Occupation = iota
	Citizen

	occupationCount
)

var occupationInfo = [occupationCount]struct {
	Name            string
	StartingGear    map[CosmeticType]*Cosmetic
	StartingPickaxe *Pickaxe
	StartingHatchet *Hatchet
}{
	Adventurer: {
		Name: "adventurer",
		StartingGear: map[CosmeticType]*Cosmetic{
			Shirt: &Cosmetic{
				Type: Shirt,
				ID:   0,
			},
			Pants: &Cosmetic{
				Type: Pants,
				ID:   0,
			},
			Shoes: &Cosmetic{
				Type: Shoes,
				ID:   0,
			},
		},
		StartingPickaxe: &Pickaxe{
			Head:   RustyMetal,
			Handle: RottingWood,
		},
		StartingHatchet: &Hatchet{
			Head:   RustyMetal,
			Handle: RottingWood,
		},
	},
	Citizen: {
		Name: "citizen",
		StartingGear: map[CosmeticType]*Cosmetic{
			Shirt: &Cosmetic{
				Type: Shirt,
				ID:   0,
			},
			Pants: &Cosmetic{
				Type: Pants,
				ID:   0,
			},
			Shoes: &Cosmetic{
				Type: Shoes,
				ID:   0,
			},
		},
	},
}

var userIDLock sync.Mutex

func newUserID() uint64 {
	userIDLock.Lock()
	defer userIDLock.Unlock()

	var userID uint64

	fn := filepath.Join(seedFilename(), "mUSERID.gz")
	f, err := os.Open(fn)
	if err == nil {
		g, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			panic(err)
		}
		err = gob.NewDecoder(g).Decode(&userID)
		g.Close()
		f.Close()
		if err != nil {
			panic(err)
		}
	}

	userID++

	f, err = os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	g, _ := gzip.NewWriterLevel(f, gzip.BestCompression)
	defer g.Close()
	err = gob.NewEncoder(g).Encode(&userID)
	if err != nil {
		panic(err)
	}

	return userID
}

type AncestryInfo struct {
	Hero  *Hero
	Death time.Time
}

type Player struct {
	networkID
	ID uint64
	*Hero
	ZoneX, ZoneY int64
	TileX, TileY uint8

	Seed RandomSource

	LastLoginAddr string
	LastLogin     time.Time
	Admin         bool
	Examine_      string

	messages  chan<- _Message
	kick      chan<- string
	hud       chan packetSetHUD
	inventory chan<- packetInventory
	updates   <-chan TileChange

	CharacterCreation *Hero

	Ancestry []AncestryInfo

	lock sync.Mutex

	zone *Zone
}

func (p *Player) Lock() {
	p.lock.Lock()
	if p.Hero != nil {
		p.Hero.Lock()
	}
}

func (p *Player) Unlock() {
	if p.Hero != nil {
		p.Hero.Unlock()
	}
	p.lock.Unlock()
}

func (p *Player) SendMessage(message string) {
	p.SendMessageColor(message, "")
}

func (p *Player) SendMessageColor(message string, color Color) {
	select {
	case p.messages <- _Message{
		Text:  message,
		Color: color,
	}:
	case <-time.After(time.Minute):
	}
}

func (p *Player) Kick(message string) {
	select {
	case p.kick <- message:
	default: // player was already kicked
	}
}

func (p *Player) SetHUD(name string, data map[string]interface{}) {
	for {
		select {
		case p.hud <- packetSetHUD{
			SetHUD: _SetHUD{
				Name: name,
				Data: data,
			},
		}:
			return
		case <-p.hud:
		}
	}
}

func (p *Player) CharacterCreationCommand(command string) {
	p.Lock()
	defer p.Unlock()

	if p.Hero != nil && p.Hero.Health() > 0 {
		return
	}

	if p.Hero != nil {
		zone := p.zone
		tx, ty := p.TileX, p.TileY
		p.Unlock()
		zone.Lock()
		if zone.Tile(tx, ty).Remove(p) {
			SendZoneTileChange(zone.X, zone.Y, TileChange{
				ID:      p.NetworkID(),
				Removed: true,
			})
		}
		zone.Unlock()

		if p.Admin {
			zone.Lock()
			zone.Tile(tx, ty).Add(p)
			zone.Unlock()
			SendZoneTileChange(zone.X, zone.Y, TileChange{
				ID:  p.NetworkID(),
				X:   tx,
				Y:   ty,
				Obj: p.Serialize(),
			})
			p.Lock()
			p.Damage = 0
			p.SetHUD("", nil)
			return
		}

		p.Lock()
		if command == "open" {
			p.ZoneX, p.ZoneY = 0, 0
			p.TileX, p.TileY = 127, 127
			ReleaseZone(p.zone, p)
			p.zone, p.updates = GrabZone(0, 0, p)
			p.Ancestry = append(p.Ancestry, AncestryInfo{
				Hero:  p.Hero,
				Death: time.Now().UTC(),
			})
			p.Unlock()
			p.Save()
			p.Lock()
		} else {
			const timeFormat = "2 January 2006"
			cause := "cause of death unknown"
			if p.Hero.DeathCause != "" {
				cause = p.Hero.DeathCause
			}
			p.SetHUD("death", map[string]interface{}{
				"N": p.Name(),
				"D": p.Hero.Birth.Format(timeFormat) + " - " + time.Now().UTC().Format(timeFormat),
				"C": cause,
			})
			return
		}
	}

	r := rand.New(&p.Seed)

	if p.CharacterCreation == nil {
		if p.Hero != nil {
			p.CharacterCreation = GenerateHeroOccupation(p.Hero.Race, Adventurer, r)
			p.CharacterCreation.Last1 = p.Hero.Last1
			p.CharacterCreation.Last1T = p.Hero.Last1T
			p.CharacterCreation.Last2 = p.Hero.Last2
			p.CharacterCreation.Last2T = p.Hero.Last2T
			p.CharacterCreation.Last3 = p.Hero.Last3
			p.CharacterCreation.Last3T = p.Hero.Last3T
			p.CharacterCreation.SkinToneIndex = p.Hero.SkinToneIndex
			p.Hero = nil
		} else {
			p.CharacterCreation = GenerateHeroOccupation(Human, Adventurer, r)
		}
	}

	race := raceInfo[p.CharacterCreation.Race]
	switch command {
	case "race":
		p.CharacterCreation.Race = Human // TODO: more races
		race = raceInfo[p.CharacterCreation.Race]
		found := false
		for _, g := range race.Genders {
			if p.CharacterCreation.Gender == g {
				found = true
				break
			}
		}
		if !found {
			p.CharacterCreation.Gender = race.Genders[r.Intn(len(race.Genders))]
		}
		if p.CharacterCreation.SkinToneIndex >= uint8(len(race.SkinTones)) {
			p.CharacterCreation.SkinToneIndex = uint8(r.Intn(len(race.SkinTones)))
		}
	case "gender":
		found := false
		changed := false
		for _, g := range race.Genders {
			if found {
				p.CharacterCreation.Gender = g
				changed = true
				break
			}
			if p.CharacterCreation.Gender == g {
				found = true
			}
		}
		if !changed {
			p.CharacterCreation.Gender = race.Genders[0]
		}
		fallthrough
	case "name":
		switch p.CharacterCreation.Race {
		case Human:
			p.CharacterCreation.HeroName = GenerateHumanName(r, p.CharacterCreation.Gender)
		}
	case "skin":
		p.CharacterCreation.SkinToneIndex = (p.CharacterCreation.SkinToneIndex + 1) % uint8(len(race.SkinTones))
	case "shirt":
		const pastels = "abcde"
		const earthy = "34567"
		palette := pastels
		if r.Intn(2) == 0 {
			palette = earthy
		}
		p.CharacterCreation.Worn[Shirt].Custom = []Color{Color([]byte{
			'#',
			palette[r.Intn(len(palette))],
			palette[r.Intn(len(palette))],
			palette[r.Intn(len(palette))],
		})}
	case "pants":
		const pastels = "abcde"
		const earthy = "34567"
		palette := pastels
		if r.Intn(2) == 0 {
			palette = earthy
		}
		p.CharacterCreation.Worn[Pants].Custom = []Color{Color([]byte{
			'#',
			palette[r.Intn(len(palette))],
			palette[r.Intn(len(palette))],
			palette[r.Intn(len(palette))],
		})}
	case "accept":
		p.CharacterCreation.Lock()
		p.Hero = p.CharacterCreation
		p.CharacterCreation = nil
		p.backpackDirty = make(chan struct{}, 1)
		zone := p.zone
		tx, ty := p.TileX, p.TileY
		tile := zone.Tile(p.TileX, p.TileY)
		if tile == nil {
			tx, ty = 127, 127
			p.TileX, p.TileY = 127, 127
			tile = zone.Tile(127, 127)
		}
		p.Unlock()

		zone.Lock()
		tile.Add(p)
		zone.Unlock()
		SendZoneTileChange(zone.X, zone.Y, TileChange{
			ID:  p.NetworkID(),
			X:   tx,
			Y:   ty,
			Obj: p.Serialize(),
		})

		p.Lock()
		p.SetHUD("", nil)
		return
	}

	p.SetHUD("character_creation", map[string]interface{}{
		"race":   race.Name,
		"gender": genderInfo[p.CharacterCreation.Gender].Name,
		"skin":   race.SkinTones[p.CharacterCreation.SkinToneIndex],
		"shirt":  p.CharacterCreation.Worn[Shirt].Custom[0],
		"pants":  p.CharacterCreation.Worn[Pants].Custom[0],
		"name":   p.CharacterCreation.Name(),
	})
}

func playerFilename(id uint64) string {
	var buf [binary.MaxVarintLen64]byte
	i := binary.PutUvarint(buf[:], id)
	return "player" + Base32Encode(buf[:i]) + ".gz"
}

func (p *Player) Save() {
	p.Lock()
	defer p.Unlock()

	f, err := os.Create(filepath.Join(seedFilename(), playerFilename(p.ID)))
	if err != nil {
		log.Printf("[save:%d] %v", p.ID, err)
		return
	}
	defer f.Close()

	g, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		log.Printf("[save:%d] %v", p.ID, err)
		return
	}
	defer g.Close()

	err = gob.NewEncoder(g).Encode(p)
	if err != nil {
		log.Printf("[save:%d] %v", p.ID, err)
	}
}

func LoadPlayer(id uint64) (*Player, error) {
	f, err := os.Open(filepath.Join(seedFilename(), playerFilename(id)))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	g, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer g.Close()

	d := gob.NewDecoder(g)
	var p Player
	err = d.Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Player) Serialize() *NetworkedObject {
	return p.Hero.Serialize()
}

func (p *Player) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // assault
		if player == p {
			player.SendMessage("stop hitting yourself!")
			player.Lock()
			player.schedule = nil
			player.Unlock()
		} else {
			player.Lock()
			player.schedule = &CombatSchedule{
				OpponentH: p.Hero,
				OpponentP: p,
				X:         x,
				Y:         y,
			}
			player.Unlock()
		}
	}
}

func (p *Player) ZIndex() int {
	if p.Admin {
		return 1 << 30
	}
	return p.Hero.ZIndex()
}

func (p *Player) Examine() string {
	examine := p.Hero.Examine()
	p.Lock()
	defer p.Unlock()
	if p.Admin {
		examine += "\n" + p.Name() + " is an admin."
		if p.Examine_ != "" {
			examine += "\n\"" + p.Examine_ + "\""
		}
	}
	if len(p.Ancestry) > 0 {
		examine += "\n" + genderInfo[p.Gender].PronounSubject + " is the " + genderInfo[p.Gender].NounOffspring + " of " + p.Ancestry[len(p.Ancestry)-1].Hero.Name() + "."
	}
	return examine
}

func (p *Player) Think(z *Zone, x, y uint8) {
	p.think(z, x, y, p)
}

type Hero struct {
	networkID
	*HeroName

	CustomColor Color

	Gender     Gender
	Race       Race
	Occupation Occupation

	SkinToneIndex uint8

	lock  sync.Mutex
	Delay uint

	Backpack []Object
	Worn     []Cosmetic
	Toolbelt struct {
		*Hatchet
		*Pickaxe
	}

	DeathCause string

	Damage uint64

	Birth time.Time

	schedule      Schedule
	scheduleDelay uint

	combat uint8

	backpackDirty chan struct{}
}

func (h *Hero) Lock() {
	h.lock.Lock()
}

func (h *Hero) Unlock() {
	h.lock.Unlock()
}

func (h *Hero) Examine() string {
	text := append(append(append([]byte(h.Name()), " is a "...), raceInfo[h.Race].Name...), " hero."...)
	for i := range h.Worn {
		if h.Worn[i].Exists() {
			switch CosmeticType(i) {
			case Weapon:
				text = append(append(append(text, '\n'), genderInfo[h.Gender].PronounSubject...), " is holding "...)
				text = append(append(append(text, h.Worn[i].Article()...), h.Worn[i].Name()...), " in "...)
				text = append(append(append(append(text, genderInfo[h.Gender].PronounPosessive...), ' '), CosmeticSlotName[h.Worn[i].Type]...), '.')
			default:
				text = append(append(append(text, '\n'), genderInfo[h.Gender].PronounSubject...), " is wearing "...)
				text = append(append(append(text, h.Worn[i].Article()...), h.Worn[i].Name()...), " on "...)
				text = append(append(append(append(text, genderInfo[h.Gender].PronounPosessive...), ' '), CosmeticSlotName[h.Worn[i].Type]...), '.')
			}
		}
	}
	return string(text)
}

func (h *Hero) Weight() uint64 {
	var total uint64
	// toolbelted and worn items are counted as half of their weight.
	for i := range h.Worn {
		total += h.Worn[i].Weight() / 2
	}
	if h.Toolbelt.Pickaxe != nil {
		total += h.Toolbelt.Pickaxe.Weight() / 2
	}
	if h.Toolbelt.Hatchet != nil {
		total += h.Toolbelt.Hatchet.Weight() / 2
	}
	for _, o := range h.Backpack {
		if i, ok := o.(Item); ok {
			total += i.Weight()
		}
	}
	return total
}

func (h *Hero) MaxWeight() uint64 {
	return 100 * 1000
}

func (h *Hero) Volume() uint64 {
	var total uint64
	// toolbelted and worn items are not counted in volume.
	for _, o := range h.Backpack {
		if i, ok := o.(Item); ok {
			total += i.Volume()
		}
	}
	return total
}

func (h *Hero) MaxVolume() uint64 {
	return 10 * 100 * 100 * 100
}

func (h *Hero) Health() uint64 {
	health := h.MaxHealth() - h.Damage
	if health > h.MaxHealth() {
		return 0
	}
	return health
}

func (h *Hero) MaxHealth() uint64 {
	health := raceInfo[h.Race].BaseHealth
	for i := range h.Worn {
		if h.Worn[i].Exists() {
			health += h.Worn[i].HealthBonus()
		}
	}
	return health
}

func (h *Hero) Armor() uint64 {
	armor := raceInfo[h.Race].BaseArmor
	for i := range h.Worn {
		if h.Worn[i].Exists() {
			armor += h.Worn[i].ArmorBonus()
		}
	}
	return armor
}

func (h *Hero) MaxDamage() uint64 {
	damage := raceInfo[h.Race].BaseDamage
	for i := range h.Worn {
		if h.Worn[i].Exists() {
			damage += h.Worn[i].DamageBonus()
		}
	}
	return damage
}

func (h *Hero) Accuracy() uint64 {
	accuracy := raceInfo[h.Race].BaseAccuracy
	for i := range h.Worn {
		if h.Worn[i].Exists() {
			accuracy += h.Worn[i].AccuracyBonus()
		}
	}
	return accuracy
}

func (h *Hero) Blocking() bool {
	return false
}

func (h *Hero) Serialize() *NetworkedObject {
	h.Lock()
	defer h.Unlock()

	colors := []Color{h.CustomColor}
	if colors[0] == "" {
		colors[0] = raceInfo[h.Race].SkinTones[h.SkinToneIndex]
	}
	var attached []*NetworkedObject
	for _, slot := range CosmeticSlotOrder {
		if h.Worn[slot].Exists() {
			attached = append(attached, h.Worn[slot].Serialize())
		}
	}
	var health *float64
	if h.combat > 0 {
		health_ := float64(h.Health()) / float64(h.MaxHealth())
		health = &health_
	}
	return &NetworkedObject{
		Name:     h.Name(),
		Sprite:   "body_" + raceInfo[h.Race].Name,
		Colors:   colors,
		Attached: attached,
		Moves:    true,
		Health:   health,
		Options:  []string{"assault"},
	}
}

func (h *Hero) Interact(x, y uint8, player *Player, zone *Zone, opt int) {
	switch opt {
	case 0: // assault
		player.Lock()
		player.schedule = &CombatSchedule{
			OpponentH: h,
			X:         x,
			Y:         y,
		}
		tx, ty := player.TileX, player.TileY
		player.Unlock()

		h.Lock()
		h.schedule = &CombatSchedule{
			OpponentH: player.Hero,
			OpponentP: player,
			X:         tx,
			Y:         ty,
		}
		h.Unlock()
	}
}

func (h *Hero) TakeDamage(damage uint64, deathcause string) {
	max := h.MaxHealth()
	if h.Damage >= max {
		return
	}
	if damage >= max || h.Damage >= max-damage {
		h.DeathCause = deathcause
	}
	h.Damage += damage
}

func (h *Hero) Think(z *Zone, x, y uint8) {
	h.think(z, x, y, nil)
}

func (h *Hero) think(z *Zone, x, y uint8, p *Player) {
	h.Lock()

	if h.Health() == 0 {
		h.Unlock()
		if p != nil {
			p.CharacterCreationCommand("")
		}
		return
	}

	if len(h.Worn) != int(cosmeticTypeCount) {
		worn := make([]Cosmetic, cosmeticTypeCount)
		copy(worn, h.Worn)
		h.Worn = worn
	}

	if p == nil || !p.Admin {
		for i := 0; i < len(h.Backpack); i++ {
			o := h.Backpack[i]
			if a, ok := o.(Item); !ok || a.AdminOnly() {
				if p != nil {
					AdminLog.Printf("AUTOREMOVE ADMIN ITEM [%d:%q] (%d:%d %d:%d) %q %q", p.ID, p.Name(), p.ZoneX, p.TileX, p.ZoneY, p.TileY, o.Name(), o.Examine())
				}
				h.Backpack = append(h.Backpack[:i], h.Backpack[i+1:]...)
				i--
			}
		}
		for i := range h.Worn {
			o := &h.Worn[i]
			if o.Exists() && o.AdminOnly() {
				if p != nil {
					AdminLog.Printf("AUTOREMOVE ADMIN ITEM [%d:%q] (%d:%d %d:%d) %q %q", p.ID, p.Name(), p.ZoneX, p.TileX, p.ZoneY, p.TileY, o.Name(), o.Examine())
				}
				h.Worn[i] = Cosmetic{}
			}
		}
		if o := h.Toolbelt.Pickaxe; o != nil && o.AdminOnly() {
			if p != nil {
				AdminLog.Printf("AUTOREMOVE ADMIN ITEM [%d:%q] (%d:%d %d:%d) %q %q", p.ID, p.Name(), p.ZoneX, p.TileX, p.ZoneY, p.TileY, o.Name(), o.Examine())
			}
			h.Toolbelt.Pickaxe = nil
		}
		if o := h.Toolbelt.Hatchet; o != nil && o.AdminOnly() {
			if p != nil {
				AdminLog.Printf("AUTOREMOVE ADMIN ITEM [%d:%q] (%d:%d %d:%d) %q %q", p.ID, p.Name(), p.ZoneX, p.TileX, p.ZoneY, p.TileY, o.Name(), o.Examine())
			}
			h.Toolbelt.Hatchet = nil
		}
	}

	select {
	case <-h.backpackDirty:
		var obj Object = h
		if p != nil {
			obj = p
			inventory := make([]InventoryItem, 0, len(h.Backpack))
			for _, o := range h.Backpack {
				inventory = append(inventory, InventoryItem{
					ID:  o.NetworkID(),
					Obj: o.Serialize(),
				})
			}
			p.inventory <- packetInventory{
				Inventory: inventory,
			}
		}
		h.Unlock()
		SendZoneTileChange(z.X, z.Y, TileChange{
			ID:  obj.NetworkID(),
			X:   x,
			Y:   y,
			Obj: obj.Serialize(),
		})
		h.Lock()
	default:
	}

	if h.combat != 0 {
		h.combat--
		if h.combat == 0 {
			var obj Object = h
			if p != nil {
				obj = p
			}
			h.Unlock()
			SendZoneTileChange(z.X, z.Y, TileChange{
				ID:  obj.NetworkID(),
				X:   x,
				Y:   y,
				Obj: obj.Serialize(),
			})
			h.Lock()
		}
	}

	if h.Delay > 0 {
		if h.scheduleDelay > 0 {
			h.scheduleDelay--
		}
		h.Delay--
		h.Unlock()
		return
	}

	if h.scheduleDelay > 0 {
		h.scheduleDelay--
		h.Unlock()
		return
	}

	if schedule := h.schedule; schedule != nil {
		h.Unlock()
		if !schedule.Act(z, x, y, h, p) {
			h.Lock()
			h.schedule = nil
			h.Unlock()
		}
		if p == nil {
			h.Lock()
			h.scheduleDelay += h.scheduleDelay / 2
			h.Unlock()
		}
		return
	}

	if h.Damage != 0 {
		if p != nil {
			h.Damage--
		} else {
			h.Damage /= 2
		}
	}

	h.Unlock()

	if p != nil {
		return
	}

	goalX, goalY := x+uint8(rand.Intn(16)-rand.Intn(16)), y+uint8(rand.Intn(16)-rand.Intn(16))
	z.Lock()
	blocked := z.Blocked(goalX, goalY)
	z.Unlock()

	h.Lock()
	if !blocked {
		schedule := MoveSchedule(FindPath(z, x, y, goalX, goalY, true))
		h.schedule = &schedule
	}
	h.Delay = uint(rand.Intn(5) + 1)
	h.scheduleDelay = uint(rand.Intn(100) + 1)
	h.Unlock()
}

func (h *Hero) ZIndex() int {
	return 50
}

func (h *Hero) GiveItem(o Object) bool {
	if i, ok := o.(Item); ok {
		if h.MaxVolume() < h.Volume() || h.MaxVolume()-h.Volume() < i.Volume() {
			return false
		}
	}
	h.ForceGiveItem(o)
	return true
}

func (h *Hero) ForceGiveItem(o Object) {
	h.Backpack = append(h.Backpack, o)
	select {
	case h.backpackDirty <- struct{}{}:
	default:
	}
}

func (h *Hero) RemoveItem(slot int, o Object) bool {
	if slot < 0 || slot >= len(h.Backpack) || h.Backpack[slot] != o {
		return false
	}
	h.Backpack = append(h.Backpack[:slot], h.Backpack[slot+1:]...)
	select {
	case h.backpackDirty <- struct{}{}:
	default:
	}
	return true
}

func (h *Hero) Equip(o Object, inventoryOnly bool) {
	h.Lock()
	defer h.Unlock()

	index := -1
	for i, io := range h.Backpack {
		if io == o {
			index = i
			break
		}
	}
	if inventoryOnly && index == -1 {
		return
	}
	var old Object
	switch i := o.(type) {
	case *Cosmetic:
		oldMax := h.MaxHealth()
		if h.Worn[i.Type].Exists() {
			tmp := h.Worn[i.Type]
			old = &tmp
		}
		h.Worn[i.Type] = *i

		// keep health at the same percent when switching armor
		h.Damage = h.Damage * h.MaxHealth() / oldMax
	case *Pickaxe:
		if h.Toolbelt.Pickaxe != nil {
			old = h.Toolbelt.Pickaxe
		}
		h.Toolbelt.Pickaxe = i
	case *Hatchet:
		if h.Toolbelt.Hatchet != nil {
			old = h.Toolbelt.Hatchet
		}
		h.Toolbelt.Hatchet = i
	default:
		return
	}

	if index == -1 && old != nil {
		h.GiveItem(old)
	} else if index != -1 && old == nil {
		h.Backpack = append(h.Backpack[:index], h.Backpack[index+1:]...)
	} else if index != -1 && old != nil {
		h.Backpack[index] = old
	}
	select {
	case h.backpackDirty <- struct{}{}:
	default:
	}
}

type Schedule interface {
	Act(*Zone, uint8, uint8, *Hero, *Player) bool
	NextMove(uint8, uint8) (uint8, uint8)
}

type ScheduleSchedule []Schedule

func (s *ScheduleSchedule) Act(z *Zone, x, y uint8, h *Hero, p *Player) bool {
	if len(*s) == 0 {
		return false
	}

	if (*s)[0].Act(z, x, y, h, p) {
		return true
	}

	*s = (*s)[1:]
	return len(*s) != 0
}

func (s *ScheduleSchedule) NextMove(x, y uint8) (uint8, uint8) {
	if len(*s) == 0 {
		return x, y
	}
	return (*s)[0].NextMove(x, y)
}

type MoveSchedule [][2]uint8

func (s *MoveSchedule) Act(z *Zone, x, y uint8, h *Hero, p *Player) bool {
	if len(*s) == 0 {
		return false
	}
	pos := (*s)[0]
	*s = (*s)[1:]

	z.Lock()
	if z.Blocked(pos[0], pos[1]) {
		z.Unlock()
		return false
	}
	z.Unlock()

	h.Lock()
	h.Delay = uint(2 + h.Weight()/h.MaxWeight()*5)
	h.scheduleDelay = 3
	obj := Object(h)
	if p != nil {
		obj = p
		p.TileX = pos[0]
		p.TileY = pos[1]
	}
	h.Unlock()

	z.Lock()
	if z.Tile(x, y).Remove(obj) {
		z.Tile(pos[0], pos[1]).Add(obj)
	}
	z.Unlock()
	SendZoneTileChange(z.X, z.Y, TileChange{
		X:  pos[0],
		Y:  pos[1],
		ID: obj.NetworkID(),
	})
	return true
}

func (s *MoveSchedule) NextMove(x, y uint8) (uint8, uint8) {
	if len(*s) == 0 {
		return x, y
	}
	if (*s)[0][0] == x && (*s)[0][1] == y && len(*s) > 1 {
		return (*s)[1][0], (*s)[1][1]
	}
	return (*s)[0][0], (*s)[0][1]
}

type TakeSchedule struct {
	Item Object
}

func (s *TakeSchedule) Act(z *Zone, x, y uint8, h *Hero, p *Player) bool {
	tile := z.Tile(x, y)
	z.Lock()
	removed := tile.Remove(s.Item)
	z.Unlock()

	if removed {
		h.Lock()
		success := h.GiveItem(s.Item)
		h.Unlock()
		if success {
			SendZoneTileChange(z.X, z.Y, TileChange{
				ID:      s.Item.NetworkID(),
				Removed: true,
			})
		} else {
			z.Lock()
			tile.Add(s.Item)
			z.Unlock()
		}
	}
	return false
}

func (s *TakeSchedule) NextMove(x, y uint8) (uint8, uint8) {
	return x, y
}

type CombatSchedule struct {
	OpponentH *Hero
	OpponentP *Player
	X, Y      uint8
}

func (s *CombatSchedule) Act(z *Zone, x, y uint8, h *Hero, p *Player) bool {
	if dx := int(x) - int(s.X); dx < -1 || dx > 1 {
		return false
	}
	if dy := int(y) - int(s.Y); dy < -1 || dy > 1 {
		return false
	}

	z.Lock()
	found := false
	for _, o := range z.Tile(s.X, s.Y).Objects {
		if o == s.OpponentH || o == s.OpponentP {
			found = true
			break
		}
	}
	z.Unlock()
	if !found {
		return false
	}

	h.Lock()
	wasCombat := h.combat > 0
	h.combat = 50
	if h.Worn[Weapon].Exists() {
		h.Delay = uint(h.Worn[Weapon].AttackSpeed() - 1)
	} else {
		h.Delay = 14
	}
	h.scheduleDelay = h.Delay + 1

	var damage uint64
	if true { // TODO: accuracy/armor
		damage = uint64(rand.Int63n(int64(h.MaxDamage())))
	}
	h.Unlock()
	if !wasCombat {
		if p == nil {
			SendZoneTileChange(z.X, z.Y, TileChange{
				ID:  h.NetworkID(),
				X:   x,
				Y:   y,
				Obj: h.Serialize(),
			})
		} else {
			SendZoneTileChange(z.X, z.Y, TileChange{
				ID:  p.NetworkID(),
				X:   x,
				Y:   y,
				Obj: p.Serialize(),
			})
		}
	}

	s.OpponentH.Lock()
	shouldUpdate := s.OpponentH.combat == 0
	s.OpponentH.combat = 50
	if damage > 0 {
		shouldUpdate = true
		s.OpponentH.TakeDamage(damage, "slain by "+h.Name())
	}
	s.OpponentH.Unlock()
	if shouldUpdate {
		if s.OpponentP == nil {
			SendZoneTileChange(z.X, z.Y, TileChange{
				ID:  s.OpponentH.NetworkID(),
				X:   s.X,
				Y:   s.Y,
				Obj: s.OpponentH.Serialize(),
			})
		} else {
			SendZoneTileChange(z.X, z.Y, TileChange{
				ID:  s.OpponentP.NetworkID(),
				X:   s.X,
				Y:   s.Y,
				Obj: s.OpponentP.Serialize(),
			})
		}
	}

	return true
}

func (s *CombatSchedule) NextMove(x, y uint8) (uint8, uint8) {
	return x, y
}

func GenerateHero(race Race, r *rand.Rand) *Hero {
	return GenerateHeroOccupation(race, raceInfo[race].Occupations[r.Intn(len(raceInfo[race].Occupations))], r)
}

func GenerateHeroOccupation(race Race, occupation Occupation, r *rand.Rand) *Hero {
	h := &Hero{
		Race:       race,
		Gender:     raceInfo[race].Genders[r.Intn(len(raceInfo[race].Genders))],
		Occupation: occupation,
		Worn:       make([]Cosmetic, cosmeticTypeCount),
		Birth:      time.Now().UTC(),
	}
	switch race {
	case Human:
		h.HeroName = GenerateHumanName(r, h.Gender)
	}

	for slot, item := range occupationInfo[occupation].StartingGear {
		h.Worn[slot] = *item
	}
	if occupationInfo[occupation].StartingPickaxe != nil {
		pickaxe := *occupationInfo[occupation].StartingPickaxe
		h.Toolbelt.Pickaxe = &pickaxe
	}
	if occupationInfo[occupation].StartingHatchet != nil {
		hatchet := *occupationInfo[occupation].StartingHatchet
		h.Toolbelt.Hatchet = &hatchet
	}

	h.SkinToneIndex = uint8(r.Intn(len(raceInfo[race].SkinTones)))
	const pastels = "abcde"
	const earthy = "34567"
	palette := pastels
	if r.Intn(2) == 0 {
		palette = earthy
	}
	h.Worn[Shirt].Custom = []Color{Color([]byte{
		'#',
		palette[r.Intn(len(palette))],
		palette[r.Intn(len(palette))],
		palette[r.Intn(len(palette))],
	})}
	palette = earthy
	if r.Intn(3) == 0 {
		palette = pastels
	}
	h.Worn[Pants].Custom = []Color{Color([]byte{
		'#',
		palette[r.Intn(len(palette))],
		palette[r.Intn(len(palette))],
		palette[r.Intn(len(palette))],
	})}
	return h
}
