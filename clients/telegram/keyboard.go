package telegram

import (
	. "main/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Creating a new bot keyboard for the current user
func (cl *TgClient) instantiateKeyboardNavigator(userID int64) (tgbotapi.ReplyKeyboardMarkup, error) {
	// Getting available files from the database
	files, err := cl.db.GetAvailableFiles(userID)
	if err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}, err
	}
	// Getting available directories from the database
	directories, err := cl.db.GetAvailableDirectories(userID)
	if err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}, err
	}
	// Return keyboard instance
	return createKeyboardNavigator(directories, files), nil
}
// Creating an empty instance of keyboard
func createEmptyKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.NewRemoveKeyboard(true)
}
// Creating a keyboard instance with files and directories
func createKeyboardNavigator(directories []Directory, files []File) tgbotapi.ReplyKeyboardMarkup {
	// Inicializing empty keyboard rows
	rows := [][]tgbotapi.KeyboardButton{}
	// Iteration over the files and getting its names
	for _, file := range files {
		current := tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(file.Name),
		)
		rows = append(rows, current)
	}
	// Iteration over the directories and getting its names
	for _, dir := range directories {
		current := tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(dir.Name),
		)
		rows = append(rows, current)
	}
	// Creating 'go back' button and adding it to the list of rows
	goBackBtn := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("../"),
	)
	rows = append(rows, goBackBtn)
	
	return tgbotapi.NewOneTimeReplyKeyboard(rows...)
}