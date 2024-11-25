package data

// TODO: Keep only needed fields
type User struct {
	ID                      int    `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`                   // [Optional]
	Username                string `json:"username"`                    // [Optional]
	LanguageCode            string `json:"language_code"`               // [Optional]
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu"`    // [Optional]
	CanJoinGroups           bool   `json:"can_join_groups"`             // [Optional]
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"` // [Optional]
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`     // [Optional]
	HasMainWebApp           bool   `json:"has_main_web_app"`            // [Optional]
}
