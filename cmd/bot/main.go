package main

import (
	"github.com/SadGodSee/telegram-bot-one/pkg/config"
	"github.com/SadGodSee/telegram-bot-one/pkg/repository"
	"github.com/SadGodSee/telegram-bot-one/pkg/repository/boltdb"
	"github.com/SadGodSee/telegram-bot-one/pkg/server"
	"github.com/SadGodSee/telegram-bot-one/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURL, cfg.Messages)

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.TelegramBotURL)

	go func() { // используем горутину, для блок операции, чтобы потом можно было запустит сервер
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil { // аналогично является блокирующей операцией, из за ListenServe
		log.Fatal(err)
	}

}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			log.Fatal(err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
