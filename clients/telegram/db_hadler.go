package telegram

import (
	"errors"
	"fmt"
	"main/db"
)

func createNewDirectory(db db.DataBase, directory Directory) (int64, error) {
	id, err := db.Insert("insert into directories (Name) values ($1) returning Id", directory.Name)
	return id, err
}

func createNewUser(db db.DataBase, user User) (int64, error) {
	exist, err := userExists(db, user.UserId)
	if err != nil {
		return -1, err
	}
	if exist {
		return -1, errors.New("user already exists")
	}

	directoryId, err := createNewDirectory(db, Directory{Name : "/"})
	if err != nil {
		return -1, err
	}
	id, err := db.Insert("insert into users (Username, UserID, FirstName, LastName, CurrentDirectory) values ($1, $2, $3, $4, $5) returning Id", 
		user.UserName, user.UserId, user.FirstName, user.LastName, directoryId)
	return id, err
}

func userExists(db db.DataBase, userId int64) (bool, error) {
	req := fmt.Sprintf("select id from users where userId = %d", userId)
	exist, err := db.Exist(req)
	if err != nil {
		return true, err
	}
	return exist, nil
}