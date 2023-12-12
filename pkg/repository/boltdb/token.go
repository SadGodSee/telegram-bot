package boltdb

import (
	"errors"
	"github.com/SadGodSee/telegram-bot-one/pkg/repository"
	"github.com/boltdb/bolt"
	"strconv"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(chatID int64, token string, bucket repository.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		// дальше ам надо сохранить в базу
		return b.Put(intToBytes(chatID), []byte(token))
	})

}

func (r *TokenRepository) Get(chatID int64, bucket repository.Bucket) (string, error) {
	var token string
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		// дальше нам надо достать из базы значение по ключу
		token = string(b.Get(intToBytes(chatID)))
		return nil
	})
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", errors.New("token not found")
	}
	return token, nil
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
