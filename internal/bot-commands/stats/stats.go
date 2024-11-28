package stats

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func New(botAPI *tgbotapi.BotAPI) *Stats {
	return NewStats(botAPI)
}

type Stats struct {
	*tgbotapi.BotAPI
}

func NewStats(botAPI *tgbotapi.BotAPI) *Stats {
	return &Stats{
		BotAPI: botAPI,
	}
}

// TODO: Deliver user stats like user id (+username), chat id (+title), topic number (message_thread_id, +reply_to_message.forum_topic_created.name)

func (this *Stats) Run(message *tgbotapi.Message) error {
	// TODO: Continue here... Create the message and reply (send) the message back

	return nil
}
