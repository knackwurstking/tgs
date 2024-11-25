package tgs

const (
	CommandGetMe       = Command("getMe")
	CommandGetUpdates  = Command("getUpdates")
	CommandSendMessage = Command("sendMessage")
)

type Command string
