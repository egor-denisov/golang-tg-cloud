package types

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Content struct {
	Text     string `json:"text"`
	Document *tgbotapi.Document `json:"document"`
	Photo    []tgbotapi.PhotoSize `json:"photo"`
	Keyboard interface{} `json:"keyboard"`
}

type User struct {
	Id int `json:"id"`
	UserName string `json:"username"`
	UserId int `json:"user_id"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	CurrentDirectory int `json:"current_directory"`
}

type Directory struct {
	Id int `json:"id"`
	ParentId int `json:"parent_id"`
	Name string `json:"name"`
	UserId int `json:"user_id"`
	Files []int `json:"files"`
	Directories []int `json:"directories"`
	Size int `json:"size"`
}

type File struct {
	Id int `json:"id"`
	Name string `json:"name"`
	FileId string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize int `json:"file_size"`
	ThumbnailFileId string `json:"thumbnail_file_id"`
	ThumbnailSource string `json:"thumbnail_source"`
	FileSource string `json:"file_source"`
}

type DirectoryContent struct {
	Directories []Directory `json:"directories"`
	Files []File `json:"files"`
}