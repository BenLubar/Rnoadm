package main

import (
	"compress/gzip"
	"encoding/base32"
	"encoding/binary"
	"encoding/gob"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
)

var nextNetworkID uint64

type networkID uint64

func (id *networkID) NetworkID() uint64 {
	if i := atomic.LoadUint64((*uint64)(id)); i != 0 {
		return i
	}
	atomic.CompareAndSwapUint64((*uint64)(id), 0, atomic.AddUint64(&nextNetworkID, 1))
	return atomic.LoadUint64((*uint64)(id))
}

func (id *networkID) Serialize() *NetworkedObject {
	// TODO: remove this method
	return &NetworkedObject{
		Sprite: "ui_r1",
		Colors: []Color{"#f0f"},
	}
}

func (id *networkID) Interact(x, y uint8, player *Player, zone *Zone, opt int) {}

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
	ref  map[*Player]chan<- TileChange
}

type zoneInfo struct {
	Name    string
	Element Element
	Biome   Biome
}

var loadedZones = make(map[[2]int64]loadedZone)
var loadedZoneLock sync.Mutex
var ZoneInfo map[[2]int64]zoneInfo

func GrabZone(x, y int64, p *Player) (*Zone, <-chan TileChange) {
	loadedZoneLock.Lock()
	defer loadedZoneLock.Unlock()

	ch := make(chan TileChange, 64)

	if ZoneInfo == nil {
		f, err := os.Open(filepath.Join(seedFilename(), "mZONEMETA.gz"))
		if err != nil {
			ZoneInfo = make(map[[2]int64]zoneInfo)
		} else {
			defer f.Close()

			g, err := gzip.NewReader(f)
			if err != nil {
				ZoneInfo = make(map[[2]int64]zoneInfo)
			} else {
				defer g.Close()

				err = gob.NewDecoder(g).Decode(&ZoneInfo)
				if err != nil {
					ZoneInfo = make(map[[2]int64]zoneInfo)
				}
			}
		}
	}

	if z, ok := loadedZones[[2]int64{x, y}]; ok {
		z.ref[p] = ch
		return z.zone, ch
	}

	z, err := LoadZone(x, y)
	if err != nil {
		log.Printf("ZONE %d %d: %v", x, y, err)
		z = &Zone{X: x, Y: y}
		z.Generate()
	}

	ZoneInfo[[2]int64{x, y}] = zoneInfo{
		Name:    z.Name(),
		Element: z.Element,
		Biome:   z.Biome,
	}
	f, err := os.Create(filepath.Join(seedFilename(), "mZONEMETA.gz"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	g, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	err = gob.NewEncoder(g).Encode(&ZoneInfo)
	if err != nil {
		panic(err)
	}

	lz := loadedZone{
		zone: z,
		ref:  make(map[*Player]chan<- TileChange),
	}
	lz.ref[p] = ch
	loadedZones[[2]int64{x, y}] = lz
	return z, ch
}

func ReleaseZone(z *Zone, p *Player) {
	loadedZoneLock.Lock()
	defer loadedZoneLock.Unlock()

	l := loadedZones[[2]int64{z.X, z.Y}]
	close(l.ref[p])
	delete(l.ref, p)
	if len(l.ref) == 0 {
		l.zone.Save()
		delete(loadedZones, [2]int64{z.X, z.Y})
	}
}

func SendZoneTileChange(x, y int64, c TileChange) {
	loadedZoneLock.Lock()
	for _, ch := range loadedZones[[2]int64{x, y}].ref {
		select {
		case ch <- c:
		default:
		}
	}
	loadedZoneLock.Unlock()
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

func Base32Encode(b []byte) string {
	encoded := base32.StdEncoding.EncodeToString(b)

	l := len(encoded)
	for encoded[l-1] == '=' {
		l--
	}

	return encoded[:l]
}

func seedFilename() string {
	var buf [binary.MaxVarintLen64]byte
	i := binary.PutVarint(buf[:], Seed)
	return "rnoadm-" + Base32Encode(buf[:i])
}

func zoneFilename(x, y int64) string {
	var buf [binary.MaxVarintLen64 * 2]byte
	i := binary.PutVarint(buf[:], x)
	i += binary.PutVarint(buf[i:], y)
	return "z" + Base32Encode(buf[:i]) + ".gz"
}

type Zone struct {
	*ZoneName
	Seed    RandomSource
	X, Y    int64
	Element Element
	Biome   Biome
	Tiles   [zoneTiles]Tile
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

type Biome uint8

const (
	Plains Biome = iota
	Forest
	Hills
	Lake

	biomeCount = Forest + 1 // TODO: implement other biomes
)

func (z *Zone) Generate() {
	z.Lock()
	defer z.Unlock()

	z.Seed.Seed(Seed ^ z.X ^ int64(uint64(z.Y)<<32|uint64(z.Y)>>32))
	r := z.Rand()
	z.Element = Nature
	z.Biome = Biome(r.Intn(int(biomeCount)))
	switch z.Biome {
	case Plains:
		z.generatePlains(r)
	case Forest:
		z.generateForest(r)
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
				j--
				players[i] = append(players[i], p)
				z.Unlock()
				p.Save()
				z.Lock()
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
	return &z, nil
}

func (z *Zone) Think() {
	z.Lock()

	for i := range z.Tiles {
		var x, y uint8
		for j, start := range rowOffset {
			if start <= i {
				x = uint8(i-start) + zoneOffset[j]
				y = uint8(j)
			}
		}
		for _, o := range z.Tiles[i].Objects {
			if t, ok := o.(Thinker); ok {
				z.Unlock()
				t.Think(z, x, y)
				z.Lock()
			}
		}
	}
	z.Unlock()
}

func (z *Zone) AllTileChange() []TileChange {
	z.Lock()
	defer z.Unlock()

	var changes []TileChange
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			t := z.Tile(uint8(x), uint8(y))
			if t != nil {
				for _, o := range t.Objects {
					changes = append(changes, TileChange{
						X:   uint8(x),
						Y:   uint8(y),
						ID:  o.NetworkID(),
						Obj: o.Serialize(),
					})
				}
			}
		}
	}

	return changes
}

func (z *Zone) Name() string {
	return "unknown zone"
}

type zindexsort []Object

func (z zindexsort) Len() int {
	return len(z)
}

func (z zindexsort) Swap(i, j int) {
	z[i], z[j] = z[j], z[i]
}

func (z zindexsort) Less(i, j int) bool {
	return z[i].ZIndex() > z[j].ZIndex()
}

type Tile struct {
	Objects []Object
	Sprite  uint8
}

func (t *Tile) Add(o Object) {
	t.Objects = append(t.Objects, o)
	sort.Sort(zindexsort(t.Objects))
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

/*func (t *Tile) Paint(z *Zone, i, j int, setcell func(int, int, PaintCell)) {
	if t.Sprite == 0 {
		t.Sprite = uint8(rand.Intn(4) + 1)
	}
	setcell(i, j, PaintCell{
		Sprite: "grass_r1",
		Color:  "#268f1e",
		SheetX: t.Sprite - 1,
		ZIndex: -50,
	})
	for k := len(t.Objects) - 1; k >= 0; k-- {
		t.Objects[k].Paint(i, j, setcell)
	}
}*/

type Object interface {
	Name() string
	Examine() string
	Blocking() bool
	ZIndex() int

	NetworkID() uint64
	Serialize() *NetworkedObject

	Interact(uint8, uint8, *Player, *Zone, int)
}

type Item interface {
	Object
	//	Mass() uint64   // grams
	//	Volume() uint64 // cubic centimeters
	AdminOnly() bool
}

type Thinker interface {
	Think(*Zone, uint8, uint8)
}

var _ Thinker = (*Player)(nil)
var _ Thinker = (*Hero)(nil)

func init() {
	gob.Register(Object(&Hero{}))
	// Players are removed from Zones before saving.
	gob.Register(Object(&Flora{}))
	gob.Register(Object(&Tree{}))
	gob.Register(Object(&Rock{}))
	gob.Register(Object(&WallStone{}))
	gob.Register(Object(&WallMetal{}))
	gob.Register(Object(&WallWood{}))
	gob.Register(Object(&FloorStone{}))
	gob.Register(Object(&FloorMetal{}))
	gob.Register(Object(&FloorWood{}))
	gob.Register(Object(&Liquid{}))
	gob.Register(Object(&Bed{}))
	gob.Register(Object(&Chest{}))

	gob.Register(Item(&Logs{}))
	gob.Register(Item(&Stone{}))
	gob.Register(Item(&Ore{}))
	gob.Register(Item(&Cosmetic{}))
	gob.Register(Item(&Pickaxe{}))
	gob.Register(Item(&Hatchet{}))
}
