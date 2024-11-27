package tgs

const (
	CommandGetMe            = Command("getMe")
	CommandGetUpdates       = Command("getUpdates")
	CommandSetMyCommands    = Command("setMyCommands")
	CommandDeleteMyCommands = Command("deleteMyCommands")
	CommandSendMessage      = Command("sendMessage")
)

type Command string
