package data

type User struct {
	ID        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`

	LastName              string `json:"last_name,omitempty"`                // [Optional]
	Username              string `json:"username,omitempty"`                 // [Optional]
	LanguageCode          string `json:"language_code,omitempty"`            // [Optional]
	AddedToAttachmentMenu *bool  `json:"added_to_attachment_menu,omitempty"` // [Optional]
	SupportsInlineQueries *bool  `json:"supports_inline_queries,omitempty"`  // [Optional]
}
