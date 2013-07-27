package main

import (
	"compress/gzip"
	"encoding/base32"
	"encoding/binary"
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/nsf/termbox-go"
)

type Player struct {
	ID uint64
	Hero
	ZoneX, ZoneY int64
	TileX, TileY uint8
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
	return &p, nil
}

type Hero struct {
	Name_ string
}

func (h *Hero) Name() string {
	if h.Name_ == "" {
		h.Name_ = "hero" // TODO: name generator
	}
	return h.Name_
}

func (h *Hero) Examine() string {
	return "a hero."
}

func (h *Hero) Blocking() bool {
	return false
}

func (h *Hero) Paint() (rune, termbox.Attribute) {
	return '☻', termbox.ColorWhite
}
