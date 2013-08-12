package world

import (
	"fmt"
)

type Tile struct {
	objects []ObjectLike

	x, y uint8
	zone *Zone
}

func (t *Tile) Objects() []ObjectLike {
	t.zone.lock()
	defer t.zone.unlock()

	objects := make([]ObjectLike, len(t.objects))
	copy(objects, t.objects)

	return objects
}

func (t *Tile) Add(obj ObjectLike) {
	t.zone.lock()
	defer t.zone.unlock()

	if t.add(obj) {
		t.zone.notifyAdd(t, obj)
	}
}

func (t *Tile) add(obj ObjectLike) bool {
	for _, o := range t.objects {
		if o == obj {
			return false
		}
	}

	t.objects = append(t.objects, obj)
	return true
}

// Remove an object from this tile. Returns true if the object was found.
func (t *Tile) Remove(obj ObjectLike) bool {
	t.zone.lock()
	defer t.zone.unlock()

	if t.remove(obj) {
		t.zone.notifyRemove(t, obj)
		return true
	}
	return false
}

func (t *Tile) remove(obj ObjectLike) bool {
	for i, o := range t.objects {
		if o == obj {
			t.objects = append(t.objects[:i], t.objects[i+1:]...)
			return true
		}
	}
	return false
}

// Move an object from this tile to another tile. Cross-zone movement is not
// performed by this function. Instead, Remove the object from this tile and
// Add it to the other.
func (t *Tile) Move(obj ObjectLike, to *Tile) {
	t.zone.lock()
	defer t.zone.unlock()

	if t == to || t.zone != to.zone {
		return
	}

	if t.remove(obj) && to.add(obj) {
		t.zone.notifyMove(t, to, obj)
	}
}

func (t *Tile) Position() (x, y uint8) {
	// this is safe - Tiles never move between locations or zones.
	return t.x, t.y
}

func (t *Tile) Zone() *Zone {
	// this is safe - Tiles never move between locations or zones.
	return t.zone
}

func (t *Tile) save() (uint, interface{}) {
	contents := make([]interface{}, 0, len(t.objects))

	for _, o := range t.objects {
		version, data := o.Save()
		contents = append(contents, map[string]interface{}{
			"t": getObjectTypeIdentifier(o),
			"v": version,
			"d": data,
		})
	}

	return 0, map[string]interface{}{
		"c": contents,
	}
}

func (t *Tile) load(version uint, data interface{}) {
	switch version {
	case 0:
		contents := data.(map[string]interface{})["c"].([]interface{})
		t.objects = make([]ObjectLike, 0, len(contents))
		for _, c := range contents {
			o := c.(map[string]interface{})
			obj := getObjectByIdentifier(o["t"].(string))
			obj.Load(o["v"].(uint), o["d"])
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
	for _, obj := range t.objects {
		obj.NotifyPosition(t)
	}
}
