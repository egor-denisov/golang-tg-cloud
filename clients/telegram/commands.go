package telegram

import (
	"fmt"
	"main/lib/h"
	. "main/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Function for proccessing commands in telegram bot
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

// Action for replying to the 'start' command
func startCommand(cl *TgClient, msg *tgbotapi.Message) error {
	// Creating content for replying
	content := Content{
		Text: "Hello! Send me a file or create a new folder",
	}
	// Creating instance of a new user
	userInfo := User{
		ChatId:    int(msg.Chat.ID),
		UserId:    int(msg.From.ID),
		UserName:  msg.From.UserName,
		FirstName: msg.From.FirstName,
		LastName:  msg.From.LastName,
	}
	// Adding a new user to the database
	if _, err := cl.db.CreateNewUser(userInfo); err != nil {
		return err
	}
	// Sending replying message
	return cl.sendMedia(msg.Chat.ID, content)
}

// Action for replying to the 'help' command
func helpCommand(cl *TgClient, msg *tgbotapi.Message) error {
	// Creating content for replying
	content := Content{
		Text: "You can use /search for search file. Also you can see all files in current directory with /showAll.",
	}
	// Sending replying message
	return cl.sendMedia(msg.Chat.ID, content)
}

// Action for replying to the 'create folder' command
func createFolderCommand(cl *TgClient, msg *tgbotapi.Message) error {
	// Creating instance of a new directory
	directory := Directory{
		Name: msg.CommandArguments(),
	}
	// Checking the correctness for a directory name
	if !h.IsValidName(directory.Name) {
		return fmt.Errorf("wrong folder name")
	}
	// Adding prefix for highlighting directory
	directory.Name = "./" + directory.Name
	// Creating a new directory in database
	if _, err := cl.db.CreateNewDirectory(msg.From.ID, directory); err != nil {	
		return err
	}
	// Making response for user after creating directory
	return cl.makeReplyAfterCreatingDirectory(msg.From.ID, directory.Name)
}

// Action for replying to the 'reset' command
func resetCommand(cl *TgClient, msg *tgbotapi.Message) error {
	// Creating content for replying
	content := Content{
		Text: "You successfully reset data",
		Keyboard: createEmptyKeyboard(),
	}
	// Resetting user data from database
	if err := cl.db.ResetUserData(msg.From.ID); err != nil {
		return err
	}
	// Sending replying message
	return cl.sendMedia(msg.Chat.ID, content)
}

// Action for replying to an unknown command
func unknownCommand(cl *TgClient, msg *tgbotapi.Message) error {
	// Creating content for replying
	content := Content{
		Text: "I don't know this command",
	}
	// Sending replying message
	return cl.sendMedia(msg.Chat.ID, content)
}