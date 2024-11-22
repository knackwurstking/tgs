package tgs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Token struct {
	token string
}

func (t *Token) Token() string {
	return t.token
}

func (t *Token) SetToken(token string) {
	t.token = token
}

type RequestGetMe struct {
	Token
}

// TODO: Remove... Moved to "./tgs.go"
func (r *RequestGetMe) Get() (*ResponseGetMe, error) {
	if r.token == "" {
		return nil, fmt.Errorf("missing token")
	}

	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(data)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", r.token)
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var getMe ResponseGetMe
	err = json.Unmarshal(bodyData, &getMe)
	return &getMe, err
}

type RequestGetUpdates struct {
	Token

	Offset         *int     `json:"offset"`
	Limit          *int     `json:"limit"`   // Limit defaults to 100
	Timeout        *int     `json:"timeout"` // Timeout defaults to 0 (Short Polling)
	AllowedUpdates []string `json:"allowed_updates"`
}

// TODO: Remove... Moved to "./tgs.go"
func (r *RequestGetUpdates) Get() (*ResponseGetUpdates, error) {
	if r.token == "" {
		return nil, fmt.Errorf("missing token")
	}

	// TODO: ...

	return nil, fmt.Errorf("under construction")
}

// TODO: "setMyCommands"
// TODO: "sendMessage"
