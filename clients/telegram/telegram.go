package telegram

import (
	"errors"
	"fmt"
	"log"
	"main/db"
	"main/lib/e"
	. "main/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgClient struct {
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	db 		db.DataBase
}
// Function for creating a new TgClient instance
func New(token string, db db.DataBase) TgClient {
	// Create a new tg bot instance
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	// Enabling debugging output
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	// Setting updating parameters
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	// Return instance of telegram client
	return TgClient{
		bot: bot,
		updates: updates,
		db: db,
	}
}
// Function which listens for updates
func (cl *TgClient) Listen() {
	for update := range cl.updates {
		if update.Message == nil {
			continue
		}
		// Processing new message and creating answer if error occurs
		if err := cl.proccessMessage(update.Message); err != nil {
			answer := Content{Text: proccessError(err)}
			cl.sendMedia(update.Message.Chat.ID, answer)
			log.Print(err)
			continue
		}
		log.Print("Message received")
	}
}
// Function for proccessing message and defining type of message
func (cl *TgClient) proccessMessage(msg *tgbotapi.Message) (err error) {
	defer func() { err = e.WrapIfErr("can`t proccess message", err) }()
	// Defining type of message
	if msg.IsCommand() {
		return cl.proccessCommand(msg)
	}else if msg.Photo != nil || msg.Document != nil {
		return cl.proccessFile(msg)
	}else if msg.Text != ""{
		return cl.proccessText(msg)
	}
	// If unknown type return error
	return errors.New("Unknown type of message")
}
// Function for processing message including file
func (cl *TgClient) proccessFile(msg *tgbotapi.Message) error {
	var fileInfo File
	// Setting information about the file
	if msg.Photo != nil {
		mainPhoto := msg.Photo[len(msg.Photo)-1]
		fileInfo = File{
			Name: "photo" + mainPhoto.FileUniqueID + ".jpg",
			FileId: mainPhoto.FileID,
			FileUniqueId: mainPhoto.FileUniqueID,
			FileSize: mainPhoto.FileSize,
			ThumbnailFileId: msg.Photo[0].FileID,
		}
	}else{
		fileInfo = File{
			Name: msg.Document.FileName,
			FileId: msg.Document.FileID,
			FileUniqueId: msg.Document.FileUniqueID,
			FileSize: msg.Document.FileSize,
		}
	}
	// Getting current directory
	directoryId, err := cl.db.GetCurrentDirectory(msg.From.ID)
	if err != nil {
		return err
	}
	// Creating a new file in the database
	if _, err := cl.db.CreateNewFile(msg.From.ID, directoryId, fileInfo); err != nil {	
		return err
	}
	// Making reply to the user
	return cl.makeReplyAfterAdding(msg.From.ID, fileInfo.Name)
}
// Function for processing message including text
func (cl *TgClient) proccessText(msg *tgbotapi.Message) error{
	// If the path starts with './' or equals '../' then we request a directory
	if msg.Text[:2] == "./" || msg.Text[:3] == "../" {
		return cl.makeReplyAfterRequestingDirectory(msg.From.ID, msg.Text)
	}
	// Else requesting a file
	return cl.makeReplyAfterRequestingFile(msg.From.ID, msg.Text)
}
// Function for creating respond after adding a new file
func (cl *TgClient) makeReplyAfterAdding(userId int64, fileName string) (err error) {
	defer func() { err = e.WrapIfErr("can`t make reply after adding", err) }()
	// Creating new keyboard with updated data
	keyboard, err := cl.instantiateKeyboardNavigator(userId)
	if err != nil {
		return err
	}
	// Creating content for replying
	replyContent := Content{
		Text: fmt.Sprintf("File %s successfully added", fileName),
		Keyboard: keyboard,
	}
	// Sending replying message
	return cl.sendMedia(userId, replyContent)
}
// Function for creating respond after requesting a file
func (cl *TgClient) makeReplyAfterRequestingFile(userId int64, reqString string) (err error) {
	defer func() { err = e.WrapIfErr("can`t make reply after requesting file", err) }()

	var file File
	// Getting current directory
	directoryId, err := cl.db.GetCurrentDirectory(userId)
	if err != nil {
		return err
	}
	// Getting available files from the database
	availableFiles, err := cl.db.GetAvailableFilesInDiretory(userId, directoryId)
	// Iterating through available files and finding matching name
	for _, f := range availableFiles {
		if f.Name == reqString {
			file = f
			break
		}
	}
	// Checking result
	if file.Name != reqString {
		return fmt.Errorf("file is not available from this folder")
	}
	// Creating content for replying
	var content Content
	if file.Name[:5] == "photo" {
		var photo []tgbotapi.PhotoSize
		photo = append(photo, tgbotapi.PhotoSize{
			FileID: file.FileId,
			FileUniqueID: file.FileUniqueId,
			FileSize: file.FileSize,
		}) 
		content = Content{Photo: photo}
	}else{
		document := tgbotapi.Document{
			FileID: file.FileId,
			FileUniqueID: file.FileUniqueId,
			FileName: file.Name,
			FileSize: file.FileSize,
		}
		content = Content{Document: &document}
	}
	// Sending replying message
	return cl.sendMedia(userId, content)
}
// Function for creating respond after requesting a directory
func (cl *TgClient) makeReplyAfterRequestingDirectory(userId int64, reqString string) (err error) {
	defer func() { err = e.WrapIfErr("can`t make reply after requesting directory", err) }()
	
	var directory Directory
	// If it`s a 'go back' then go back to the parent directory
	if reqString == "../" {
		// Getting current directory id
		currentDirectoryId, err := cl.db.GetCurrentDirectory(userId)
		if err != nil {
			return err
		}
		// And getting parent directory for the current
		directory, err = cl.db.GetParentDirectory(currentDirectoryId)
		if err != nil {
			return err
		}
	}else{
		// Getting current directory
		directoryId, err := cl.db.GetCurrentDirectory(userId)
		if err != nil {
			return err
		}
		// Getting available directories
		availableDirectory, err := cl.db.GetAvailableDirectoriesInDiretory(userId, directoryId)
		if err != nil {
			return err
		}
		// Iterating through available directories and finding matching name
		for _, d := range availableDirectory {
			if d.Name == reqString {
				directory = d
				break
			}
		}
		// Checking result
		if directory.Name != reqString {
			return fmt.Errorf("directory is not available from this folder")
		}
	}
	// Moving to the found directory
	if err := cl.db.JumpToDirectory(userId, directory.Id); err != nil {
		return err
	}
	// Creating new keyboard with updated data
	keyboard, err := cl.instantiateKeyboardNavigator(userId)
	if err != nil {
		return err
	}
	// Creating content for replying
	content := Content{
		Text: fmt.Sprintf("Now you in '%s' folder", directory.Name),
		Keyboard: keyboard,
	}
	// Sending replying message
	return cl.sendMedia(userId, content)
}
// Function for creating respond after creating a new directory
func (cl *TgClient) makeReplyAfterCreatingDirectory(userId int64, directoryName string) (err error) {
	defer func() { err = e.WrapIfErr("can`t make reply after creating directory", err) }()
	// Creating new keyboard with updated data
	keyboard, err := cl.instantiateKeyboardNavigator(userId)
	if err != nil {
		return err
	}
	// Creating content for replying
	replyContent := Content{
		Text: fmt.Sprintf("Directory %s created successfully", directoryName),
		Keyboard: keyboard,
	}
	// Sending replying message
	return cl.sendMedia(userId, replyContent)
}
// Function for sending message to user
func (cl *TgClient) sendMedia(chatID int64, content Content) (err error) {
	defer func() { err = e.WrapIfErr("can`t send media", err) }()

	var msg tgbotapi.Chattable
	if content.Document != nil {
		msg = tgbotapi.NewDocument(chatID, tgbotapi.FileID(content.Document.FileID))
	}else if content.Photo != nil {
		msg = tgbotapi.NewPhoto(chatID, tgbotapi.FileID(content.Photo[0].FileID))
	}else if content.Text != "" {
		msg := tgbotapi.NewMessage(chatID, content.Text)
		if content.Keyboard != nil{
			msg.ReplyMarkup = content.Keyboard
		}
		// Sending the message
		_, err = cl.bot.Send(msg)
		return err
	}else {
		return errors.New("can`t find required type of message")
	}
	// Sending the message
	_, err = cl.bot.Send(msg)
	return err
}
// Function for error handling
func proccessError(err error) string {
	switch err.Error() {
		case "can`t proccess message: directory or user not found!": 
			return "You need create account before sending files! For it send me /start"
		case "can`t proccess message: user already exists":
			return "You already have account!"
		case "can`t proccess message: file already exists in folder":
			return "This file already in folder"
		case "can`t proccess message: directory with this name already exists in current folder":
			return "Directory with this name already exists in current folder"
		case "can`t proccess message: wrong folder name":
			return "Sorry, but folder cannot have this name"
		case "can`t proccess message: can`t make reply after requesting: file is not available from this folder":
			return "Sorry, this file does not exist"
		case "can`t proccess message: can`t make reply after requesting directory: this directory is root":
			return "Sorry, this directory is root"
		default: 
			return "Sorry, i can`t understand u want to do :("
	}
}