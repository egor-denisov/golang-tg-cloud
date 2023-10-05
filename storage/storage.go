package storage

import (
	"io/ioutil"
	"log"
	"main/db"
	"main/lib/e"
	"main/lib/h"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UploadingItem struct {
	idChannel chan int
	path     string
	filename string
	user_id     int
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
		id, err := s.UploadFile(item)
		if err != nil {
			log.Print(err)
		}
		item.idChannel <- id
	}
}

// Function return file bytes by file ID
func (s *Storage) GetFileAsBytes(url string) (fileBytes []byte, err error) {
	defer func() { err = e.WrapIfErr("can`t get file from storage", err) }()

	// Doing GET request to url
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Return bytes of content
	return ioutil.ReadAll(resp.Body)
}

// Function for getting url based on telegram.org
func (s *Storage) GetFileURL(fileId string) (string, error) {
	return s.bot.GetFileDirectURL(fileId)
}

// Function upload a file to the tg server by sending message to user
func (s *Storage) UploadFile(item UploadingItem) (id int, err error) {
	defer func() { err = e.WrapIfErr("can`t upload file to storage", err) }()
	// Rename unique filename to original
	newPath := "./assets/" + item.filename
	if err := os.Rename(item.path, newPath); err != nil {
		return -1, err
	}
	defer func () { err = os.Remove(newPath) }()
	// Create new instance of document and sending it to user
	document := tgbotapi.NewDocument(int64(item.user_id), tgbotapi.FilePath(newPath))
	msg, err := s.bot.Send(document)
	if err != nil {
		return -1, err
	}
	// Create new instance of file
	file := h.GetFileDataFromMessage(msg)
	
	// Upload file to database
	id, err = s.db.CreateNewFile(int64(item.user_id), item.directoryId, file)
	if err != nil {
		return -1, err
	}

	return id, err
}

// Function for adding new item to queue channel
func (s *Storage) AddToUploadingQueue(idChannel chan int, path string, filename string, user_id int, directoryId int) {
	// Adding new item to queue channel
	s.uploadingQueue <- UploadingItem{
		idChannel: idChannel,
		path: path, 
		filename: filename, 
		user_id: user_id,
		directoryId: directoryId,
	}
}