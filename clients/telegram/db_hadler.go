package telegram

import (
	"errors"
	"fmt"
	"main/db"
	"strconv"
)

func createNewDirectory(db db.DataBase, directory Directory) (string, error) {
	id, err := db.Insert("insert into directories (Name) values ($1) returning Id", directory.Name)
	return id, err
}

func createNewFile(db db.DataBase, userId int, file File) (string, error) {
	fileId := ""
	id, err := getFileId(db, file.FileId)
	if err != nil {
		return "", err
	}
	if id != "" {
		fileId = id
	}else{
		fileId, err = db.Insert("insert into files (Name, FileId, FileUniqueId, FileSize) values ($1, $2, $3, $4) returning Id", 
		file.Name, file.FileId, file.FileUniqueId, file.FileSize)
		if err != nil {
			return "", err
		}
	}

	directoryId, err := getCurrentDirectory(db, strconv.Itoa(userId))
	if err != nil {
		return "", err
	}
	err = addFileToDirectory(db, directoryId, fileId)
	if err != nil {
		return "", err
	}
	return id, err
}

func createNewUser(db db.DataBase, user User) (string, error) {
	exist, err := userExists(db, strconv.Itoa(user.UserId))
	if err != nil {
		return "", err
	}
	if exist {
		return "", errors.New("user already exists")
	}

	directoryId, err := createNewDirectory(db, Directory{Name : "/"})
	if err != nil {
		return "", err
	}
	id, err := db.Insert("insert into users (Username, UserID, FirstName, LastName, CurrentDirectory) values ($1, $2, $3, $4, $5) returning Id", 
		user.UserName, user.UserId, user.FirstName, user.LastName, directoryId)
	return id, err
}

func userExists(db db.DataBase, userId string) (bool, error) {
	req := fmt.Sprintf("select id from users where userId = %s", userId)
	id, err := db.RawSelect(req)
	if err != nil {
		return true, err
	}
	return id != "", nil
}

func getFileId(db db.DataBase, fileId string) (string, error) {
	req := fmt.Sprintf("select id from files where fileId = %s", fileId)
	id, err := db.RawSelect(req)
	if err != nil {
		return "", err
	}
	return id, nil
}

func getCurrentDirectory(db db.DataBase, userId string) (string, error) {
	req := fmt.Sprintf("select currentDirectory from users where userId = %s", userId)
	directoryId, err := db.RawSelect(req)
	if err != nil {
		return "", err
	}
	if directoryId == "" {
		return "", errors.New("directory or user not found!")
	}
	return directoryId, nil
}

func addFileToDirectory(db db.DataBase, directoryId string, fileId string) error {
	req := fmt.Sprintf("update directories set files = array_append(files, %s) where id = %s", fileId, directoryId)
	err := db.MakeQuery(req)
	if err != nil {
		return err
	}
	return nil
}