package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Content struct {
	Text     string
	Document *tgbotapi.Document
	Photo    []tgbotapi.PhotoSize
}

type User struct {
	Id int64
	UserName string
	UserId int64
	FirstName string
	LastName string
}

type Directory struct {
	Name string
	Files []int
}