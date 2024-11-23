package tgs

const (
	CommandGetMe      = Command("getMe")
	CommandGetUpdates = Command("getUpdates")
)

type Command string
