package telegram

import (
	"github.com/SadGodSee/telegram-bot-one/pkg/config"
	"github.com/SadGodSee/telegram-bot-one/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string

	messages config.Messages
}

// инициализация полей извне, понижаем зону ответственности
// также получаем возомжность легко тестировать
func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, tr repository.TokenRepository, redirectURL string, messages config.Messages) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, redirectURL: redirectURL, tokenRepository: tr, messages: messages}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel() // создание канала на получение обновлений
	if err != nil {
		return err
	}

	b.handleUpdates(updates) // обработка обновлений
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // Если не сообщение, то скип
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue // добавляем, чтобы не обрабатывать 2 раза команду
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
