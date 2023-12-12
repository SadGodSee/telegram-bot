package telegram

import (
	"context"
	"fmt"
	"github.com/SadGodSee/telegram-bot-one/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// здесь будет все связанное с авторизацией нашего приложения

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return nil
	}
	msg := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf(b.messages.Start, authLink))
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessTokens)
}

func (b *Bot) generateAuthorizationLink(chatID int64) (string, error) {
	redirectURL := b.generateRedirectLink(chatID)

	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectURL)
	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

func (b *Bot) generateRedirectLink(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}
