package world

import (
	"sync"
)

type Zone struct {
	X, Y  int64           // unprotected by mutex; never changes
	tiles [256 * 256]Tile // individual tiles are protected by the mutex

	mtx sync.Mutex
}

func (z *Zone) lock()   { z.mtx.Lock() }
func (z *Zone) unlock() { z.mtx.Unlock() }

// Tile returns the Tile at a position in the Zone.
func (z *Zone) Tile(x, y uint8) *Tile {
	z.lock()
	defer z.unlock()

	return z.tile(x, y)
}

func (z *Zone) tile(x, y uint8) *Tile {
	t := &z.tiles[uint(x)|uint(y)<<8]
	t.x, t.y, t.zone = x, y, z

	return t
}

func (z *Zone) notifyAdd(t *Tile, obj ObjectLike) {
	old := obj.notifyPosition(t)
	if old != nil {
		old.remove(obj)
	}
}

func (z *Zone) notifyRemove(t *Tile, obj ObjectLike) {
	old := obj.notifyPosition(nil)
	if old != t && old != nil {
		old.remove(obj)
	}
}

func (z *Zone) notifyMove(from *Tile, to *Tile, obj ObjectLike) {
	old := obj.notifyPosition(to)
	if old != from && old != nil {
		old.remove(obj)
	}
}
