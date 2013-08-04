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

type Player struct {
	ID uint64
	*Hero
	ZoneX, ZoneY int64
	TileX, TileY uint8

	Seed RandomSource

	hud interface {
		//Paint(func(int, int, PaintCell))
		Key(int, bool) bool
		Click(int, int) bool
	}

	LastLoginAddr string
	LastLogin     time.Time
	Admin         bool
	Examine_      string

	messages chan<- string
	kick     chan<- string

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
	select {
	case p.messages <- message:
		p.Repaint()
	case <-time.After(time.Second):
	}
}

func (p *Player) Kick(message string) {
	select {
	case p.kick <- message:
	default: // player was already kicked
	}
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

func (p *Player) ZIndex() int {
	if p.Admin {
		return 1 << 30
	}
	return p.Hero.ZIndex()
}

func (p *Player) Repaint() {
	// TODO: remove this no-op
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

	frame         uint8
	tileX, tileY  uint8
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

/*func (h *Hero) Paint(x, y int, setcell func(int, int, PaintCell)) {
	h.Lock()
	defer h.Unlock()

	frame := h.frame
	var offsetX, offsetY int8
	if h.schedule != nil {
		cx, cy := h.tileX, h.tileY
		nx, ny := h.schedule.NextMove(cx, cy)
		switch h.scheduleDelay {
		case 3, 1:
			if cx > nx {
				frame = 6
			} else if cx < nx {
				frame = 9
			} else if cy > ny {
				frame = 3
			} else if cy < ny {
				frame = 0
			}
		case 2, 0:
			if cy > ny {
				frame = 3
			} else if cy < ny {
				frame = 0
			} else if cx > nx {
				frame = 6
			} else if cy > ny {
				frame = 9
			}
		}
		if h.scheduleDelay&1 == 1 {
			frame = frame%3 + uint8(h.scheduleDelay/2+1)
		}
		h.frame = frame
		offsetX = int8(nx-cx) * 16 * int8(4-h.scheduleDelay) / 4
		offsetY = int8(ny-cy) * 16 * int8(4-h.scheduleDelay) / 4
	} else {
		frame -= frame % 3
	}

	color := h.CustomColor
	if color == "" {
		color = raceInfo[h.Race].SkinTones[h.SkinToneIndex]
	}
	setcell(x, y, PaintCell{
		Sprite: "body_human",
		Color:  color,
		SheetX: frame,
		X:      offsetX,
		Y:      offsetY,
		ZIndex: 500,
	})
	if h.Feet != nil {
		h.Feet.PaintWorn(x, y, setcell, frame, offsetX, offsetY)
	}
	if h.Legs != nil {
		h.Legs.PaintWorn(x, y, setcell, frame, offsetX, offsetY)
	}
	if h.Top != nil {
		h.Top.PaintWorn(x, y, setcell, frame, offsetX, offsetY)
	}
	if h.Head != nil {
		h.Head.PaintWorn(x, y, setcell, frame, offsetX, offsetY)
	}
	if h.Toolbelt.Pickaxe != nil {
		h.Toolbelt.Pickaxe.PaintWorn(x, y, setcell, frame, offsetX, offsetY)
	}
	if h.Toolbelt.Hatchet != nil {
		h.Toolbelt.Hatchet.PaintWorn(x, y, setcell, frame, offsetX, offsetY)
	}
}*/

func (h *Hero) Think(z *Zone, x, y uint8) {
	h.think(z, x, y, nil)
}

func (h *Hero) think(z *Zone, x, y uint8, p *Player) {
	h.Lock()

	h.tileX, h.tileY = x, y

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
			z.Repaint()
		}
		h.Delay--
		h.Unlock()
		return
	}

	if h.scheduleDelay > 0 {
		h.scheduleDelay--
		h.Unlock()
		z.Repaint()
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
		h.GiveItem(s.Item)
		h.Unlock()
		z.Repaint()
	}
	return false
}

func (s *TakeSchedule) NextMove(x, y uint8) (uint8, uint8) {
	return x, y
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
