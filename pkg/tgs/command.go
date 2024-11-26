package tgs

const (
	CommandGetMe         = Command("getMe")
	CommandGetUpdates    = Command("getUpdates")
	CommandSetMyCommands = Command("setMyCommands") // TODO: Add get and delete my commands
	CommandSendMessage   = Command("sendMessage")
)

type Command string
