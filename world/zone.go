package world

import (
	"sync"
)

type Zone struct {
	X, Y      int64 // unprotected by mutex; never changes
	tiles     [256 * 256]Tile
	listeners map[*ZoneListener]bool

	mtx sync.Mutex
}

type ZoneListener struct {
	Add    func(*Tile, ObjectLike)
	Remove func(*Tile, ObjectLike)
	Move   func(*Tile, *Tile, ObjectLike)
}

func (z *Zone) lock()   { z.mtx.Lock() }
func (z *Zone) unlock() { z.mtx.Unlock() }

func (z *Zone) AddListener(l *ZoneListener) {
	z.lock()
	defer z.unlock()

	if z.listeners == nil {
		z.listeners = make(map[*ZoneListener]bool)
	}
	z.listeners[l] = true
}

func (z *Zone) RemoveListener(l *ZoneListener) {
	z.lock()
	defer z.unlock()

	delete(z.listeners, l)
}

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
	for l := range z.listeners {
		z.unlock()
		if l.Add != nil {
			l.Add(t, obj)
		}
		z.lock()
	}
}

func (z *Zone) notifyRemove(t *Tile, obj ObjectLike) {
	old := obj.notifyPosition(nil)
	if old != t && old != nil {
		old.remove(obj)
	}
	for l := range z.listeners {
		z.unlock()
		if l.Remove != nil {
			l.Remove(t, obj)
		}
		z.lock()
	}
}

func (z *Zone) notifyMove(from *Tile, to *Tile, obj ObjectLike) {
	old := obj.notifyPosition(to)
	if old != from && old != nil {
		old.remove(obj)
	}
	for l := range z.listeners {
		z.unlock()
		if l.Move != nil {
			l.Move(from, to, obj)
		}
		z.lock()
	}
}
