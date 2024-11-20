# Telegram Scripts

## Telegram API:

https://core.telegram.org/bots/api

```
curl https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getMe | jq
```

```
curl https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates | jq
```

Set a Bot command

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"commands":[{"command":"test","description":"description with spaces"}]}' \
    https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setMyCommands
```

Send a message with newlines

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"parse_mode": "HTML", "chat_id": "<CHAT_ID>", "text": "<b>My title</b>%0AMy body"}' \
    https://api.telegram.org/bot<YOUR_BOT_TOKEN>/sendMessage
```

## TODOs

- [ ] Get IP command (no deps.)
- [ ] Get server logs (journalctl) command (no deps.)
