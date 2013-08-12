package world

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type ObjectLike interface {
	notifyOuter(ObjectLike)
	notifyPosition(*Tile) *Tile

	Outer() ObjectLike

	// Position returns the tile this object is located in.
	Position() *Tile

	// Think is called every tick (200ms) for objects that are in a zone.
	Think()

	// Save returns a version number, data, and attached objects describing
	// the object. Data may only contain primitives, []interface{}, and
	// map[string]interface{}, recursively. If this object contains other
	// objects, return them in the third return argument.
	Save() (uint, interface{}, []ObjectLike)

	// Load is given a version number, data, and attached objects from a
	// previous call to Save. Panicking on errors is suggested.
	Load(uint, interface{}, []ObjectLike)
}

type Object struct {
	outer ObjectLike
	tile  unsafe.Pointer // *Tile
}

func init() {
	Register("object", (*Object)(nil))
}

func InitObject(obj ObjectLike) ObjectLike {
	obj.notifyOuter(obj)
	return obj
}

func (o *Object) notifyOuter(outer ObjectLike) {
	o.outer = outer
}

func (o *Object) Outer() ObjectLike {
	return o.outer
}

func (o *Object) notifyPosition(t *Tile) *Tile {
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

func (o *Object) Save() (uint, interface{}, []ObjectLike) {
	return 0, 0, nil
}

func (o *Object) Load(version uint, data interface{}, attached []ObjectLike) {
	switch version {
	case 0:
		// no fields in version 0
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}
