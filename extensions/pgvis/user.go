package pgvis

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
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
	user := &User{
		ID:       id,
		UserName: userName,
	}

	// Check if the user already exists in the pg-vis database
	cmd := exec.Command("pg-vis", "user", "show", "--api-key", fmt.Sprintf("%d", id))
	out, err := cmd.Output()
	if err != nil {
		// Error handling
		c, ok := err.(*exec.ExitError)
		if !ok {
			return nil, err
		}

		log.Debugf("Command failed with %d", c.ExitCode())

		// Check if the error is due to the user not being found
		if c.ExitCode() != PGVisExitCodeNotFound {
			return nil, fmt.Errorf(
				"pg-vis command failed with an exit code %d",
				c.ExitCode(),
			)
		}

		// Generate a new API key first
		user.ApiKey, err = generateApiKey()
		if err != nil {
			return nil, err
		}

		// Add new user with the generated API key
		cmd = exec.Command("pg-vis", "user", "add", strconv.Itoa(int(user.ID)), user.UserName, user.ApiKey)
		if err = cmd.Run(); err != nil {
			return nil, err
		}
	} else {
		if out := strings.TrimSpace(string(out)); out != "" {
			user.ApiKey = strings.TrimSpace(string(out))
		} else {
			user.ApiKey, err = generateApiKey()
			if err != nil {
				return nil, err
			}

			// Modify the users api-key using the pg-vis command
			cmd = exec.Command("pg-vis", "user", "mod", "--api-key", user.ApiKey, strconv.Itoa(int(user.ID)))
			if err = cmd.Run(); err != nil {
				return nil, err
			}
		}
	}

	// Final safety check, should never happen :)
	if user.ApiKey == "" {
		return nil, errors.New("failed to generate API key")
	}

	return user, nil
}

func generateApiKey() (string, error) {
	cmd := exec.Command("pg-vis", "api-key")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
