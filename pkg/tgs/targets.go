package tgs

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Targets struct {
	Users []UserTarget `yaml:"users,omitempty"`
	Chats []ChatTarget `yaml:"chats,omitempty"`
	All   bool         `yaml:"all,omitempty"`
}

func NewTargets() *Targets {
	return &Targets{
		Users: make([]UserTarget, 0),
		Chats: make([]ChatTarget, 0),
	}
}

type UserTarget struct {
	ID int64 `yaml:"id"`
}

type ChatTarget struct {
	Type            string       `yaml:"type"`
	ID              int64        `yaml:"id"`
	MessageThreadID int          `yaml:"message_thread_id,omitempty"`
	Users           []UserTarget `yaml:"users,omitempty"`
}

func CheckTargets(message *tgbotapi.Message, targets *Targets) bool {
	if targets == nil {
		return false
	}

	if targets.All {
		return true
	}

	checkUserID := func(id int64, users []UserTarget) bool {
		if users == nil {
			users = targets.Users
		}

		for _, user := range users {
			if user.ID == id {
				return true
			}
		}

		return false
	}

	// User ID check
	if targets.Users != nil {
		checkUserID(message.From.ID, nil)
	}

	// Chat ID check & message thread ID if chat is forum

	if targets.Chats != nil {
		for _, chat := range targets.Chats {
			if chat.ID == message.Chat.ID && (chat.Type == message.Chat.Type || chat.Type == "") {
				if !message.Chat.IsForum {
					if chat.Users == nil {
						return true
					}

					return checkUserID(message.From.ID, chat.Users)
				}

				if chat.MessageThreadID <= 0 || chat.MessageThreadID == message.MessageThreadID {
					if chat.Users == nil {
						return true
					}

					return checkUserID(message.From.ID, chat.Users)
				}
			}
		}
	}

	return false
}
