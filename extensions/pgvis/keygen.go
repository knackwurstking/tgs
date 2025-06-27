package pgvis

import (
	"fmt"
	"log/slog"
	"os/exec"
)

const (
	PGVisExitCodeNotFound = 10
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

	// Get api key for this user (telegram id), if possible
	cmd := exec.Command("pg-vis", "user", "show", "--api-key", fmt.Sprintf("%d", id))

	if err := cmd.Run(); err != nil {
		if c, ok := err.(*exec.ExitError); !ok {
			// Command failed
			return u, err
		} else {
			slog.Debug(fmt.Sprintf("Command failed with %d", c.ExitCode()))

			// NOTE: For now, 10 is the code used for not found
			if c.ExitCode() != PGVisExitCodeNotFound {
				return u, fmt.Errorf(
					"pg-vis command failed with an exit code %d",
					c.ExitCode(),
				)
			} else {
				// TODO: Create user, using the pg-vis command here
			}
		}
	} else {
		// Get the api-key from the command output
		if apiKey, err := cmd.Output(); err != nil {
			return u, fmt.Errorf("output error: %s", err.Error())
		} else {
			u.ApiKey = string(apiKey)
		}
	}

	// If the user has no api-key, create a new one one
	if u.ApiKey == "" {
		u.ApiKey = generateApiKey()
		// TODO: Update user, using the pg-vis command here
	}

	return u, nil
}

func generateApiKey() string {
	// TODO: Generate a new api key, use the pg-vis command for this

	return "<api-key>"
}
