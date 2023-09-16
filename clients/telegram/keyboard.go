package telegram

import (
	. "main/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Creating a new bot keyboard for the current user
func (cl *TgClient) instantiateKeyboardNavigator(userId int64) (tgbotapi.ReplyKeyboardMarkup, error) {
	// Getting current directory
	directoryId, err := cl.db.GetCurrentDirectory(userId)
	if err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}, err
	}
	// Getting available items from the database
	items, err := cl.db.GetAvailableItemsInDirectory(userId, directoryId)
	if err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}, err
	}
	// Return keyboard instance
	return createKeyboardNavigator(items.Directories, items.Files), nil
}
// Creating an empty instance of keyboard
func createEmptyKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.NewRemoveKeyboard(true)
}
// Creating a keyboard instance with files and directories
func createKeyboardNavigator(directories []Directory, files []File) tgbotapi.ReplyKeyboardMarkup {
	// Inicializing empty keyboard rows
	rows := [][]tgbotapi.KeyboardButton{}
	// Iteration over the directories and getting its names
	for _, dir := range directories {
		if dir.Name != "../" {
			dir.Name = "./" + dir.Name
		}
		current := tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(dir.Name),
		)
		rows = append(rows, current)
	}
	// Iteration over the files and getting its names
	for _, file := range files {
		current := tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(file.Name),
		)
		rows = append(rows, current)
	}
	
	return tgbotapi.NewOneTimeReplyKeyboard(rows...)
}