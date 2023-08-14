package storage

import (
	"io/ioutil"
	"log"
	"main/db"
	"main/lib/e"
	. "main/types"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UploadingItem struct {
	path     string
	filename string
	user     User
	directoryId int
}

type Storage struct {
	bot *tgbotapi.BotAPI
	db  db.DataBase
	uploadingQueue chan UploadingItem
}

// Function for creating new instnce of storage
func New(token string, db db.DataBase) Storage {
	// Create a new tg bot instance
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	// Enabling debugging output
	bot.Debug = true
	log.Printf("Storage activated on account %s", bot.Self.UserName)
	// Return instance of storage
	return Storage{
		bot: bot,
		db: db,
		uploadingQueue: make(chan UploadingItem),
	}
}

// Function which listen queue channel for updates
func (s *Storage) StartUploading() {
	// Listen queue channel
	for item := range s.uploadingQueue {
		if err := s.UploadFile(item); err != nil {
			log.Print(err)
		}
	}
}

// Function return file bytes by file ID
func (s *Storage) GetFileAsBytes(fileId string) (fileBytes []byte, err error) {
	defer func() { err = e.WrapIfErr("can`t get file from storage", err) }()

	// Getting url from file id
	url, err := s.bot.GetFileDirectURL(fileId)
	if err != nil {
		return nil, err
	}
	// Doing GET request to url
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Return bytes of content
	return ioutil.ReadAll(resp.Body)
}

// Function upload a file to the tg server by sending message to user
func (s *Storage) UploadFile(item UploadingItem) (err error) {
	defer func() { err = e.WrapIfErr("can`t upload file to storage", err) }()

	// Rename unique filename to original
	newPath := "./assets/" + item.filename
	if err := os.Rename(item.path, newPath); err != nil {
		return err
	}

	// Create new instance of document and sending it to user
	document := tgbotapi.NewDocument(int64(item.user.ChatId), tgbotapi.FilePath(newPath))
	msg, err := s.bot.Send(document)
	if err != nil {
		return err
	}
	// Create new instance of file
	file := File{
		Name: msg.Document.FileName,
		FileId: msg.Document.FileID,
		FileUniqueId: msg.Document.FileUniqueID,
		FileSize: msg.Document.FileSize,
	}
	// Upload file to database
	file.Id, err = s.db.CreateNewFile(int64(item.user.UserId), file)
	if err != nil {
		return err
	}
	// Adding uploaded file to directory
	if err := s.db.AddFileToDirectory(item.directoryId, file); err != nil {
		return err
	}

	// Remove temp file
	if err := os.Remove(newPath); err != nil {
		return err
	}

	return err
}

// Function for adding new item to queue channel
func (s *Storage) AddToUploadingQueue(path string, filename string, user User, directoryId int) {
	// Adding new item to queue channel
	s.uploadingQueue <- UploadingItem{
		path: path, 
		filename: filename, 
		user: user,
		directoryId: directoryId,
	}
}