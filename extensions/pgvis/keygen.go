package pgvis

import (
	"fmt"
	"log/slog"
	"os/exec"
)

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	ApiKey   string `json:"api_key"`
}

func NewUser(id int64, userName string) (*User, error) {
	u := &User{
		ID:       id,
		UserName: userName,
		ApiKey:   "",
	}

	cmd := exec.Command("pg-vis", "user", "show", "--api-key", fmt.Sprintf("%d", id))
	if err := cmd.Run(); err != nil {
		if c, ok := err.(*exec.ExitError); !ok {
			return u, err
		} else {
			slog.Debug(fmt.Sprintf("Command failed with %d", c.ExitCode()))
			// TODO: Check the exit code, only continue if user not found,
			// 		 need to find out which exit code this is
		}
	} else {
		if apiKey, err := cmd.Output(); err != nil {
			return u, fmt.Errorf("output error: %s", err.Error())
		} else {
			u.ApiKey = string(apiKey)
		}
	}

	if u.ApiKey == "" {
		u.ApiKey = generateApiKey()
	}

	return u, nil
}

func generateApiKey() string {
	// TODO: Generate a new api key, Use the pg-vis command for this

	return "<api-key>"
}
