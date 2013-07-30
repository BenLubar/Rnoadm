package main

import (
	"math"
	"math/rand"
)

func fuzzyDistance(dx, dy uint8, radius int, r *rand.Rand) bool {
	a := int(int8(dx))*int(int8(dx)) + int(int8(dy))*int(int8(dy))
	radius += r.Intn(radius/2 + 1)
	b := radius * radius

	return a < b
}

func (z *Zone) generateForest(r *rand.Rand) {
	z.Name_ = GenerateName(r, NameZone, NameForest)

	var nodes struct {
		main [5]struct {
			x, y uint8
		}
	}

	for i := range nodes.main {
		radius, theta := r.Float64()*100, r.Float64()*2*math.Pi
		nodes.main[i].x = uint8(127.5 + math.Cos(theta)*radius)
		nodes.main[i].y = uint8(127.5 + math.Sin(theta)*radius)
	}

	for x := uint8(1); x < 255; x++ {
		for y := uint8(1); y < 255; y++ {
			if x < zoneOffset[y]+1 || x >= 255-zoneOffset[y]-1 {
				continue
			}
			tile := z.Tile(x, y)
			tile.Add(&Tree{Type: Oak})
		}
	}

	for i := range nodes.main {
		for x := nodes.main[i].x - 25; x < nodes.main[i].x+25; x++ {
			for y := nodes.main[i].y - 25; y < nodes.main[i].y+25; y++ {
				tile := z.Tile(x, y)
				if fuzzyDistance(x-nodes.main[i].x, y-nodes.main[i].y, 20, r) {
					for j := 0; j < len(tile.Objects); j++ {
						if _, ok := tile.Objects[j].(*Tree); ok {
							tile.Objects = append(tile.Objects[:j], tile.Objects[j+1:]...)
						}
					}
				}
				if fuzzyDistance(x-nodes.main[i].x, y-nodes.main[i].y, 10, r) {
					tile.Add(&Liquid{})
				}
			}
		}
	}
}
