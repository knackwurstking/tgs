package tgs

import "fmt"

type GetRequest interface {
	Token
	Get() (any, error)
}

type PostRequest interface {
	Token
	Post() (any, error)
}

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

func (r *RequestGetMe) Get() (*ResponseGetMe, error) {
	if r.token == "" {
		return nil, fmt.Errorf("missing token")
	}

	// TODO: ...

	return nil, fmt.Errorf("under construction")
}

type RequestGetUpdates struct {
	Token

	Offset         *int     `json:"offset"`
	Limit          *int     `json:"limit"`   // Limit defaults to 100
	Timeout        *int     `json:"timeout"` // Timeout defaults to 0 (Short Polling)
	AllowedUpdates []string `json:"allowed_updates"`
}

func (r *RequestGetUpdates) Get() (*ResponseGetUpdates, error) {
	if r.token == "" {
		return nil, fmt.Errorf("missing token")
	}

	// TODO: ...

	return nil, fmt.Errorf("under construction")
}

// TODO: "setMyCommands"
// TODO: "sendMessage"
