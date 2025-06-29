module github.com/knackwurstking/tgs

go 1.24

require (
	github.com/SuperPaintman/nice v0.0.0-20211001214957-a29cd3367b17
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/lmittmann/tint v1.0.5
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42 // indirect
)

//replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 => ../telegram-bot-api

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 => github.com/knackwurstking/telegram-bot-api/v5 v5.6.0
