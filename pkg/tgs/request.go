package tgs

type RequestGetMe struct{}

type RequestGetUpdates struct {
	Offset *int `json:"offset"`
	// Limit defaults to 100
	Limit *int `json:"limit"`
	// Timeout defaults to 0 (Short Polling)
	Timeout        *int     `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

// TODO: "setMyCommands"
// TODO: "sendMessage"
