package main

import (
	"compress/gzip"
	"encoding/base32"
	"encoding/binary"
	"encoding/gob"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
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

type loadedZone struct {
	zone *Zone
	ref  uint
}

var loadedZones = make(map[[2]int64]*loadedZone)
var loadedZoneLock sync.Mutex

func GrabZone(x, y int64) *Zone {
	loadedZoneLock.Lock()
	defer loadedZoneLock.Unlock()

	if z, ok := loadedZones[[2]int64{x, y}]; ok {
		z.ref++
		return z.zone
	}

	z, err := LoadZone(x, y)
	if err != nil {
		z = &Zone{X: x, Y: y}
		z.Generate()
	}

	loadedZones[[2]int64{x, y}] = &loadedZone{
		zone: z,
		ref:  1,
	}
	return z
}

func ReleaseZone(z *Zone) {
	loadedZoneLock.Lock()
	defer loadedZoneLock.Unlock()

	l := loadedZones[[2]int64{z.X, z.Y}]
	l.ref--
	if l.ref == 0 {
		l.zone.Save()
		delete(loadedZones, [2]int64{z.X, z.Y})
	}
}

func EachLoadedZone(f func(*Zone)) {
	loadedZoneLock.Lock()
	for _, z := range loadedZones {
		loadedZoneLock.Unlock()
		f(z.zone)
		loadedZoneLock.Lock()
	}
	loadedZoneLock.Unlock()
}

func seedFilename() string {
	var buf [binary.MaxVarintLen64]byte
	i := binary.PutVarint(buf[:], Seed)
	encoded := base32.StdEncoding.EncodeToString(buf[:i])

	l := len(encoded)
	for encoded[l-1] == '=' {
		l--
	}
	return "rnoadm-" + encoded[:l]
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
	Seed    RandomSource
	X, Y    int64
	Element Element
	Tiles   [zoneTiles]Tile
	Name_   *Name
	dirty   chan struct{}
	mtx     sync.Mutex
}

func (z *Zone) Lock() {
	z.mtx.Lock()
}

func (z *Zone) Unlock() {
	z.mtx.Unlock()
}

func (z *Zone) Blocked(x, y uint8) bool {
	return z.Tile(x, y).Blocked()
}

func (z *Zone) Tile(x, y uint8) *Tile {
	if zoneOffset[y] > x || x > 255-zoneOffset[y] {
		return nil
	}
	return &z.Tiles[rowOffset[y]+int(x-zoneOffset[y])]
}

func (z *Zone) Rand() *rand.Rand {
	return rand.New(&z.Seed)
}

func (z *Zone) Generate() {
	z.Lock()
	defer z.Unlock()

	z.dirty = make(chan struct{}, 1)

	z.Seed.Seed(Seed ^ z.X ^ int64(uint64(z.Y)<<32|uint64(z.Y)>>32))
	r := z.Rand()
	z.Element = Nature
	z.Name_ = GenerateName(r, NameZone, NamePlains)
	for i := r.Intn(100); i > 0; i-- {
		x := r.Float64()*192 + 32
		y := r.Float64()*192 + 32
		rock, ok := z.Element.Rock(r)
		if !ok {
			break
		}

		for j := 0; j < 40; j++ {
			radius := r.Float64() * 4
			theta := r.Float64() * 2 * math.Pi
			tile := z.Tile(uint8(x+radius*math.Cos(theta)), uint8(y+radius*math.Sin(theta)))

			if !tile.Blocked() {
				tile.Add(&Rock{
					Type: rock,
				})
			}
		}
	}
	for i := range z.Tiles {
		if !z.Tiles[i].Blocked() && r.Intn(4) == 0 {
			flora, ok := z.Element.Flora(r)
			if !ok {
				break
			}
			z.Tiles[i].Add(&Flora{
				Type: flora,
			})
		}
	}
}

func (z *Zone) Save() error {
	players := make(map[int][]*Player)

	z.Lock()
	defer z.Unlock()

	// remove and save players before saving the zone
	for i := range z.Tiles {
		for j := 0; j < len(z.Tiles[i].Objects); j++ {
			if p, ok := z.Tiles[i].Objects[j].(*Player); ok {
				z.Tiles[i].Objects = append(z.Tiles[i].Objects[:j], z.Tiles[i].Objects[j+1:]...)
				players[i] = append(players[i], p)
				p.Save()
			}
		}
	}
	defer func() {
		// add players back in
		for i, l := range players {
			for _, p := range l {
				z.Tiles[i].Add(p)
			}
		}
	}()

	dir := seedFilename()
	os.MkdirAll(dir, 0755)

	f, err := os.Create(filepath.Join(dir, zoneFilename(z.X, z.Y)))
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
	return e.Encode(z)
}

func LoadZone(x, y int64) (*Zone, error) {
	dir := seedFilename()

	f, err := os.Open(filepath.Join(dir, zoneFilename(x, y)))
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
	var z Zone
	err = d.Decode(&z)
	if err != nil {
		return nil, err
	}
	z.dirty = make(chan struct{}, 1)
	return &z, nil
}

func (z *Zone) Think() {
	z.Lock()
	defer z.Unlock()

	for i := range z.Tiles {
		for _, o := range z.Tiles[i].Objects {
			if t, ok := o.(Thinker); ok {
				t.Think()
			}
		}
	}
	select {
	case <-z.dirty:
		for i := range z.Tiles {
			for _, o := range z.Tiles[i].Objects {
				if p, ok := o.(*Player); ok {
					p.Repaint()
				}
			}
		}
	default:
	}
}

func (z *Zone) Repaint() {
	select {
	case z.dirty <- struct{}{}:
	default:
	}
}

func (z *Zone) Name() string {
	return z.Name_.String()
}

type Tile struct {
	Objects []Object
	Ground  rune
}

func (t *Tile) Add(o Object) {
	t.Objects = append(t.Objects, o)
}

func (t *Tile) Remove(o Object) bool {
	for i := range t.Objects {
		if t.Objects[i] == o {
			t.Objects = append(t.Objects[:i], t.Objects[i+1:]...)
			return true
		}
	}
	return false
}

func (t *Tile) Blocked() bool {
	if t == nil {
		return true
	}
	for _, o := range t.Objects {
		if o.Blocking() {
			return true
		}
	}
	return false
}

type Color string

func (t *Tile) Paint(z *Zone) (rune, Color) {
	if t.Ground == 0 {
		const ground = " ,.'-`"
		t.Ground = rune(ground[rand.Intn(len(ground))])
	}
	if len(t.Objects) == 0 {
		return t.Ground, "#333"
	}
	return t.Objects[len(t.Objects)-1].Paint()
}

type Object interface {
	Name() string
	Examine() string
	Paint() (rune, Color)
	Blocking() bool
	InteractOptions() []string
}

type Item interface {
	IsItem()
}

var _ Item = (*Ore)(nil)
var _ Item = (*Stone)(nil)

type Thinker interface {
	Think()
}

var _ Thinker = (*Player)(nil)
var _ Thinker = (*Hero)(nil)

func init() {
	gob.Register(Object(&Hero{}))
	// Players are removed from Zones before saving.
	gob.Register(Object(&Flora{}))
	gob.Register(Object(&Tree{}))
	gob.Register(Object(&Rock{}))
	gob.Register(Object(&Stone{}))
	gob.Register(Object(&Ore{}))
	gob.Register(Object(&Wall{}))
}
