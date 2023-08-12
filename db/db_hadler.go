package db

import (
	"errors"
	"fmt"
	"main/lib/h"
	. "main/types"
	"strconv"
	"strings"
)

func (db *DataBase) CreateRootDirectory(directory Directory) (string, error) {
	return db.insert("insert into directories (Name, UserId) values ($1, $2) returning Id", 
		directory.Name, directory.UserId)
}

func (db *DataBase) CreateNewDirectory(userId int64, directory Directory) (string, error) {
	currentDirectoryId, err := db.GetCurrentDirectory(userId)
	if err != nil {
		return "", err
	}
	directory.ParentId = currentDirectoryId

	existence, err := db.FolderExists(userId, currentDirectoryId, directory.Name)
	if err != nil {
		return "", err
	}
	if existence {
		return "", fmt.Errorf("directory with this name already exists in current folder")
	}
	
	id, err := db.insert("insert into directories (ParentId, Name, UserId) values ($1, $2, $3) returning Id", 
		currentDirectoryId, directory.Name, userId)
	if err != nil {
		return "", err
	}

	newDirectoryId, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}

	return id, db.AddNewDirectoryToDirectory(currentDirectoryId, newDirectoryId)
}

func (db *DataBase) CreateNewFile(userId int64, file File) (string, error) {
	directoryId, err := db.GetCurrentDirectory(userId)
	if err != nil {
		return "", err
	}

	id, err := db.GetIdOfFileByUniqueId(file.FileUniqueId)
	if err != nil {
		return "", err
	}

	if id > 0 {
		file.Id = id
		currentFiles, err := db.GetIdsArray(directoryId, "files")
		if err != nil {
			return "", err
		}

		if h.Contains(currentFiles, file.Id) {
			return "", fmt.Errorf("file already exists in folder")
		}
	}else{
		idStr, err := db.insert("insert into files (Name, FileId, FileUniqueId, FileSize) values ($1, $2, $3, $4) returning Id", 
			file.Name, file.FileId, file.FileUniqueId, file.FileSize)
		if err != nil {
			return "", err
		}
		file.Id, err = strconv.Atoi(idStr)
		if err != nil {
			return "", err
		}
	}

	return strconv.Itoa(file.Id), db.AddFileToDirectory(directoryId, file)
}

func (db *DataBase) AddFileToDirectory(directoryId int, file File) error {
	req := fmt.Sprintf("update directories set files = array_append(files, %d), size = size + %d where id = %d", 
		file.Id, file.FileSize, directoryId)
	return db.makeQuery(req)
}

func (db *DataBase) FolderExists(userId int64, currentDirectoryId int, directoryName string) (bool, error) {
	currentDirectories, err := db.GetNamesArray(currentDirectoryId, "directories")
	return h.Contains(currentDirectories, directoryName), err
}

func (db *DataBase) AddNewDirectoryToDirectory(currentDirectoryId int, newDirectoryId int) error {
	req := fmt.Sprintf("update directories set directories = array_append(directories, %d) where id = %d", 
		newDirectoryId, currentDirectoryId)
	return db.makeQuery(req)
}

func (db *DataBase) JumpToDirectory(userId int64, directoryId int) error {
	req := fmt.Sprintf("update users set currentDirectory = %d where userId = %d", directoryId, userId)
	return db.makeQuery(req)
}

func (db *DataBase) GetParentDirectory(directoryId int) (Directory, error) {
	currentDirectory, err := db.GetDirectory(directoryId)
	if err != nil{
		return Directory{}, err
	}
	
	if currentDirectory.ParentId == -1 {
		return Directory{}, fmt.Errorf("this directory is root")
	}

	return db.GetDirectory(currentDirectory.ParentId)
}

func (db *DataBase) CreateNewUser(user User) (string, error) {
	exist, err := db.UserExists(int64(user.UserId))
	if err != nil {
		return "", err
	}
	if exist {
		return "", errors.New("user already exists")
	}

	directoryId, err := db.CreateRootDirectory(Directory{Name : "/", UserId: user.UserId})
	if err != nil {
		return "", err
	}
	
	return db.insert("insert into users (Username, ChatID, UserID, FirstName, LastName, CurrentDirectory) values ($1, $2, $3, $4, $5, $6) returning Id", 
		user.UserName, user.ChatId, user.UserId, user.FirstName, user.LastName, directoryId)
}

