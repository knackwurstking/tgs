package pgpress

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

const (
	PGPressExitCodeNotFound = 10
)

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	ApiKey   string `json:"api_key"`
}

func NewUser(id int64, userName string) (*User, error) {
	log.Debug("Creating PGPress user",
		"user_id", id,
		"username", userName,
	)

	user := &User{
		ID:       id,
		UserName: userName,
	}

	// Check if the user already exists in the pg-press database
	log.Debug("Checking if user exists in pg-press database",
		"user_id", id,
		"command", "pg-press user show --api-key",
	)
	cmd := exec.Command("pg-press", "user", "show", "--api-key", fmt.Sprintf("%d", id))
	out, err := cmd.Output()
	if err != nil {
		// Error handling
		c, ok := err.(*exec.ExitError)
		if !ok {
			return nil, err
		}

		log.Debug("pg-press user show command failed",
			"exit_code", c.ExitCode(),
			"user_id", id,
		)

		// Check if the error is due to the user not being found
		if c.ExitCode() != PGPressExitCodeNotFound {
			log.Error("pg-press command failed with unexpected exit code",
				"exit_code", c.ExitCode(),
				"user_id", id,
				"expected_not_found_code", PGPressExitCodeNotFound,
			)
			return nil, fmt.Errorf(
				"pg-press command failed with an exit code %d",
				c.ExitCode(),
			)
		}

		log.Info("User not found in pg-press database, creating new user",
			"user_id", id,
			"username", userName,
		)

		// Generate a new API key first
		log.Debug("Generating new API key for user", "user_id", id)
		user.ApiKey, err = generateApiKey()
		if err != nil {
			log.Error("Failed to generate API key",
				"user_id", id,
				"error", err,
			)
			return nil, err
		}

		log.Debug("Generated API key successfully",
			"user_id", id,
			"api_key_length", len(user.ApiKey),
		)

		// Add new user with the generated API key
		log.Debug("Adding new user to pg-press database",
			"user_id", id,
			"username", user.UserName,
		)
		cmd = exec.Command("pg-press", "user", "add", strconv.Itoa(int(user.ID)), user.UserName, user.ApiKey)
		if err = cmd.Run(); err != nil {
			log.Error("Failed to add user to pg-press database",
				"user_id", id,
				"username", user.UserName,
				"error", err,
			)
			return nil, err
		}

		log.Info("New user added to pg-press database successfully",
			"user_id", id,
			"username", user.UserName,
		)
	} else {
		log.Info("User found in pg-press database", "user_id", id)

		if out := strings.TrimSpace(string(out)); out != "" {
			user.ApiKey = strings.TrimSpace(string(out))
			log.Debug("Retrieved existing API key",
				"user_id", id,
				"api_key_length", len(user.ApiKey),
			)
		} else {
			log.Debug("User exists but has no API key, generating new one", "user_id", id)

			user.ApiKey, err = generateApiKey()
			if err != nil {
				log.Error("Failed to generate API key for existing user",
					"user_id", id,
					"error", err,
				)
				return nil, err
			}

			log.Debug("Generated new API key for existing user",
				"user_id", id,
				"api_key_length", len(user.ApiKey),
			)

			// Modify the users api-key using the pg-press command
			log.Debug("Updating user API key in pg-press database", "user_id", id)
			cmd = exec.Command("pg-press", "user", "mod", "--api-key", user.ApiKey, strconv.Itoa(int(user.ID)))
			if err = cmd.Run(); err != nil {
				log.Error("Failed to update user API key",
					"user_id", id,
					"error", err,
				)
				return nil, err
			}

			log.Debug("User API key updated successfully", "user_id", id)
		}
	}

	// Final safety check, should never happen :)
	if user.ApiKey == "" {
		log.Error("API key is empty after user creation/retrieval",
			"user_id", id,
			"username", userName,
		)
		return nil, errors.New("failed to generate API key")
	}

	log.Info("PGPress user created/retrieved successfully",
		"user_id", id,
		"username", userName,
		"api_key_set", user.ApiKey != "",
	)

	return user, nil
}

func generateApiKey() (string, error) {
	log.Debug("Executing pg-press api-key command")

	cmd := exec.Command("pg-press", "api-key")
	out, err := cmd.Output()
	if err != nil {
		log.Error("pg-press api-key command failed",
			"command", "pg-press api-key",
			"error", err,
		)
		return "", err
	}

	apiKey := strings.TrimSpace(string(out))
	log.Debug("API key generated successfully",
		"key_length", len(apiKey),
		"output_size", len(out),
	)

	return apiKey, nil
}
