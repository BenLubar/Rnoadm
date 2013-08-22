package world

type SendMessageLike interface {
	SendMessage(string)
	SendMessageColor(string, string)
}

type AdminLike interface {
	IsAdmin() bool
}

type InventoryLike interface {
	Inventory() []Visible
	GiveItem(Visible) bool
	RemoveItem(Visible) bool
}

type CombatInventoryMessageAdmin interface {
	Combat
	InventoryLike
	SendMessageLike
	AdminLike
}

type Item interface {
	Volume() uint64
	Weight() uint64
	AdminOnly() bool
}
