package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/db"
	"strconv"
	"strings"
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

func addFileToDirectory(db db.DataBase, directoryId string, fileId string) error {
	req := fmt.Sprintf("update directories set files = array_append(files, %s) where id = %s", fileId, directoryId)
	err := db.MakeQuery(req)
	if err != nil {
		return err
	}
	return nil
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
	id, err := db.SelectRow(req)
	if err != nil {
		return true, err
	}
	return id != "", nil
}

func getFileId(db db.DataBase, fileId string) (string, error) {
	req := fmt.Sprintf("select id from files where fileId = %s", fileId)
	id, err := db.SelectRow(req)
	if err != nil {
		return "", err
	}
	return id, nil
}

func getCurrentDirectory(db db.DataBase, userId string) (string, error) {
	req := fmt.Sprintf("select currentDirectory from users where userId = %s", userId)
	directoryId, err := db.SelectRow(req)
	if err != nil {
		return "", err
	}
	if directoryId == "" {
		return "", errors.New("directory or user not found!")
	}
	return directoryId, nil
}

func getDirectory(db db.DataBase, id int) (Directory, error) {
	d := Directory{}
	req := fmt.Sprintf("select id, name, size from directories where id = %d", id)
	rows, err := db.Select(req)
	if err != nil {
		return d, err
	}
	rows.Next()
	err = rows.Scan(&d.Id, &d.Name, &d.Size)
	if err != nil {
		return d, err
	}
	return d, nil
}

func getFile(db db.DataBase, id int) (File, error) {
	f := File{}
	req := fmt.Sprintf("select * from files where id = %d", id)
	rows, err := db.Select(req)
	if err != nil {
		return f, err
	}
	rows.Next()
	err = rows.Scan(&f.Id, &f.Name, &f.FileId, &f.FileUniqueId, &f.FileSize)
	if err != nil {
		return f, err
	}
	return f, nil
}

func getAvailableDirectories(db db.DataBase, userId string) ([]Directory, error) {
	var res []Directory

	directoryId, err := getCurrentDirectory(db, userId)
	if err != nil {
		return nil, err
	}

	arr, err := getIdsArray(db, directoryId, "directories")

	for _, id := range arr {
		d, err := getDirectory(db, id)
		if err != nil {
			return res, err
		}
		res = append(res, d)
	}

	return res, nil
}

func getAvailableFiles(db db.DataBase, userId string) ([]File, error) {
	var res []File

	directoryId, err := getCurrentDirectory(db, userId)
	if err != nil {
		return nil, err
	}

	arr, err := getIdsArray(db, directoryId, "files")

	for _, id := range arr {
		d, err := getFile(db, id)
		if err != nil {
			return res, err
		}
		res = append(res, d)
	}

	return res, nil
}

func getIdsArray(db db.DataBase, id string, name string) ([]int, error) {
	req := fmt.Sprintf("select %s from directories where id = %s", name, id)
	idsStr, err := db.SelectRow(req)
	if err != nil {
		return nil, err
	}
	ids, err := parseIds(idsStr)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func parseIds(jsonBuffer string) ([]int, error) {
	ids := []int{}

	jsonBuffer = strings.Replace(jsonBuffer, "{", "[", -1)
	jsonBuffer = strings.Replace(jsonBuffer, "}", "]", -1)

    err := json.Unmarshal([]byte(jsonBuffer), &ids)
    if err != nil {
        return nil, err
    }

    return ids, nil
}