package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/SuperPaintman/nice/cli"
	"github.com/knackwurstking/tgs/pkg/tgs"
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

					if botCommand, err := parseUpdatesResponse(resp); err != nil {
						slog.Warn("Command not found for response", "response", *resp)
						continue
					} else {
						// TODO: Handle command
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

func parseUpdatesResponse(resp *tgs.ResponseGetUpdates) (*BotCommand, error) {
	// TODO: Check response for commands, dont forget to validate IDs

	return nil, fmt.Errorf("under construction")
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
