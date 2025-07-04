package pgvis

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/goombaio/namegenerator"
)

const (
	PGVisExitCodeNotFound      = 10
	PGVisExitCodeAlreadyExists = 20 // TODO: Need to implement this
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
		if u.UserName == "" {
			log.Debugf("Telegram user \"%d\" is missing a user name, generate one...", u.ID)

			var (
				cmd *exec.Cmd
				err error
				c   *exec.ExitError
				ok  bool
			)

			for true {
				u.UserName = generateUserName()

				cmd = exec.Command("pg-vis", "user", "mod", fmt.Sprintf("%d", u.ID), u.UserName)

				if err = cmd.Run(); err != nil {
					if c, ok = err.(*exec.ExitError); !ok {
						return nil, fmt.Errorf("adding api-key for user \"%d\" failed: %s", u.ID, err.Error())
					} else {
						if c.ExitCode() == PGVisExitCodeAlreadyExists {
							continue
						} else {
							log.Debugf("...Found a new user name for \"%d\": %s", u.ID, u.UserName)
							break
						}
					}
				}
			}
		}
	}

	return u, nil
}

func generateUserName() string {
	g := namegenerator.NewNameGenerator(time.Now().Unix())
	return g.Generate()
}
