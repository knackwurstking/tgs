package tgs

type ResponseGetMe struct {
	OK     bool `json:"ok"`
	Result User `json:"result"`
}

type ResponseGetUpdates struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// TODO: "setMyCommands"
// TODO: "sendMessage"
