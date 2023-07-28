package telegram

import (
	"fmt"
	"main/lib/h"

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
		UserId:    int(msg.From.ID),
		UserName:  msg.From.UserName,
		FirstName: msg.From.FirstName,
		LastName:  msg.From.LastName,
	}

	if _, err := createNewUser(cl.db, userInfo); err != nil {
		return err
	}

	if err := cl.sendMedia(msg.Chat.ID, content); err != nil {
		return err
	}
	return nil
}

func helpCommand(cl *TgClient, msg *tgbotapi.Message) error {
	var content Content
	content.Text = "You can use /search for search file. Also you can see all files in current directory with /showAll."

	if err := cl.sendMedia(msg.Chat.ID, content); err != nil {
		return err
	}
	return nil
}

func createFolderCommand(cl *TgClient, msg *tgbotapi.Message) error {
	directory := Directory{
		Name: msg.CommandArguments(),
	}
	if !h.IsValidName(directory.Name) {
		return fmt.Errorf("wrong folder name")
	}
	directory.Name = "./" + directory.Name
	
	if _, err := createNewDirectory(cl.db, msg.From.ID, directory); err != nil {	
		return err
	}
	
	if err := cl.makeReplyAfterCreatingDirectory(msg.From.ID, directory.Name); err != nil {	
		return err
	}
	return nil
}

func resetCommand(cl *TgClient, msg *tgbotapi.Message) error {
	var content Content
	content.Text = "You successfully reset data"
	content.Keyboard = createEmptyKeyboard()

	if err := resetUserData(cl.db, msg.From.ID); err != nil {
		return err
	}

	if err := cl.sendMedia(msg.Chat.ID, content); err != nil {
		return err
	}
	return nil
}

func unknownCommand(cl *TgClient, msg *tgbotapi.Message) error {
	var content Content
	content.Text = "I don't know this command"

	if err := cl.sendMedia(msg.Chat.ID, content); err != nil {
		return err
	}
	return nil
}