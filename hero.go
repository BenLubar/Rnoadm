package main

import (
	"compress/gzip"
	"encoding/base32"
	"encoding/binary"
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Player struct {
	ID uint64
	Hero
	ZoneX, ZoneY int64
	TileX, TileY uint8

	hud interface {
		Paint(func(int, int, rune, Color))
		Key(int, bool) bool
	}
	repaint chan struct{}

	Joined    time.Time
	LastLogin time.Time
	Admin     bool
	Examine_  string

	zone *Zone
}

func (p *Player) Move(dx, dy int) {
	if p.Delay > 0 {
		return
	}
	destX := dx + int(p.TileX)
	destY := dy + int(p.TileY)

	zoneChange := destX < 0 || destY < 0 || destX > 255 || destY > 255

	p.lock.Lock()
	z := p.zone
	p.lock.Unlock()
	z.Lock()
	if !zoneChange {
		zoneChange = z.Tile(uint8(destX), uint8(destY)) == nil
	}
	if !zoneChange && z.Blocked(uint8(destX), uint8(destY)) {
		z.Unlock()
		return
	}
	z.Tile(p.TileX, p.TileY).Remove(p)
	z.Unlock()
	if zoneChange {
		p.RepaintZone()
		ReleaseZone(z)
		p.lock.Lock()
		p.Delay = 2
		if destY < 0 {
			p.ZoneY -= 2
			p.TileY = 255
		} else if destY > 255 {
			p.ZoneY += 2
			p.TileY = 0
		} else if destX < 128 {
			p.TileX = 255 - p.TileX
			p.TileY = 255 - p.TileY
			if p.ZoneY&1 == 0 {
				p.ZoneX--
			}
			if destY < 128 {
				p.ZoneY--
			} else {
				p.ZoneY++
			}
		} else {
			p.TileX = 255 - p.TileX
			p.TileY = 255 - p.TileY
			if p.ZoneY&1 == 1 {
				p.ZoneX++
			}
			if destY < 128 {
				p.ZoneY--
			} else {
				p.ZoneY++
			}
		}
		z = GrabZone(p.ZoneX, p.ZoneY)
		p.zone = z
		p.lock.Unlock()
		p.Save()
		p.hud = nil
	} else {
		p.lock.Lock()
		p.TileX = uint8(destX)
		p.TileY = uint8(destY)
		p.Delay = 2
		p.lock.Unlock()
	}
	z.Lock()
	z.Tile(p.TileX, p.TileY).Add(p)
	z.Unlock()
	p.RepaintZone()
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
	p.lock.Lock()
	defer p.lock.Unlock()

	dir := seedFilename()
	os.MkdirAll(dir, 0755)

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

func (p *Player) Repaint() {
	select {
	case p.repaint <- struct{}{}:
	default:
	}
}

func (p *Player) RepaintZone() {
	go func() {
		p.lock.Lock()
		p.zone.Repaint()
		p.lock.Unlock()
	}()
}

func (p *Player) Examine() string {
	if p.Admin {
		p.lock.Lock()
		defer p.lock.Unlock()
		if p.Examine_ != "" {
			return p.Examine_
		}
		return "an admin."
	}
	return p.Hero.Examine()
}

func (p *Player) Paint() (rune, Color) {
	if p.Admin {
		return '♚', "#fa0"
	}
	return p.Hero.Paint()
}

func (p *Player) Think() {
	p.think(false)
}

type ZoneEntryHUD string

func (h ZoneEntryHUD) Paint(setcell func(int, int, rune, Color)) {
	i := 0
	for _, r := range h {
		setcell(i, 0, r, "#fff")
		i++
	}
}

func (h ZoneEntryHUD) Key(code int, special bool) bool {
	return false
}

type Hero struct {
	Name_ *Name

	lock     sync.Mutex
	Delay    uint
	Backpack []Object
}

func (h *Hero) Name() string {
	return h.Name_.String()
}

func (h *Hero) Examine() string {
	return "a hero."
}

func (h *Hero) Blocking() bool {
	return false
}

func (h *Hero) Paint() (rune, Color) {
	return '☻', "#fff"
}

func (h *Hero) Think() {
	h.think(true)
}

func (h *Hero) think(ai bool) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.Delay > 0 {
		h.Delay--
		return
	}
}

func (h *Hero) InteractOptions() []string {
	return nil
}

func (h *Hero) GiveItem(o Object) {
	h.Backpack = append(h.Backpack, o)
}
