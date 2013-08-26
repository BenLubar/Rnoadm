package world

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

type CombatInventoryMessageAdminHUD interface {
	Combat
	InventoryLike
	SendMessageLike
	AdminLike
	SetHUDLike
}

type Item interface {
	Volume() uint64
	Weight() uint64
	AdminOnly() bool
}
