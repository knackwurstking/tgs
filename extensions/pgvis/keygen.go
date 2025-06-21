package pgvis

func GenApiKey() string {
	// TODO: ...

	return ""
}

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	ApiKey   string `json:"api_key"`
}

func NewUser(id int64, userName string, apiKey string) *User {
	// TODO: ...

	return nil
}
