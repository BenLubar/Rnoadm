package main

import (
	"encoding/base32"
	"encoding/binary"
	"encoding/gob"
	"math"

	"github.com/nsf/termbox-go"
)

const root3 = 1.7320508075688772935274463415058723669428052538103806
const zoneTiles = 46872

var zoneOffset [256]uint8
var rowOffset [256]int

func init() {
	sum := 0
	for i := range zoneOffset {
		zoneOffset[i] = uint8(math.Abs(float64(i)-127.5) / root3)
		rowOffset[i] = sum
		sum += 256 - int(zoneOffset[i])*2
	}
	if zoneTiles != sum {
		panic(sum)
	}
}

func zoneFilename(x, y int64) string {
	var buf [binary.MaxVarintLen64 * 2]byte
	i := binary.PutVarint(buf[:], x)
	i += binary.PutVarint(buf[i:], y)
	encoded := base32.StdEncoding.EncodeToString(buf[:i])

	l := len(encoded)
	for encoded[l-1] == '=' {
		l--
	}
	return "z" + encoded[:l] + ".gz"
}

type Zone struct {
	X, Y    int64
	Element Element
	Tiles   [zoneTiles]Tile
}

func (z *Zone) Blocked(x, y uint8) bool {
	tile := z.Tile(x, y)
	if tile == nil {
		return true
	}
	return tile.Blocked()
}

func (z *Zone) Tile(x, y uint8) *Tile {
	if zoneOffset[y] > x || x > 255-zoneOffset[y] {
		return nil
	}
	return &z.Tiles[rowOffset[y]+int(x-zoneOffset[y])]
}

type Tile struct {
	Objects []Object
}

func (t *Tile) Add(o Object) {
	t.Objects = append(t.Objects, o)
}

func (t *Tile) Blocked() bool {
	for _, o := range t.Objects {
		if o.Blocking() {
			return true
		}
	}
	return false
}

func (t *Tile) Paint() (rune, termbox.Attribute) {
	if len(t.Objects) == 0 {
		return '.', termbox.ColorWhite
	}
	return t.Objects[len(t.Objects)-1].Paint()
}

type Object interface {
	Name() string
	Examine() string
	Paint() (rune, termbox.Attribute)
	Blocking() bool
}

func init() {
	gob.Register(Object(&Rock{}))
}
