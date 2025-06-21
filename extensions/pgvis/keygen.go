package pgvis

import (
	"database/sql"
	"errors"
)

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	ApiKey   string `json:"api_key"`
}

func NewUser(id int64, userName string) *User {
	// TODO: Generate api key for user (id), but check database first, return existing user if available
	apiKey := "<api-key>"

	return &User{
		ID:       id,
		UserName: userName,
		ApiKey:   apiKey,
	}
}

func createUser(db *sql.DB, string, user *User) error {
	// TODO: Create a new user inside the pg vis server database

	return errors.New("under construction")
}
