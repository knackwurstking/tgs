package data

type Chat struct {
	ID        int     `json:"id" yaml:"id"`                 // Unique identifier for this chat. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
	Type      string  `json:"type" yaml:"type"`             // Possible types: "private", “group”, “supergroup” or “channel”
	Title     *string `json:"title" yaml:"title"`           // Optional. Title, for supergroups, channels and group chats
	Username  *string `json:"username" yaml:"username"`     // Optional. Username, for private chats, supergroups and channels if available
	FirstName *string `json:"first_name" yaml:"first_name"` // Optional. First name of the other party in a private chat
	LastName  *string `json:"last_name" yaml:"last_name"`   // Optional. Last name of the other party in a private chat
	IsForum   *bool   `json:"is_forum" yaml:"is_forum"`     // Optional. True, if the supergroup chat is a forum (has topics enabled)
}
