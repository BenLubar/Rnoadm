package world

import (
	"math/big"
	"strings"
	"sync"
)

type Zone struct {
	X, Y      int64 // unprotected by mutex; never changes
	tiles     [256 * 256]Tile
	listeners map[*ZoneListener]bool

	impersonation map[PlayerLike]Visible
	impersonated  map[Visible]PlayerLike

	mtx sync.Mutex
}

type ZoneListener struct {
	Add    func(*Tile, ObjectLike)
	Remove func(*Tile, ObjectLike)
	Move   func(*Tile, *Tile, ObjectLike)
	Update func(*Tile, ObjectLike)
	Damage func(Combat, Combat, *big.Int)
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
	if l.Add != nil {
		for i := range z.tiles {
			t := &z.tiles[i]
			for _, o := range t.objects {
				z.unlock()
				l.Add(t, o)
				z.lock()
			}
		}
	}
}

func (z *Zone) RemoveListener(l *ZoneListener) {
	z.lock()
	defer z.unlock()

	delete(z.listeners, l)
	if l.Remove != nil {
		for i := range z.tiles {
			t := &z.tiles[i]
			for _, o := range t.objects {
				z.unlock()
				l.Remove(t, o)
				z.lock()
			}
		}
	}
}

// Tile returns the Tile at a position in the Zone.
func (z *Zone) Tile(x, y uint8) *Tile {
	z.lock()
	defer z.unlock()

	return z.tile(x, y)
}

func (z *Zone) tile(x, y uint8) *Tile {
	t := &z.tiles[uint(x)|uint(y)<<8]
	if t.zone == nil {
		t.x, t.y, t.zone = x, y, z
	}

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

func (z *Zone) Update(t *Tile, obj ObjectLike) {
	z.lock()
	defer z.unlock()

	if visible, ok := obj.(Visible); ok {
		if p, ok := z.impersonated[visible]; ok {
			for l := range z.listeners {
				z.unlock()
				if l.Update != nil {
					l.Update(t, p)
				}
				z.lock()
			}
			return
		}
	}

	for l := range z.listeners {
		z.unlock()
		if l.Update != nil {
			l.Update(t, obj)
		}
		z.lock()
	}
}

func (z *Zone) Damage(attacker, victim Combat, amount *big.Int) {
	z.lock()
	for l := range z.listeners {
		z.unlock()
		if l.Damage != nil {
			l.Damage(attacker, victim, amount)
		}
		z.lock()
	}
	z.unlock()
}

func (z *Zone) Impersonate(player PlayerLike, o Visible) {
	z.lock()
	defer z.unlock()

	if z.impersonation == nil {
		z.impersonation = make(map[PlayerLike]Visible)
		z.impersonated = make(map[Visible]PlayerLike)
	}

	if old, ok := z.impersonation[player]; ok {
		delete(z.impersonation, player)
		delete(z.impersonated, old)
	}

	if o == nil {
		return
	}

	z.impersonation[player] = o
	z.impersonated[o] = player
	o.notifyPosition(player.Position())
}

func (z *Zone) think(wg *sync.WaitGroup) {
	defer wg.Done()

	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			z.Tile(uint8(x), uint8(y)).think()
		}
	}
}

func (z *Zone) Chat(sender Visible, message string) {
	message = strings.Join(strings.Fields(message), " ")
	if message == "" {
		return
	}
	message = "‹" + sender.Name() + "› says, “" + message + "”"
	color := "#ccc"
	if a, ok := sender.(AdminLike); ok {
		if a.IsAdmin() {
			color = "#fd8"
		} else {
			color = "#eee"
		}
	}

	z.lock()
	defer z.unlock()

	for i := range z.tiles {
		for _, o := range z.tiles[i].objects {
			if recipient, ok := o.(SendMessageLike); ok {
				recipient.SendMessageColor(message, color)
			}
		}
	}
}
