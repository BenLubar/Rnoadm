package world

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type ObjectLike interface {
	// NotifyPosition is called when an object is added to, removed from,
	// or moved between tiles. The argument is the new tile, or nil for
	// removal. This function must return the previous value given to it.
	NotifyPosition(*Tile) *Tile

	// This function is not called, but is provided in the default
	// implementation for convenience.
	Position() *Tile

	// Think is called every tick (200ms) for objects that are in a zone.
	Think()

	// Save returns a version number and data describing the object.
	// Data may only contain primitives, []interface{}, and
	// map[string]interface{}, recursively.
	Save() (uint, interface{})

	// Load is given a version number and data from a previous call to Save.
	// Panicing on errors is suggested. From a zero value of a type, after
	// a call to Load, the NotifyPosition method's next call should return
	// nil.
	Load(uint, interface{})
}

type Object struct {
	tile unsafe.Pointer // *Tile
}

func init() {
	Register("object", (*Object)(nil))
}

func (o *Object) NotifyPosition(t *Tile) *Tile {
	for {
		old := atomic.LoadPointer(&o.tile)
		if atomic.CompareAndSwapPointer(&o.tile, old, unsafe.Pointer(t)) {
			return (*Tile)(old)
		}
	}
}

func (o *Object) Position() *Tile {
	return (*Tile)(atomic.LoadPointer(&o.tile))
}

func (o *Object) Think() {
	// do nothing
}

func (o *Object) Save() (uint, interface{}) {
	return 0, 0
}

func (o *Object) Load(version uint, data interface{}) {
	switch version {
	case 0:
		// no fields in version 0
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}
