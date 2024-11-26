package botcommands

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IP struct {
	*tgbotapi.BotAPI
}

func NewIP(botAPI *tgbotapi.BotAPI) *IP {
	return &IP{
		BotAPI: botAPI,
	}
}

func (*IP) URL() string {
	return "https://ifconfig.io"
}

func (this *IP) Run(chatID int64, replyID *int) error {
	address, err := this.FetchAddressFromURL()
	if err != nil {
		return err
	}

	msgConfig := tgbotapi.NewMessage(chatID, fmt.Sprintf("`%s`", address))
	msgConfig.ReplyToMessageID = *replyID
	msgConfig.ParseMode = "MarkdownV2"

	_, err = this.BotAPI.Send(msgConfig)
	return err
}

func (this *IP) FetchAddressFromURL() (address string, err error) {
	resp, err := http.Get(this.URL())
	if err != nil {
		return address, err
	}
	if resp.StatusCode != http.StatusOK {
		return address, fmt.Errorf("request to %s: %d (%s)",
			this.URL(), resp.StatusCode, resp.Status,
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return address, err
	}

	return strings.Trim(string(data), "\n\r\t "), nil
}
