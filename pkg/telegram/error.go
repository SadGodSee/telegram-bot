package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnathoraized = errors.New("user is not authorized")
	errUnableToSave = errors.New("unavle to safe")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.messages.Default)
	switch err {
	case errInvalidURL:
		msg.Text = b.messages.InvalidURL
		b.bot.Send(msg)
	case errUnathoraized:
		msg.Text = b.messages.Unauthorized
		b.bot.Send(msg)
	case errUnableToSave:
		msg.Text = b.messages.UnableToSave
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}
}
