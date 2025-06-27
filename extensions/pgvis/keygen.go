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

			// NOTE: For now, 1 is the exit code in use for not found
			if c.ExitCode() != 1 {
				return u, fmt.Errorf(
					"pg-vis command failed with an exit code %d",
					c.ExitCode(),
				)
			}
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
	// TODO: Generate a new api key, use the pg-vis command for this, and
	// 		 create the user

	return "<api-key>"
}
