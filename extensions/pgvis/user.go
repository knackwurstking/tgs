package pgvis

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
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

	{ // Get user from "pg-vis" or create a new one
		cmd := exec.Command("pg-vis", "user", "show", "--api-key", fmt.Sprintf("%d", id))

		if out, err := cmd.Output(); err != nil {
			// Error handling
			if c, ok := err.(*exec.ExitError); !ok {
				return nil, err
			} else {
				log.Debugf("Command failed with %d", c.ExitCode())

				// NOTE: For now, 10 is the code used for not found
				if c.ExitCode() != PGVisExitCodeNotFound {
					return nil, fmt.Errorf(
						"pg-vis command failed with an exit code %d",
						c.ExitCode(),
					)
				}

				// Add a new user to the pg-vis database
				cmd = exec.Command("pg-vis", "user", "add", fmt.Sprintf("%d", id), userName)
				if err := cmd.Run(); err != nil {
					return nil, fmt.Errorf("creating user failed: %s", err.Error())
				}
			}
		} else {
			// Get the api-key from the command output
			u.ApiKey = strings.Trim(string(out), "\n\r\t ")
		}
	}

	{ // Generate a new api key for this user if needed
		if u.ApiKey == "" {
			cmd := exec.Command("pg-vis", "api-key")

			out, err := cmd.Output()
			if err != nil {
				return nil, fmt.Errorf("generating a new api key failed: %s", err.Error())
			}

			u.ApiKey = strings.Trim(string(out), "\n\r\t ")
		}
	}

	{ // Mod: Api Key
		cmd := exec.Command("pg-vis", "user", "mod", "--api-key", u.ApiKey, fmt.Sprintf("%d", u.ID))
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("adding api-key for user \"%d\" failed: %s", u.ID, err.Error())
		}
	}

	{ // Mod: User Name
		if u.UserName == "" && userName != "" {
			cmd := exec.Command("pg-vis", "user", "mod", "--name", userName, fmt.Sprintf("%d", u.ID))
			if err := cmd.Run(); err != nil {
				log.Error("Update user name for \"%d\" failed: %s", u.ID, err.Error())
			} else {
				u.UserName = userName
			}
		}
	}

	return u, nil
}