func (db *DataBase) UserExists(userId int64) (bool, error) {
	req := fmt.Sprintf("select id from users where userId = %d", userId)
	id, err := db.selectRow(req)
	return id != "", err
}

func (db *DataBase) GetIdOfFileByUniqueId(fileUniqueId string) (int, error) {
	req := fmt.Sprintf("select id from files where fileUniqueId = '%s'", fileUniqueId)
	id, err := db.selectRow(req)
	if err != nil {
		return -1, err
	}
	if id == "" {
		return -1, nil
	}
	return strconv.Atoi(id)
}

func (db *DataBase) GetCurrentDirectory(userId int64) (int, error) {
	req := fmt.Sprintf("select currentDirectory from users where userId = %d", userId)
	directoryId, err := db.selectRow(req)
	if err != nil {
		return -1, err
	}
	if directoryId == "" {
		return -1, errors.New("directory or user not found!")
	}
	return strconv.Atoi(directoryId)
}

func (db *DataBase) GetDirectory(id int) (Directory, error) {
	d := Directory{}
	req := fmt.Sprintf("select id, parentId, name, userId, size from directories where id = %d", id)
	rows, err := db.selectRows(req)
	if err != nil {
		return d, err
	}
	
	rows.Next()
	err = rows.Scan(&d.Id, &d.ParentId, &d.Name, &d.UserId, &d.Size)
	return d, err
}

func (db *DataBase) GetFileById(id int) (File, error) {
	f := File{}
	req := fmt.Sprintf("select * from files where id = %d", id)
	rows, err := db.selectRows(req)
	if err != nil {
		return f, err
	}

	rows.Next()
	err = rows.Scan(&f.Id, &f.Name, &f.FileId, &f.FileUniqueId, &f.FileSize)
	return f, err
}

func (db *DataBase) GetAvailableDirectories(userId int64) ([]Directory, error) {
	var res []Directory

	directoryId, err := db.GetCurrentDirectory(userId)
	if err != nil {
		return nil, err
	}

	arr, err := db.GetIdsArray(directoryId, "directories")
	if err != nil {
		return nil, err
	}

	for _, id := range arr {
		d, err := db.GetDirectory(id)
		if err != nil {
			return res, err
		}
		res = append(res, d)
	}

	return res, nil
}

func (db *DataBase) GetAvailableFiles(userId int64) ([]File, error) {
	var res []File

	directoryId, err := db.GetCurrentDirectory(userId)
	if err != nil {
		return nil, err
	}

	arr, err := db.GetIdsArray(directoryId, "files")
	if err != nil {
		return nil, err
	}

	for _, id := range arr {
		d, err := db.GetFileById(id)
		if err != nil {
			return res, err
		}
		res = append(res, d)
	}

	return res, nil
}

func (db *DataBase) GetIdsArray(directoryId int, name string) ([]int, error) {
	req := fmt.Sprintf("select %s from directories where id = %d", name, directoryId)
	idsStr, err := db.selectRow(req)
	if err != nil {
		return nil, err
	}

	return h.ParseIds(idsStr)
}

func (db *DataBase) GetNamesArray(directoryId int, name string) ([]string, error) {
	req := fmt.Sprintf("select %s from directories where id = %d", name, directoryId)
	idsStr, err := db.selectRow(req)
	if err != nil {
		return nil, err
	}

	ids, err := h.ParseIds(idsStr)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	req = fmt.Sprintf("select name from directories where id in ( %s )", strings.Join(h.IntArrayToStrArray(ids), ", "))
	rows, err := db.selectRows(req)
	if err != nil {
		return nil, err
	}

	var res []string
	for rows.Next() {
		current := ""
		err = rows.Scan(&current)
		if err != nil {
			return nil, err
		}
		res = append(res, current)
	}

	return res, nil
}

func (db *DataBase) ResetUserData(userId int64) error {
	directoryId, err := db.CreateRootDirectory(Directory{Name : "/", UserId: int(userId)})
	if err != nil {
		return err
	}

	req := fmt.Sprintf("update users set currentDirectory = %s where userId=%d", directoryId, userId)
	if err := db.makeQuery(req); err != nil {
		return err
	}

	req = fmt.Sprintf("delete from directories where userId=%d and id != %s", userId, directoryId)

	return db.makeQuery(req)
}
