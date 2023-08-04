package types

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Content struct {
	Text     string
	Document *tgbotapi.Document
	Photo    []tgbotapi.PhotoSize
	Keyboard interface{}
}

type User struct {
	Id int
	UserName string
	UserId int
	FirstName string
	LastName string
}

type Directory struct {
	Id int `json:"id"`
	ParentId int `json:"parentId"`
	Name string `json:"name"`
	UserId int `json:"userId"`
	Files []int `json:"files"`
	Directories []int `json:"directories"`
	Size int `json:"size"`
}

type File struct {
	Id int
	Name string
	FileId string
	FileUniqueId string
	FileSize int
}
