package world

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// ObjectLike is the base interface which all persistent game objects extend.
type ObjectLike interface {
	notifyOuter(ObjectLike)
	notifyPosition(*Tile) *Tile

	// NotifyPosition is called whenever the tile containing the object
	// changes. Either argument may be nil.
	NotifyPosition(old, new *Tile)

	// Outer returns the object given as an argument to InitObject.
	Outer() ObjectLike

	// Position returns the tile containing this object.
	Position() *Tile

	// Think is called every tick (200ms) for objects that are in a zone.
	Think()

	// Save returns a serialized version of this object. The data may
	// contain primitives, []interface{}, map[string]interface{}, time.Time,
	// and *big.Int. Sub-objects and objects owned by this object should
	// be returned in the third argument.
	Save() (version uint, data interface{}, attached []ObjectLike)

	// Load unserializes the data returned from a previous call to Save.
	// On unrecoverable errors, it is recommended to panic.
	Load(version uint, data interface{}, attached []ObjectLike)
}

// Object implements ObjectLike.
type Object struct {
	outer ObjectLike
	tile  unsafe.Pointer // *Tile
}

func init() {
	Register("object", ObjectLike((*Object)(nil)))
}

// InitObject sets the value returned by ObjectLike.Outer. Objects should only
// be initialized before they are inserted into the world.
func InitObject(obj ObjectLike) ObjectLike {
	obj.notifyOuter(obj)
	return obj
}

func (o *Object) notifyOuter(outer ObjectLike) {
	o.outer = outer
}

// Outer returns the object given as an argument to InitObject.
func (o *Object) Outer() ObjectLike {
	return o.outer
}

func (o *Object) notifyPosition(t *Tile) *Tile {
	if p, ok := o.outer.(PlayerLike); t != nil && ok {
		if i, ok := t.Zone().impersonation[p]; ok {
			i.notifyPosition(t)
		}
	}
	for {
		old := atomic.LoadPointer(&o.tile)
		if atomic.CompareAndSwapPointer(&o.tile, old, unsafe.Pointer(t)) {
			go o.outer.NotifyPosition((*Tile)(old), t)
			return (*Tile)(old)
		}
	}
}

// NotifyPosition is called whenever the tile containing the object
// changes. Either argument may be nil.
func (o *Object) NotifyPosition(old, new *Tile) {
	// do nothing
}

// Position returns the tile containing this object.
func (o *Object) Position() *Tile {
	return (*Tile)(atomic.LoadPointer(&o.tile))
}

// Think is called every tick (200ms) for objects that are in a zone.
func (o *Object) Think() {
	// do nothing
}

// Save returns a serialized version of this object. The data may
// contain primitives, []interface{}, map[string]interface{}, time.Time,
// and *big.Int. Sub-objects and objects owned by this object should
// be returned in the third argument.
func (o *Object) Save() (uint, interface{}, []ObjectLike) {
	return 0, uint(0), nil
}

// Load unserializes the data returned from a previous call to Save.
// On unrecoverable errors, it is recommended to panic.
func (o *Object) Load(version uint, data interface{}, attached []ObjectLike) {
	switch version {
	case 0:
		// no fields in version 0
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}
