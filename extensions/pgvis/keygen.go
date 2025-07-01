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

	cmd := exec.Command("pg-vis", "user", "show", "--api-key", fmt.Sprintf("%d", id))

	if out, err := cmd.Output(); err != nil {
		// Error handling
		if c, ok := err.(*exec.ExitError); !ok {
			return u, err
		} else {
			slog.Debug(fmt.Sprintf("Command failed with %d", c.ExitCode()))

			// NOTE: For now, 10 is the code used for not found
			if c.ExitCode() != PGVisExitCodeNotFound {
				return u, fmt.Errorf(
					"pg-vis command failed with an exit code %d",
					c.ExitCode(),
				)
			}

			// Add a new user to the pg-vis database
			cmd = exec.Command("pg-vis", "user", "add", fmt.Sprintf("%d", id), userName)
			if err := cmd.Run(); err != nil {
				return u, fmt.Errorf("creating user failed: %s", err.Error())
			}
		}
	} else {
		// Get the api-key from the command output
		u.ApiKey = string(out)
	}

	// If the user has no api-key, create a new one one
	if u.ApiKey == "" {
		cmd = exec.Command("pg-vis", "api-key")

		out, err := cmd.Output()
		if err != nil {
			return u, fmt.Errorf("generating a new api key failed: %s", err.Error())
		}

		u.ApiKey = string(out)
	}

	cmd = exec.Command("pg-vis", "user", "mod", "--api-key", u.ApiKey, fmt.Sprintf("%d", u.ID))
	if err := cmd.Run(); err != nil {
		return u, fmt.Errorf("adding api-key for user \"%d\" failed: %s", u.ID, err.Error())
	}

	return u, nil
}
