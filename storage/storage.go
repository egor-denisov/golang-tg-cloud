package storage

import (
	"log"
	"main/db"
	"main/lib/e"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type Storage struct {
	bot *tgbotapi.BotAPI
	db  db.DataBase
}

func New(token string, db db.DataBase) Storage {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Storage activated on account %s", bot.Self.UserName)

	return Storage{
		bot: bot,
		db: db,
	}
}

func (s *Storage) GetFile(fileId string) (url string, err error) {
	defer func() { err = e.WrapIfErr("can`t get file from storage", err) }()

	return s.bot.GetFileDirectURL(fileId)
}