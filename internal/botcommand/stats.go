package botcommand

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Stats struct {
	*tgbotapi.BotAPI
}

func NewStats(botAPI *tgbotapi.BotAPI) *Stats {
	return &Stats{
		BotAPI: botAPI,
	}
}

func (this *Stats) Run(message *tgbotapi.Message) error {
	data := struct {
		UserName        string `json:"username"`
		UserID          int64  `json:"user_id"`
		ChatID          int64  `json:"chat_id"`
		MessageThreadID int    `json:"message_thread_id"`
	}{
		UserName:        message.From.UserName,
		UserID:          message.From.ID,
		ChatID:          message.Chat.ID,
		MessageThreadID: message.MessageThreadID,
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")

	msgConfig := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf("```json\n")+
			fmt.Sprintf("%s\n", string(jsonData))+
			fmt.Sprintf("```"),
	)

	msgConfig.ReplyToMessageID = message.MessageID
	msgConfig.ParseMode = "MarkdownV2"

	_, err = this.BotAPI.Send(msgConfig)
	return err
}
