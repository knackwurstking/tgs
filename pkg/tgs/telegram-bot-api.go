package tgs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TelegramBotAPI struct {
	token string
}

func NewTelegramBotAPI(token string) *TelegramBotAPI {
	return &TelegramBotAPI{
		token: token,
	}
}

func (this *TelegramBotAPI) Token() string {
	return this.token
}

func (this *TelegramBotAPI) SetToken(token string) {
	this.token = token
}

func (this *TelegramBotAPI) URL(command Command) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", this.Token(), command)
}

func (this *TelegramBotAPI) SendRequest(request Request) ([]byte, error) {
	if this.Token() == "" {
		return nil, fmt.Errorf("missing token")
	}

	if request == nil {
		return nil, fmt.Errorf("missing request")
	}

	var (
		bodyData []byte
		method   string
	)

	if request.Body() != nil {
		var err error
		bodyData, err = json.Marshal(request.Body())
		if err != nil {
			return nil, err
		}
	}

	switch request.Command() {
	case CommandGetMe, CommandGetUpdates:
		method = "GET"
		break
	case CommandSetMyCommands, CommandSendMessage:
		method = "POST"
		break
	default:
		return nil, fmt.Errorf(fmt.Sprintf("unknown command: %s", request.Command()))
	}

	body := bytes.NewBuffer(bodyData)
	r, err := http.NewRequest(method, this.URL(request.Command()), body)
	if err != nil {
		return nil, err
	}

	if request.Body() != nil {
		r.Header.Add("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"telegram request to \"%s\" returned status code %d (%s)",
			r.URL, resp.StatusCode, resp.Status,
		)
	}

	return io.ReadAll(resp.Body)
}
