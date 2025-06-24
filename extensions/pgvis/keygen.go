package pgvis

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	ApiKey   string `json:"api_key"`
}

func NewUser(id int64, userName string) *User {
	u := &User{
		ID:       id,
		UserName: userName,
		ApiKey:   "",
	}

	// TODO: Generate api key for user (id), but check database first, return existing user if available
	apiKey := "<api-key>"

	return u
}
