package main

import (
	"fmt"

	"github.com/SuperPaintman/nice/cli"
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

				// TODO: Check token, register commands, enter the main loop and wait for updates

				return nil
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.0.0"),
		},
	}

	app.HandleError(app.Run())
}

func loadConfig(config *Config, path string) error {
	// TODO: Load YAML configuration

	return fmt.Errorf("under construction")
}
