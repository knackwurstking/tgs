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

func NewTelegramBotAPI(token string) API {
	return &TelegramBotAPI{
		token: token,
	}
}

func (api *TelegramBotAPI) Token() string {
	return api.token
}

func (api *TelegramBotAPI) SetToken(token string) {
	api.token = token
}

func (api *TelegramBotAPI) URL(command Command) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", api.Token(), command)
}

func (api *TelegramBotAPI) Send(request Request) ([]byte, error) {
	if api.Token() == "" {
		return nil, fmt.Errorf("missing token")
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var method string
	switch request.Command() {
	case CommandGetMe, CommandGetUpdates:
		method = "GET"
		break
	default:
		return nil, fmt.Errorf(fmt.Sprintf("unknown command: %s", request.Command()))
	}

	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(method, api.URL(request.Command()), body)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}
