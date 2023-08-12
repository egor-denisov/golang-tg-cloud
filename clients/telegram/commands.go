package telegram

import (
	"fmt"
	"main/lib/h"
	. "main/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (cl *TgClient) proccessCommand(msg *tgbotapi.Message) error{
	switch msg.Command() {
		case "help":
			return helpCommand(cl, msg)
		case "start":
			return startCommand(cl, msg)
		case "create_folder":
			return createFolderCommand(cl, msg)
		case "reset":
			return resetCommand(cl, msg)
		default:
			return unknownCommand(cl, msg)
	}
}

func startCommand(cl *TgClient, msg *tgbotapi.Message) error {
	var content Content
	content.Text = "Hello! Send me a file or create a new folder"
	userInfo := User{
		ChatId:    int(msg.Chat.ID),
		UserId:    int(msg.From.ID),
		UserName:  msg.From.UserName,
		FirstName: msg.From.FirstName,
		LastName:  msg.From.LastName,
	}

	if _, err := cl.db.CreateNewUser(userInfo); err != nil {
		return err
	}
	
	return cl.sendMedia(msg.Chat.ID, content)
}

func helpCommand(cl *TgClient, msg *tgbotapi.Message) error {
	return cl.sendMedia(msg.Chat.ID, Content{Text: "You can use /search for search file. Also you can see all files in current directory with /showAll."})
}

func createFolderCommand(cl *TgClient, msg *tgbotapi.Message) error {
	directory := Directory{
		Name: msg.CommandArguments(),
	}
	if !h.IsValidName(directory.Name) {
		return fmt.Errorf("wrong folder name")
	}
	directory.Name = "./" + directory.Name
	
	if _, err := cl.db.CreateNewDirectory(msg.From.ID, directory); err != nil {	
		return err
	}
	
	return cl.makeReplyAfterCreatingDirectory(msg.From.ID, directory.Name)
}

func resetCommand(cl *TgClient, msg *tgbotapi.Message) error {
	content := Content{
		Text: "You successfully reset data",
		Keyboard: createEmptyKeyboard(),
	}

	if err := cl.db.ResetUserData(msg.From.ID); err != nil {
		return err
	}

	return cl.sendMedia(msg.Chat.ID, content)
}

func unknownCommand(cl *TgClient, msg *tgbotapi.Message) error {
	return cl.sendMedia(msg.Chat.ID, Content{Text: "I don't know this command"})
}