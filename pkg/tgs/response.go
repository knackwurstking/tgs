package tgs

type ResponseGetMe struct {
	OK     bool `json:"ok" yaml:"ok"`
	Result User `json:"result" yaml:"result"`
}

type ResponseGetUpdates struct {
	OK     bool     `json:"ok" yaml:"ok"`
	Result []Update `json:"result" yaml:"result"`
}

// TODO: "setMyCommands"
// TODO: "sendMessage"
