package main

import (
	"testing"
)

func TestZoneTilesUnique(t *testing.T) {
	seen := make(map[*Tile][2]uint8, zoneTiles)
	zone := &Zone{X: 0, Y: 0}
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			tile := zone.Tile(uint8(x), uint8(y))
			if tile == nil {
			} else if pos, ok := seen[tile]; ok {
				t.Errorf("Tile at (%d, %d) is the same as the tile at (%d, %d)", x, y, pos[0], pos[1])
			} else {
				seen[tile] = [2]uint8{uint8(x), uint8(y)}
			}
		}
	}
}
