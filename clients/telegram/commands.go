package telegram

import (
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
		Text: "Send the files to add in folder and use command /create_folder to create new folder",
	}
	// Sending replying message
	return cl.sendMedia(msg.Chat.ID, content)
}

// Action for replying to the 'create folder' command
func createFolderCommand(cl *TgClient, msg *tgbotapi.Message) error {
	// Getting current directory and setting as parent directory
	currentDirectoryId, err := cl.db.GetCurrentDirectory(msg.From.ID)
	if err != nil {
		return err
	}
	currentDirectory, err := cl.db.GetDirectory(currentDirectoryId)
	if err != nil {
		return err
	}
	// Creating instance of a new directory
	directory := Directory{
		Name: msg.CommandArguments(),
		ParentId: currentDirectoryId,
		UserId: int(msg.From.ID),
		Path: currentDirectory.Path + "/" + msg.CommandArguments(),
	}
	// Creating a new directory in database
	if _, err := cl.db.CreateNewDirectory(directory); err != nil {	
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