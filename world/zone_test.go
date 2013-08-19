package world

import (
	"testing"
)

func TestZoneTile(t *testing.T) {
	var z Zone
	seen := make(map[[2]uint8]*Tile, 256*256)

	for i := 0; i < 10; i++ {
		for x := 0; x < 256; x++ {
			x8 := uint8(x)
			for y := 0; y < 256; y++ {
				y8 := uint8(y)
				tile := z.Tile(x8, y8)
				if tile == nil {
					t.Errorf("tile at (%d, %d) is nil on iteration %d", x, y, i)
				}
				if i != 0 {
					seenTile := seen[[2]uint8{x8, y8}]
					if seenTile != tile {
						t.Errorf("tile at (%d, %d) differs on iteration %i: %p != %p", x, y, i, tile, seenTile)
					}
				}
				seen[[2]uint8{x8, y8}] = tile
			}
		}
		for p, tile := range seen {
			if zone := tile.Zone(); zone != &z {
				t.Errorf("tile at (%d, %d) reports the wrong zone on iteration %d: %p != %p", p[0], p[1], i, zone, &z)
			}
			x, y := tile.Position()
			if x != p[0] || y != p[1] {
				t.Errorf("tile at (%d, %d) reports the wrong position on iteration %d: (%d, %d)", p[0], p[1], i, x, y)
			}
		}
	}
}
