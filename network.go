package main

import (
	"sync/atomic"
)

var nextNetworkID uint64

// networkID implements the NetworkID() function for Object in a safe, atomic
// way. During the execution of this program, each networkID is guaranteed to
// return a unique ID from the aforementioned function.
type networkID uint64

func (id *networkID) NetworkID() uint64 {
	if i := atomic.LoadUint64((*uint64)(id)); i != 0 {
		// simple case: we already have a network ID; return it
		return i
	}
	// set our network ID to the next available ID, but do nothing if
	// we already have an ID set.
	atomic.CompareAndSwapUint64((*uint64)(id), 0, atomic.AddUint64(&nextNetworkID, 1))
	// we definitely have a network ID at this point; return it.
	return atomic.LoadUint64((*uint64)(id))
}

func (id *networkID) Serialize() *NetworkedObject {
	// TODO: remove this method
	return &NetworkedObject{
		Sprite: "ui_r1",
		Colors: []Color{"#f0f"},
	}
}

// stub function for non-interactable Objects.
func (id *networkID) Interact(x, y uint8, player *Player, zone *Zone, opt int) {}
