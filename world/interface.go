package world

type SendMessageLike interface {
	SendMessage(string)
	SendMessageColor(string, string)
}

type AdminLike interface {
	IsAdmin() bool
}

type Item interface {
	Volume() uint64
	Weight() uint64
	AdminOnly() bool
}
