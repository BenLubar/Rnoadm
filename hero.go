package main

import (
	"compress/gzip"
	"encoding/base32"
	"encoding/binary"
	"encoding/gob"
	"os"
	"path/filepath"
)

type Player struct {
	ID uint64
	Hero
	ZoneX, ZoneY int64
	TileX, TileY uint8
	hud          interface {
		Paint(func(int, int, rune, Color))
		Key(int) bool
	}
	repaint chan struct{}
}

func (p *Player) Move(dx, dy int) {
	if p.Delay > 0 {
		return
	}
	destX := dx + int(p.TileX)
	destY := dy + int(p.TileY)

	zoneChange := destX < 0 || destY < 0 || destX > 255 || destY > 255

	z := GrabZone(p.ZoneX, p.ZoneY)
	z.Lock()
	if !zoneChange {
		zoneChange = z.Tile(uint8(destX), uint8(destY)) == nil
	}
	if !zoneChange && z.Blocked(uint8(destX), uint8(destY)) {
		z.Unlock()
		ReleaseZone(z)
		return
	}
	z.Tile(p.TileX, p.TileY).Remove(p)
	if zoneChange {
		z.Unlock()
		p.RepaintZone()
		ReleaseZone(z) // player has left zone
		ReleaseZone(z)
		if destY < 0 {
			p.ZoneY -= 2
			p.TileX = 127
			p.TileY = 255
		} else if destY > 255 {
			p.ZoneY += 2
			p.TileX = 127
			p.TileY = 0
		} else if destX < 128 {
			if destY < 128 {
				p.ZoneY--
				p.TileX = 255 - zoneOffset[255-64]
				p.TileY = 255 - 64
			} else {
				p.ZoneX--
				p.TileX = 255 - zoneOffset[64]
				p.TileY = 64
			}
		} else {
			if destY < 128 {
				p.ZoneX++
				p.ZoneY--
				p.TileX = zoneOffset[255-64]
				p.TileY = 255 - 64
			} else {
				p.ZoneY++
				p.TileX = zoneOffset[64]
				p.TileY = 64
			}
		}
		p.Save()
		z = GrabZone(p.ZoneX, p.ZoneY)
		GrabZone(p.ZoneX, p.ZoneY) // player has entered zone
		z.Lock()
		p.hud = ZoneEntryHUD(z.Name())
	} else {
		p.TileX = uint8(destX)
		p.TileY = uint8(destY)
	}
	p.Delay = 2
	z.Tile(p.TileX, p.TileY).Add(p)
	z.Unlock()
	ReleaseZone(z)
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

func (p *Player) Save() error {
	dir := seedFilename()
	os.MkdirAll(dir, 0755)

	f, err := os.Create(filepath.Join(dir, playerFilename(p.ID)))
	if err != nil {
		return err
	}
	defer f.Close()

	g, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer g.Close()

	e := gob.NewEncoder(g)
	return e.Encode(p)
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
	z := GrabZone(p.ZoneX, p.ZoneY)
	z.Repaint()
	ReleaseZone(z)
}

func (p *Player) Think() {
	p.think(false)
}

type ZoneEntryHUD string

func (zeh ZoneEntryHUD) Paint(setcell func(int, int, rune, Color)) {
	i := 0
	for _, r := range zeh {
		setcell(i, 0, r, "#fff")
		i++
	}
}

func (zeh ZoneEntryHUD) Key(code int) bool {
	return false
}

type Hero struct {
	Name_ *Name
	Delay uint
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
	if h.Delay > 0 {
		h.Delay--
		return
	}
}
