package tgs

type RequestGetMe struct{}

type RequestGetUpdates struct {
	Offset         *int     `json:"offset"`
	Limit          *int     `json:"limit"`
	Timeout        *int     `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

// TODO: "setMyCommands"
// TODO: "sendMessage"
