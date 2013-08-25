package world

import (
	"math/rand"
)

func (z *Zone) Path(start *Tile, end *Tile, stopEarly bool) [][2]uint8 {
	queue := []*Tile{start}
	from := map[*Tile]*Tile{start: nil}
	closest := start
	ex, ey := end.Position()
	distance := 1 << 30

	for len(queue) != 0 {
		t := queue[0]
		queue = queue[1:]
		if t != start && t.Blocked() {
			continue
		}
		if t == end && (!stopEarly || start != end) {
			if stopEarly {
				return z.constructPath(from, from[end])
			}
			return z.constructPath(from, end)
		}
		x, y := t.Position()
		d := (int(x)-int(ex))*(int(x)-int(ex)) + (int(y)-int(ey))*(int(y)-int(ey))
		if d < distance && (!stopEarly || t != start) {
			closest = t
			distance = d
		}
		for _, i := range rand.Perm(4) {
			switch i {
			case 0:
				if x > 0 {
					next := t.Zone().Tile(x-1, y)
					if _, ok := from[next]; !ok {
						from[next] = t
						queue = append(queue, next)
					}
				}
			case 1:
				if x < 255 {
					next := t.Zone().Tile(x+1, y)
					if _, ok := from[next]; !ok {
						from[next] = t
						queue = append(queue, next)
					}
				}
			case 2:
				if y > 0 {
					next := t.Zone().Tile(x, y-1)
					if _, ok := from[next]; !ok {
						from[next] = t
						queue = append(queue, next)
					}
				}
			case 3:
				if y < 255 {
					next := t.Zone().Tile(x, y+1)
					if _, ok := from[next]; !ok {
						from[next] = t
						queue = append(queue, next)
					}
				}
			}
		}
	}

	return z.constructPath(from, closest)
}

func (z *Zone) constructPath(from map[*Tile]*Tile, end *Tile) [][2]uint8 {
	var path [][2]uint8

	for end != nil {
		x, y := end.Position()
		path = append(path, [2]uint8{x, y})
		end = from[end]
	}

	// reverse the slice
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}
