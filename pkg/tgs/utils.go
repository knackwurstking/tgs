package tgs

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func CheckTargets(message *tgbotapi.Message, targets *Targets) bool {
	if targets == nil {
		return false
	}

	if targets.All {
		return true
	}

	// User ID check
	if targets.Users != nil {
		if checkUserID(message.From.ID, targets.Users) {
			return true
		}
	}

	if targets.Chats != nil {
		for _, chat := range targets.Chats {
			if chat.ID == message.Chat.ID && (string(chat.Type) == message.Chat.Type || chat.Type == "") {
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

func CheckTargetsForUser(id int64, targets *Targets) bool {
	if targets == nil {
		return false
	}

	if targets.All {
		return true
	}

	// TODO: Check target

	return false
}

func CheckCallbackQueryTargets(callbackQuery *tgbotapi.CallbackQuery, targets *Targets) bool {
	if targets == nil {
		return false
	}

	if targets.All {
		return true
	}

	qFrom := callbackQuery.From

	if targets.Users != nil {
		if checkUserID(qFrom.ID, targets.Users) {
			return true
		}
	}

	qMessage := callbackQuery.Message
	qChat := qMessage.Chat

	if targets.Chats != nil {
		for _, chat := range targets.Chats {
			if chat.ID == qChat.ID && (string(chat.Type) == qChat.Type || chat.Type == "") {
				if !qChat.IsForum {
					if chat.Users == nil {
						return true
					}

					return checkUserID(qFrom.ID, chat.Users)
				}

				if chat.MessageThreadID <= 0 || chat.MessageThreadID == qMessage.MessageThreadID {
					if chat.Users == nil {
						return true
					}

					return checkUserID(qMessage.From.ID, chat.Users)
				}
			}
		}
	}

	return false
}

func checkUserID(id int64, users []UserTarget) bool {
	for _, user := range users {
		if user.ID == id {
			return true
		}
	}

	return false
}
