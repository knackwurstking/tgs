package tgs

const (
	CommandGetMe         = Command("getMe")
	CommandGetUpdates    = Command("getUpdates")
	CommandSetMyCommands = Command("setMyCommands")
	CommandSendMessage   = Command("sendMessage")
)

type Command string
