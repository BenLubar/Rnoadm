package main

import (
	"compress/gzip"
	"encoding/base32"
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

var genderData = [genderCount]struct {
	Name string
}{
	Male: {
		Name: "male",
	},
	Female: {
		Name: "female",
	},
}

type Race uint16

const (
	Human Race = iota

	raceCount
)

var raceInfo = [raceCount]struct {
	Name      string
	Genders   []Gender
	SkinTones []Color
}{
	Human: {
		Name:      "human",
		Genders:   []Gender{Male, Female},
		SkinTones: []Color{"#ffe3cc", "#ffdbbd", "#e6c0a1", "#edd0b7", "#e3c3a8", "#ffcda3", "#e8d1be", "#e6d2c1", "#f7e9dc", "#cfbcab", "#c2a38a", "#c9a281", "#d9a980", "#ba9c82", "#ad8f76", "#a17a5a", "#876d58", "#6e5948"},
	},
}

type Occupation uint16

const (
	Civilian Occupation = iota

	occupationCount
)

type Player struct {
	ID uint64
	Hero
	ZoneX, ZoneY int64
	TileX, TileY uint8

	hud interface {
		Paint(func(int, int, PaintCell))
		Key(int, bool) bool
		Click(int, int) bool
	}
	repaint chan struct{}

	Joined    time.Time
	LastLogin time.Time
	Admin     bool
	Examine_  string

	messages chan<- string

	zone *Zone
}

func (p *Player) SendMessage(message string) {
	select {
	case p.messages <- message:
		p.Repaint()
	case <-time.After(time.Second):
	}
}

func (p *Player) Move(dx, dy int) {
	if p.Delay > 0 {
		return
	}
	destX := dx + int(p.TileX)
	destY := dy + int(p.TileY)

	zoneChange := destX < 0 || destY < 0 || destX > 255 || destY > 255

	p.Lock()
	z := p.zone
	p.Unlock()
	zoneChange = zoneChange || z.Tile(uint8(destX), uint8(destY)) == nil
	z.Lock()
	if !zoneChange && z.Blocked(uint8(destX), uint8(destY)) {
		z.Unlock()
		return
	}
	z.Unlock()

	if zoneChange {
		if destY >= 0 && destY <= 255 {
			// TEMPORARY: no zone changes to the sides due to bugs
			return
		}
		z.Lock()
		z.Tile(p.TileX, p.TileY).Remove(p)
		z.Unlock()
		z.Repaint()
		ReleaseZone(z)
		p.Lock()
		if destY < 0 {
			p.ZoneY -= 2
			p.TileY = 255
		} else if destY > 255 {
			p.ZoneY += 2
			p.TileY = 0
		}
		p.Delay = 2
		z = GrabZone(p.ZoneX, p.ZoneY)
		p.zone = z
		p.Unlock()
		p.Save()
		p.hud = nil
	} else {
		z.Lock()
		z.Tile(p.TileX, p.TileY).Remove(p)
		z.Unlock()
		p.Lock()
		p.TileX = uint8(destX)
		p.TileY = uint8(destY)
		p.Delay = 2
		p.Unlock()
	}
	z.Lock()
	z.Tile(p.TileX, p.TileY).Add(p)
	z.Unlock()
	z.Repaint()
}

func playerFilename(id uint64) string {
	var buf [binary.MaxVarintLen64]byte
	i := binary.PutUvarint(buf[:], id)
	encoded := base32.StdEncoding.EncodeToString(buf[:i])

	l := len(encoded)
	for encoded[l-1] == '=' {
		l--
	}
	return "p" + encoded[:l] + ".gz"
}

func (p *Player) Save() {
	p.Lock()
	defer p.Unlock()

	dir := seedFilename()

	f, err := os.Create(filepath.Join(dir, playerFilename(p.ID)))
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
	dir := seedFilename()

	f, err := os.Open(filepath.Join(dir, playerFilename(id)))
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
	p.repaint = make(chan struct{}, 1)
	return &p, nil
}

func (p *Player) ZIndex() int {
	if p.Admin {
		return 1 << 30
	}
	return p.Hero.ZIndex()
}

func (p *Player) Repaint() {
	select {
	case p.repaint <- struct{}{}:
	default:
	}
}

func (p *Player) Examine() string {
	if p.Admin {
		p.Lock()
		defer p.Unlock()
		if p.Examine_ != "" {
			return p.Examine_
		}
		return "an admin."
	}
	return p.Hero.Examine()
}

func (p *Player) Think(z *Zone, x, y uint8) {
	p.think(z, x, y, p)
}

type Hero struct {
	*HeroName

	CustomColor Color

	Gender     Gender
	Race       Race
	Occupation Occupation

	SkinToneIndex uint8

	lock  sync.Mutex
	Delay uint

	Backpack []Object
	Head     *Hat
	Top      *Shirt
	Legs     *Pants
	Feet     *Shoes
	Toolbelt struct {
		*Hatchet
		*Pickaxe
	}

	schedule      Schedule
	scheduleDelay uint
}

func (h *Hero) Lock() {
	h.lock.Lock()
}

func (h *Hero) Unlock() {
	h.lock.Unlock()
}

func (h *Hero) Examine() string {
	return "a hero."
}

func (h *Hero) Blocking() bool {
	return false
}

func (h *Hero) Paint(x, y int, setcell func(int, int, PaintCell)) {
	h.Lock()
	defer h.Unlock()

	color := h.CustomColor
	if color == "" {
		color = raceInfo[h.Race].SkinTones[h.SkinToneIndex]
	}
	setcell(x, y, PaintCell{
		Sprite: "body_human",
		Color:  color,
	})
	if h.Feet != nil {
		h.Feet.PaintWorn(x, y, setcell)
	}
	if h.Legs != nil {
		h.Legs.PaintWorn(x, y, setcell)
	}
	if h.Top != nil {
		h.Top.PaintWorn(x, y, setcell)
	}
	if h.Head != nil {
		h.Head.PaintWorn(x, y, setcell)
	}
	if h.Toolbelt.Pickaxe != nil {
		h.Toolbelt.Pickaxe.PaintWorn(x, y, setcell)
	}
	if h.Toolbelt.Hatchet != nil {
		h.Toolbelt.Hatchet.PaintWorn(x, y, setcell)
	}
}

func (h *Hero) Think(z *Zone, x, y uint8) {
	h.think(z, x, y, nil)
}

func (h *Hero) think(z *Zone, x, y uint8, p *Player) {
	h.Lock()

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
		if o := h.Head; o != nil && o.AdminOnly() {
			if p != nil {
				AdminLog.Printf("AUTOREMOVE ADMIN ITEM [%d:%q] (%d:%d %d:%d) %q %q", p.ID, p.Name(), p.ZoneX, p.TileX, p.ZoneY, p.TileY, o.Name(), o.Examine())
			}
			h.Head = nil
		}
		if o := h.Top; o != nil && o.AdminOnly() {
			if p != nil {
				AdminLog.Printf("AUTOREMOVE ADMIN ITEM [%d:%q] (%d:%d %d:%d) %q %q", p.ID, p.Name(), p.ZoneX, p.TileX, p.ZoneY, p.TileY, o.Name(), o.Examine())
			}
			h.Top = nil
		}
		if o := h.Legs; o != nil && o.AdminOnly() {
			if p != nil {
				AdminLog.Printf("AUTOREMOVE ADMIN ITEM [%d:%q] (%d:%d %d:%d) %q %q", p.ID, p.Name(), p.ZoneX, p.TileX, p.ZoneY, p.TileY, o.Name(), o.Examine())
			}
			h.Legs = nil
		}
		if o := h.Feet; o != nil && o.AdminOnly() {
			if p != nil {
				AdminLog.Printf("AUTOREMOVE ADMIN ITEM [%d:%q] (%d:%d %d:%d) %q %q", p.ID, p.Name(), p.ZoneX, p.TileX, p.ZoneY, p.TileY, o.Name(), o.Examine())
			}
			h.Feet = nil
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

func (h *Hero) InteractOptions() []string {
	return nil
}

func (h *Hero) Interact(x uint8, y uint8, player *Player, zone *Zone, opt int) {
}

func (h *Hero) ZIndex() int {
	return 50
}

func (h *Hero) GiveItem(o Object) {
	h.Backpack = append(h.Backpack, o)
}

func (h *Hero) Equip(o Object, inventoryOnly bool) {
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
	case *Hat:
		if h.Head != nil {
			old = h.Head
		}
		h.Head = i
	case *Shirt:
		if h.Top != nil {
			old = h.Top
		}
		h.Top = i
	case *Pants:
		if h.Legs != nil {
			old = h.Legs
		}
		h.Legs = i
	case *Shoes:
		if h.Feet != nil {
			old = h.Feet
		}
		h.Feet = i
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
}

type Schedule interface {
	Act(*Zone, uint8, uint8, *Hero, *Player) bool
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
	h.Delay = 2
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
	z.Repaint()
	return true
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
		h.GiveItem(s.Item)
		h.Unlock()
		z.Repaint()
	}
	return false
}

func GenerateHero(race Race, r *rand.Rand) *Hero {
	h := &Hero{
		Race:   race,
		Gender: raceInfo[race].Genders[r.Intn(len(raceInfo[race].Genders))],
	}
	switch race {
	case Human:
		h.HeroName = GenerateHumanName(r, h.Gender)
	}
	h.SkinToneIndex = uint8(r.Intn(len(raceInfo[race].SkinTones)))
	const pastels = "abcde"
	const earthy = "34567"
	palette := pastels
	if r.Intn(2) == 0 {
		palette = earthy
	}
	h.Top = &Shirt{
		Type: HipHopTeeShirt,
		CustomColor: [5]Color{Color([]byte{
			'#',
			palette[r.Intn(len(palette))],
			palette[r.Intn(len(palette))],
			palette[r.Intn(len(palette))],
		})},
	}
	palette = earthy
	if r.Intn(3) == 0 {
		palette = pastels
	}
	h.Legs = &Pants{
		Type: OffBrandJeans,
		CustomColor: [5]Color{Color([]byte{
			'#',
			palette[r.Intn(len(palette))],
			palette[r.Intn(len(palette))],
			palette[r.Intn(len(palette))],
		})},
	}
	return h
}
