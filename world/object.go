package world

type ObjectLike interface {
	// Name is the user-visible name of this object. Capitalization is only
	// used for living objects. Non-living object names are fully lowercase.
	Name() string

	// NotifyPosition is called when an object is added to, removed from,
	// or moved between tiles. The argument is the new tile, or nil for
	// removal. This function must return the previous value given to it.
	NotifyPosition(*Tile) *Tile

	// Think is called every tick (200ms) for objects that are in a zone.
	Think()
}

type Object struct {
	tile *Tile
}

func (o *Object) Name() string {
	return "unknown"
}

func (o *Object) NotifyPosition(t *Tile) *Tile {
	old := o.tile
	o.tile = t
	return old
}

func (o *Object) Think() {
	// do nothing
}
