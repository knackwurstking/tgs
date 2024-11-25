package data

type User struct {
	ID                    int    `json:"id"`
	IsBot                 bool   `json:"is_bot"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`                // [Optional]
	Username              string `json:"username"`                 // [Optional]
	LanguageCode          string `json:"language_code"`            // [Optional]
	AddedToAttachmentMenu bool   `json:"added_to_attachment_menu"` // [Optional]
	SupportsInlineQueries bool   `json:"supports_inline_queries"`  // [Optional]
}
