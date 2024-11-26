package commands

type TelegramCommandHandler interface {
	Run(chatID int) error
}
