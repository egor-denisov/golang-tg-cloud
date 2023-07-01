package telegram

import (
	"errors"
	"log"
	"main/db"
	"main/lib/e"

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
			cl.sendMessage(update.Message.Chat.ID, Content{Text: "Sorry, i can`t understand what u want to do :("})
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
				Id: msg.From.ID,
				UserName: msg.From.UserName,
				FirstName: msg.From.FirstName,
				LastName: msg.From.LastName,
			}
			_, err := createNewUser(cl.db, userInfo)
			if err.Error() == "user already exists" {
				content.Text = "You already have account!"
			}else if err != nil {	
				return err
			}
		case "search":
			content.Text = "Input search string."
		case "showAll":
			content.Text = "Show all files in current directory"
		case "mainFolder":
			content.Text = "Jump to main directory"
		default:
			content.Text = "I don't know that command"
		}
	}else if msg.Photo != nil {
		content = Content{
			Photo: msg.Photo,
		}
	}else if msg.Document != nil {
		content = Content{
			Document: msg.Document,
		}
	}else{
		return errors.New("Unknown type of message")
	}

	if err := cl.sendMessage(msg.Chat.ID, content); err != nil{
		return err
	}
	return nil
}


func (cl *TgClient) sendMessage(chatID int64, content Content) (err error) {
	defer func() { err = e.WrapIfErr("can`t send message", err) }()

	var msg tgbotapi.Chattable

	if content.Document != nil {
		msg = tgbotapi.NewDocument(chatID, tgbotapi.FileID(content.Document.FileID))
	}else if content.Photo != nil {
		msg = tgbotapi.NewPhoto(chatID, tgbotapi.FileID(content.Photo[0].FileID))
	}else if content.Text != "" {
		msg = tgbotapi.NewMessage(chatID, content.Text)
	}else {
		return errors.New("can`t find required type of message")
	}

	if _, err := cl.bot.Send(msg); err != nil {
		return err
	}
	return nil
}
