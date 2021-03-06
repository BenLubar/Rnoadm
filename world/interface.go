package world

import (
	"time"
)

type InventoryLike interface {
	Inventory() []Visible
	GiveItem(Visible) bool
	RemoveItem(Visible) bool
}

type SendMessageLike interface {
	SendMessage(string)
	SendMessageColor(string, string)
}

type AdminLike interface {
	IsAdmin() bool
	Impersonate(Visible)
}

type SetHUDLike interface {
	SetHUD(string, map[string]interface{})
}

type PlayerLike interface {
	Combat
	InventoryLike
	SendMessageLike
	AdminLike
	SetHUDLike
	Instance(*Tile) Instance
	Weight() uint64
	WeightMax() uint64
	Volume() uint64
	VolumeMax() uint64
}

type Instance interface {
	Items(func([]Visible) []Visible)
	Last(func(time.Time) time.Time)
}

type Item interface {
	Volume() uint64
	Weight() uint64
	AdminOnly() bool
}
