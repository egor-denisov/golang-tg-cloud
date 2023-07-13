package telegram

import (
	"errors"
	"log"
	"main/db"
	"main/lib/e"
	"strconv"

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
		err := cl.proccessMessage(update.Message)
		if err != nil {
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

	var content Content

	if msg.IsCommand() {
		content = Content{}
		switch msg.Command() {
		case "help":
			content.Text = "You can use /search for search file. Also you can see all files in current directory with /showAll."
		case "start":
			content.Text = "Hello! Send me a file or create a new folder"
			userInfo := User{
				UserId: int(msg.From.ID),
				UserName: msg.From.UserName,
				FirstName: msg.From.FirstName,
				LastName: msg.From.LastName,
			}
			_, err := createNewUser(cl.db, userInfo)
			if err != nil {	
				return err
			}
		case "search":
			content.Text = "Input search string."
		case "show_all":
			content.Text = "Show all files and directories"
		case "mainFolder":
			content.Text = "Jump to main directory"
		default:
			content.Text = "I don't know that command"
		}
		if err := cl.sendMedia(msg.Chat.ID, content); err != nil {
			return err
		}
		return nil
	}else if msg.Photo != nil {
		fileInfo := File{
			Name: "photo" + msg.Photo[0].FileUniqueID,
			FileId: msg.Photo[0].FileID,
			FileUniqueId: msg.Photo[0].FileUniqueID,
			FileSize: msg.Photo[0].FileSize,
		}
		
		_, err := createNewFile(cl.db, int(msg.From.ID), fileInfo)
		if err != nil {	
			return err
		}
	}else if msg.Document != nil {
		fileInfo := File{
			Name: msg.Document.FileName,
			FileId: msg.Document.FileID,
			FileUniqueId: msg.Document.FileUniqueID,
			FileSize: msg.Document.FileSize,
		}

		_, err := createNewFile(cl.db, int(msg.From.ID), fileInfo)
		if err != nil {	
			return err
		}
	}else if msg.Text != ""{
		if err := cl.MakeReplyAfterRequest(msg.From.ID, msg.Text); err != nil{
			return err
		}
		return nil
	}else{
		return errors.New("Unknown type of message")
	}
	
	if err := cl.makeReplyAfterAdding(msg.From.ID); err != nil{
		return err
	}
	return nil
}

func (cl *TgClient) makeReplyAfterAdding(userId int64) (err error) {
	defer func() { err = e.WrapIfErr("can`t make reply after adding", err) }()

	keyboard, err := cl.instantiateKeyboardNavigator(userId)

	replyContent := Content{
		Text: "File successfully added",
		Keyboard: keyboard,
	}

	if err := cl.sendMedia(userId, replyContent); err != nil{
		return err
	}
	return nil
}

func (cl *TgClient) MakeReplyAfterRequest(userId int64, fileName string) (err error) {
	defer func() { err = e.WrapIfErr("can`t make reply after adding", err) }()

	file, err := getFileByName(cl.db, fileName) 
	if err != nil {
		return err
	}

	if(file.Name[:5] == "photo") {
		var photo []tgbotapi.PhotoSize
		photo = append(photo, tgbotapi.PhotoSize{
			FileID: file.FileId,
			FileUniqueID: file.FileUniqueId,
			FileSize: file.FileSize,
		}) 
		if err := cl.sendMedia(userId, Content{Photo: photo}); err != nil{
			return err
		}
	}else{
		document := tgbotapi.Document{
			FileID: file.FileId,
			FileUniqueID: file.FileUniqueId,
			FileName: file.Name,
			FileSize: file.FileSize,
		}
		if err := cl.sendMedia(userId, Content{Document: &document}); err != nil{
			return err
		}
	}
	return nil
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
		if content.Keyboard.Keyboard != nil{
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

func (cl *TgClient) instantiateKeyboardNavigator(userID int64) (tgbotapi.ReplyKeyboardMarkup, error) {
	files, err := getAvailableFiles(cl.db, strconv.Itoa(int(userID)))
	if err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}, err
	}
	directories, err := getAvailableDirectories(cl.db, strconv.Itoa(int(userID)))
	if err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}, err
	}
	return createKeyboardNavigator(directories, files), nil
}

func createKeyboardNavigator(directories []Directory, files []File) tgbotapi.ReplyKeyboardMarkup {
	rows := [][]tgbotapi.KeyboardButton{}

	createNewFolderBtn := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Create new folder"),
	)

	rows = append(rows, createNewFolderBtn)

	for _, file := range files {
		current := tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(file.Name),
		)
		rows = append(rows, current)
	}

	for _, dir := range directories {
		current := tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(dir.Name),
		)
		rows = append(rows, current)
	}

	goBackBtn := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("../"),
	)
	rows = append(rows, goBackBtn)
	
	return tgbotapi.NewOneTimeReplyKeyboard(rows...)
}

func proccessError(err error) string {
	switch err.Error() {
		case "directory or user not found!": 
			return "You need create account before sending files! For it send me /start"
		case "user already exists":
			return "You already have account!"
		case "can`t proccess message: file already exists in folder":
			return "This file already in folder"
		default: 
			return "Sorry, i can`t understand what u want to do :("
	}
}










// func createTextKeyboard(directories []Directory, files []File) string {
// 	var rows []string
	
// 	for _, file := range files {
// 		rows = append(rows, file.Name)
// 	}

// 	for _, dir := range directories {
// 		rows = append(rows, dir.Name)
// 	}

// 	rows = append(rows, "../")
// 	return strings.Join(rows, "\n")
// }




