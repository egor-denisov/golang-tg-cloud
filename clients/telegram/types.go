package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Content struct {
	Text     string
	Document *tgbotapi.Document
	Photo    []tgbotapi.PhotoSize
}

type User struct {
	Id int
	UserName string
	UserId int
	FirstName string
	LastName string
}

type Directory struct {
	Id int
	Name string
	Files []int
	Directories []int
	Size int
}

type File struct {
	Id int
	Name string
	FileId string
	FileUniqueId string
	FileSize int
}
