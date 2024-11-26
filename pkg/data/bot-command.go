package data

type BotCommand struct {
	Command     string `json:"command"`     //  Text of the command; 1-32 characters. Can contain only lowercase English letters, digits and underscores.
	Description string `json:"description"` // Description of the command; 1-256 characters.
}
