package world

func (z *Zone) Path(start *Tile, end *Tile, stopEarly bool) [][2]uint8 {
	queue := []*Tile{start}
	from := map[*Tile]*Tile{start: nil}
	closest := start
	sx, sy := start.Position()
	ex, ey := end.Position()
	distance := (int(sx)-int(ex))*(int(sx)-int(ex)) + (int(sy)-int(ey))*(int(sy)-int(ey))

	for len(queue) != 0 {
		t := queue[0]
		queue = queue[1:]
		if t.Blocked() {
			continue
		}
		if t == end {
			if stopEarly {
				return z.constructPath(from, from[end])
			}
			return z.constructPath(from, end)
		}
		x, y := t.Position()
		d := (int(x)-int(ex))*(int(x)-int(ex)) + (int(y)-int(ey))*(int(y)-int(ey))
		if d < distance {
			closest = t
			distance = d
		}
		if x > 0 {
			next := t.Zone().Tile(x-1, y)
			if _, ok := from[next]; !ok {
				from[next] = t
				queue = append(queue, next)
			}
		}
		if x < 255 {
			next := t.Zone().Tile(x+1, y)
			if _, ok := from[next]; !ok {
				from[next] = t
				queue = append(queue, next)
			}
		}
		if y > 0 {
			next := t.Zone().Tile(x, y-1)
			if _, ok := from[next]; !ok {
				from[next] = t
				queue = append(queue, next)
			}
		}
		if y < 255 {
			next := t.Zone().Tile(x, y+1)
			if _, ok := from[next]; !ok {
				from[next] = t
				queue = append(queue, next)
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
