package telegram

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

func createEmptyKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.NewRemoveKeyboard(true)
}

func createKeyboardNavigator(directories []Directory, files []File) tgbotapi.ReplyKeyboardMarkup {
	rows := [][]tgbotapi.KeyboardButton{}

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