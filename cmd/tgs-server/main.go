package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/SuperPaintman/nice/cli"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

var (
	handledUpdateIDs = make([]int, 0)
)

func main() {
	app := cli.App{
		Name:  "tgs-server",
		Usage: cli.Usage("Telegram scripts server, the scripts part was kicked from the project :)"),
		Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
			var configPath string

			cli.StringVar(cmd, &configPath, "config",
				cli.Usage("Path to server configuration (yaml)"),
				cli.Required,
			)

			return func(cmd *cli.Command) error {
				config := NewConfig()

				if err := loadConfig(config, configPath); err != nil {
					return err
				}

				if err := checkConfig(config); err != nil {
					return err
				}

				requestTimeout := 60 // 1 Minute
				getUpdates := tgs.RequestGetUpdates{
					Timeout: &requestTimeout,
				}
				for {
					resp, err := getUpdates.Send()
					if err != nil {
						slog.Warn("Request updates", "error", err)
						continue
					}

					if !resp.OK {
						slog.Error("Request updates", "response", *resp)
						return fmt.Errorf("request updates failed")
					}

					if err := handleUpdates(resp.Result); err != nil {
						slog.Warn("Command not found for response", "response", *resp)
						continue
					}
				}
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.0.0"),
		},
	}

	app.HandleError(app.Run())
}

func handleUpdates(result []tgs.Update) error {
	defer cleanUpHandledUpdateIDs(result)

	for _, update := range result {
		if !isNewUpdateID(update.UpdateID) {
			continue
		}

		// TODO: ...
	}

	return fmt.Errorf("under construction")
}

func isNewUpdateID(id int) bool {
	for _, handledID := range handledUpdateIDs {
		if handledID == id {
			return false
		}
	}

	return true
}

func cleanUpHandledUpdateIDs(result []tgs.Update) {
	newHandledUpdateIDs := make([]int, 0)

	for _, handledID := range handledUpdateIDs {
		for _, update := range result {
			if update.UpdateID == handledID {
				newHandledUpdateIDs = append(newHandledUpdateIDs, handledID)
				break
			}
		}
	}

	handledUpdateIDs = newHandledUpdateIDs
}

func checkConfig(config *Config) error {
	if config.Token == "" {
		return fmt.Errorf("token missing")
	}

	if config.Targets.Chats == nil && config.Targets.Users == nil {
		return fmt.Errorf("missing targets")
	}

	return nil
}

func loadConfig(config *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, config)
}
