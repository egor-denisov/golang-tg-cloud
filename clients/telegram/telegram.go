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

func New(token string, db db.DataBase) TgClient {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	return TgClient{
		bot: bot,
		updates: updates,
		db: db,
	}
}

func (cl *TgClient) Listen() {
	for update := range cl.updates {
		if update.Message == nil {
			continue
		}
		
		if err := cl.proccessMessage(update.Message); err != nil {
			answer := Content{Text: proccessError(err)}
			cl.sendMedia(update.Message.Chat.ID, answer)
			log.Print(err)
			continue
		}
		log.Print("Message received")
	}
}

func (cl *TgClient) proccessMessage(msg *tgbotapi.Message) (err error) {
	defer func() { err = e.WrapIfErr("can`t proccess message", err) }()

	if msg.IsCommand() {
		return cl.proccessCommand(msg)
	}else if msg.Photo != nil || msg.Document != nil {
		return cl.proccessFile(msg)
	}else if msg.Text != ""{
		return cl.proccessText(msg)
	}

	return errors.New("Unknown type of message")
}

func (cl *TgClient) proccessFile(msg *tgbotapi.Message) error {
	var fileInfo File
	if msg.Photo != nil {
		fileInfo = File{
			Name: "photo" + msg.Photo[0].FileUniqueID + ".jpg",
			FileId: msg.Photo[0].FileID,
			FileUniqueId: msg.Photo[0].FileUniqueID,
			FileSize: msg.Photo[0].FileSize,
		}
	}else{
		fileInfo = File{
			Name: msg.Document.FileName,
			FileId: msg.Document.FileID,
			FileUniqueId: msg.Document.FileUniqueID,
			FileSize: msg.Document.FileSize,
			
		}
	}

	if _, err := cl.db.CreateNewFile(msg.From.ID, fileInfo); err != nil {	
		return err
	}

	return cl.makeReplyAfterAdding(msg.From.ID, fileInfo.Name)
}

func (cl *TgClient) proccessText(msg *tgbotapi.Message) error{
	if msg.Text[:2] == "./" || msg.Text[:3] == "../" {
		return cl.makeReplyAfterRequestingDirectory(msg.From.ID, msg.Text)
	}

	return cl.makeReplyAfterRequestingFile(msg.From.ID, msg.Text)
}

func (cl *TgClient) makeReplyAfterAdding(userId int64, fileName string) (err error) {
	defer func() { err = e.WrapIfErr("can`t make reply after adding", err) }()

	keyboard, err := cl.instantiateKeyboardNavigator(userId)

	replyContent := Content{
		Text: fmt.Sprintf("File %s successfully added", fileName),
		Keyboard: keyboard,
	}

	return cl.sendMedia(userId, replyContent)
}

func (cl *TgClient) makeReplyAfterRequestingFile(userId int64, reqString string) (err error) {
	defer func() { err = e.WrapIfErr("can`t make reply after requesting file", err) }()

	var file File

	availableFiles, err := cl.db.GetAvailableFiles(userId)
	for _, f := range availableFiles {
		if f.Name == reqString {
			file = f
			break
		}
	}

	if file.Name != reqString {
		return fmt.Errorf("file is not available from this folder")
	}

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

	return cl.sendMedia(userId, content)
}

func (cl *TgClient) makeReplyAfterRequestingDirectory(userId int64, reqString string) (err error) {
	defer func() { err = e.WrapIfErr("can`t make reply after requesting directory", err) }()
	
	var directory Directory

	if reqString == "../" {
		currentDirectoryId, err := cl.db.GetCurrentDirectory(userId)
		if err != nil {
			return err
		}

		directory, err = cl.db.GetParentDirectory(currentDirectoryId)
		if err != nil {
			return err
		}
		
	}else{
		availableDirectory, err := cl.db.GetAvailableDirectories(userId)
		if err != nil {
			return err
		}
		for _, d := range availableDirectory {
			if d.Name == reqString {
				directory = d
				break
			}
		}

		if directory.Name != reqString {
			return fmt.Errorf("directory is not available from this folder")
		}
	}

	if err := cl.db.JumpToDirectory(userId, directory.Id); err != nil {
		return err
	}

	keyboard, err := cl.instantiateKeyboardNavigator(userId)
	if err != nil {
		return err
	}

	content := Content{
		Text: fmt.Sprintf("Now you in '%s' folder", directory.Name),
		Keyboard: keyboard,
	}

	return cl.sendMedia(userId, content)
}

func (cl *TgClient) makeReplyAfterCreatingDirectory(userId int64, directoryName string) (err error) {
	defer func() { err = e.WrapIfErr("can`t make reply after creating directory", err) }()
	
	keyboard, err := cl.instantiateKeyboardNavigator(userId)
	replyContent := Content{
		Text: fmt.Sprintf("Directory %s created successfully", directoryName),
		Keyboard: keyboard,
	}

	return cl.sendMedia(userId, replyContent)
}

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
		if _, err := cl.bot.Send(msg); err != nil {
			return err
		}
		return nil
	}else {
		return errors.New("can`t find required type of message")
	}

	if _, err := cl.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

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