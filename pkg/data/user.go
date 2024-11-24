package data

type User struct {
	ID                      int     `json:"id" yaml:"id"`
	IsBot                   bool    `json:"is_bot" yaml:"is_bot"`
	FirstName               string  `json:"first_name" yaml:"first_name"`
	LastName                *string `json:"last_name" yaml:"last_name"`
	Username                *string `json:"username" yaml:"username"`
	LanguageCode            *string `json:"language_code" yaml:"language_code"`
	IsPremium               *bool   `json:"is_premium" yaml:"is_premium"`
	AddedToAttachmentMenu   *bool   `json:"added_to_attachment_menu" yaml:"added_to_attachment_menu"`
	CanJoinGroups           *bool   `json:"can_join_groups" yaml:"can_join_groups"`
	CanReadAllGroupMessages *bool   `json:"can_read_all_group_messages" yaml:"can_read_all_group_messages"`
	SupportsInlineQueries   *bool   `json:"supports_inline_queries" yaml:"supports_inline_queries"`
	CanConnectToBusiness    *bool   `json:"can_connect_to_business" yaml:"can_connect_to_business"`
	HasMainWebApp           *bool   `json:"has_main_web_app" yaml:"has_main_web_app"`
}
