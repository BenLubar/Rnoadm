package world

type SendMessageLike interface {
	SendMessage(string)
	SendMessageColor(string, string)
}

type AdminLike interface {
	IsAdmin() bool
}
