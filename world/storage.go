package world

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
)

type savedZone struct {
	X, Y     int64
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
	z.X, z.Y = data.X, data.Y

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
	data.X, data.Y = z.X, z.Y
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
