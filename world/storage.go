package world

import (
	"bytes"
	"compress/gzip"
	"encoding/base32"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type savedZone struct {
	X, Y     int64
	Z        int8
	Version  uint
	TileData []savedTile
}

type savedTile struct {
	Version uint
	Data    interface{}
}

func readZone(r io.Reader) *Zone {
	// this function panics on errors as all encoding errors are bugs and
	// need to be manually resolved in the code.

	g, err := gzip.NewReader(r)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	var data savedZone
	err = gob.NewDecoder(g).Decode(&data)
	if err != nil {
		panic(err)
	}

	var z Zone
	z.lock()
	defer z.unlock()
	z.X, z.Y, z.Z = data.X, data.Y, data.Z

	switch data.Version {
	case 0:
		if len(data.TileData) != 256*256 {
			panic(fmt.Sprintf("wrong size for TileData: %d != %d", len(data.TileData), 256*256))
		}
		i := 0
		for x := 0; x < 256; x++ {
			x8 := uint8(x)
			for y := 0; y < 256; y++ {
				y8 := uint8(y)
				t := z.tile(x8, y8)
				t.load(data.TileData[i].Version, data.TileData[i].Data)
				i++
			}
		}
	default:
		panic(fmt.Sprintf("version %d unknown", data.Version))
	}

	return &z
}

func writeZone(z *Zone, w io.Writer) {
	// this function panics on errors as all encoding errors are bugs and
	// need to be manually resolved in the code.

	z.lock()
	defer z.unlock()

	var data savedZone
	data.X, data.Y, data.Z = z.X, z.Y, z.Z
	data.Version = 0
	data.TileData = make([]savedTile, 256*256)
	i := 0
	for x := 0; x < 256; x++ {
		x8 := uint8(x)
		for y := 0; y < 256; y++ {
			y8 := uint8(y)
			t := z.tile(x8, y8)
			data.TileData[i].Version, data.TileData[i].Data = t.save()
			i++
		}
	}

	g, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	err = gob.NewEncoder(g).Encode(&data)
	if err != nil {
		panic(err)
	}
}

type zoneCoord struct {
	X, Y int64
	Z    int8
}

type loadedZone struct {
	zone *Zone
	ref  uint
}

var loadedZones = make(map[zoneCoord]*loadedZone)
var loadedZoneLock sync.Mutex

// run zone loading and saving in a separate goroutine - corrupt zones cause
// full server crashes to avoid further corruption.
func zoneCritical(f func()) {
	ch := make(chan struct{}, 1)

	go func() {
		loadedZoneLock.Lock()
		defer loadedZoneLock.Unlock()

		f()

		ch <- struct{}{}
	}()

	<-ch
}

func zoneFilename(x, y int64, z int8) string {
	var buf [binary.MaxVarintLen64*2 + 1]byte
	i := binary.PutVarint(buf[:], x)
	i += binary.PutVarint(buf[i:], y)
	buf[i] = uint8(z)
	i++
	encoded := base32.StdEncoding.EncodeToString(buf[:i])
	return filepath.Join("rnoadm-AA", "zone"+encoded+".gz")
}

func zoneFilenameV1(x, y int64) string {
	var buf [binary.MaxVarintLen64 * 2]byte
	i := binary.PutVarint(buf[:], x)
	i += binary.PutVarint(buf[i:], y)
	encoded := base32.StdEncoding.EncodeToString(buf[:i])
	for i := range encoded {
		if encoded[i] == '=' {
			encoded = encoded[:i]
			break
		}
	}
	return filepath.Join("rnoadm-AA", "zone"+encoded+".gz")
}

func GetZone(x, y int64, z int8) *Zone {
	var zone *Zone
	zoneCritical(func() {
		lz, ok := loadedZones[zoneCoord{x, y, z}]
		if ok {
			lz.ref++
			zone = lz.zone
			return
		}

		f, err := os.Open(zoneFilename(x, y, z))
		if os.IsNotExist(err) && z == 0 {
			if os.Rename(zoneFilenameV1(x, y), zoneFilename(x, y, z)) == nil {
				f, err = os.Open(zoneFilename(x, y, z))
			}
		}
		if err == nil {
			defer f.Close()
			zone = readZone(f)
		} else {
			if !os.IsNotExist(err) {
				panic(err)
			}
			zone = generateZone(x, y, z)
		}

		lz = &loadedZone{
			zone: zone,
			ref:  1,
		}
		loadedZones[zoneCoord{x, y, z}] = lz
	})
	return zone
}

func ReleaseZone(z *Zone) {
	zoneCritical(func() {
		lz := loadedZones[zoneCoord{z.X, z.Y, z.Z}]
		lz.ref--
		if lz.ref != 0 {
			return
		}

		delete(loadedZones, zoneCoord{z.X, z.Y, z.Z})

		// write to a memory buffer first to avoid corruption on failed
		// saves
		var buf bytes.Buffer
		writeZone(lz.zone, &buf)

		f, err := os.Create(zoneFilename(z.X, z.Y, z.Z))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = buf.WriteTo(f)
		if err != nil {
			panic(err)
		}
	})
}

func Think() {
	var wg sync.WaitGroup
	zoneCritical(func() {
		for _, lz := range loadedZones {
			wg.Add(1)
			go lz.zone.think(&wg)
		}
	})
	wg.Wait()
}

func SaveAllZones() {
	zoneCritical(func() {
		for _, lz := range loadedZones {
			var buf bytes.Buffer
			writeZone(lz.zone, &buf)

			f, err := os.Create(zoneFilename(lz.zone.X, lz.zone.Y, lz.zone.Z))
			if err != nil {
				panic(err)
			}
			_, err = buf.WriteTo(f)
			f.Close()
			if err != nil {
				panic(err)
			}
		}
	})
}
