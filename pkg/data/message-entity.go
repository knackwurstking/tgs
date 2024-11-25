package data

type MessageEntity struct {
	Type          string `json:"type"`
	Offset        int    `json:"offset"`
	Length        int    `json:"length"`
	URL           string `json:"url"`             // [Optional]
	User          *User  `json:"user"`            // [Optional]
	Language      string `json:"language"`        // [Optional]
	CustomEmojiID string `json:"custom_emoji_id"` // [Optional]
}
