package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/db"
	"main/lib/h"
	"strconv"
	"strings"
)

func createRootDirectory(db db.DataBase, directory Directory) (string, error) {
	id, err := db.Insert("insert into directories (Name, UserId) values ($1, $2) returning Id", directory.Name, directory.UserId)
	return id, err
}

func createNewDirectory(db db.DataBase, userId int64, directory Directory) (string, error) {

	currentDirectoryId, err := getCurrentDirectory(db, userId)
	if err != nil {
		return "", err
	}
	directory.ParentId = currentDirectoryId

	existence, err := folderExists(db, userId, currentDirectoryId, directory.Name)
	if err != nil {
		return "", err
	}
	if existence {
		return "", fmt.Errorf("directory with this name already exists in current folder")
	}
	
	id, err := db.Insert("insert into directories (ParentId, Name, UserId) values ($1, $2, $3) returning Id", 
		currentDirectoryId, directory.Name, userId)
	if err != nil {
		return "", err
	}

	newDirectoryId, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}

	return id, addNewDirectoryToDirectory(db, currentDirectoryId, newDirectoryId)
}

func createNewFile(db db.DataBase, userId int64, file File) (string, error) {
	directoryId, err := getCurrentDirectory(db, userId)
	if err != nil {
		return "", err
	}

	id, err := getIdOfFileByUniqueId(db, file.FileUniqueId)
	if err != nil {
		return "", err
	}

	if id > 0 {
		file.Id = id
		currentFiles, err := getIdsArray(db, directoryId, "files")
		if err != nil {
			return "", err
		}

		if h.Contains(currentFiles, file.Id) {
			return "", fmt.Errorf("file already exists in folder")
		}
	}else{
		idStr, err := db.Insert("insert into files (Name, FileId, FileUniqueId, FileSize) values ($1, $2, $3, $4) returning Id", 
		file.Name, file.FileId, file.FileUniqueId, file.FileSize)
		if err != nil {
			return "", err
		}
		file.Id, err = strconv.Atoi(idStr)
		if err != nil {
			return "", err
		}
	}

	if err := addFileToDirectory(db, directoryId, file); err != nil {
		return "", err
	}
	return strconv.Itoa(file.Id), nil
}

func addFileToDirectory(db db.DataBase, directoryId int, file File) error {
	req := fmt.Sprintf("update directories set files = array_append(files, %d), size = size + %d where id = %d", file.Id, file.FileSize, directoryId)
	err := db.MakeQuery(req)
	if err != nil {
		return err
	}
	return nil
}

func folderExists(db db.DataBase, userId int64, currentDirectoryId int, directoryName string) (bool, error) {
	currentDirectories, err := getNamesArray(db, currentDirectoryId, "directories")
	if err != nil {
		return false, err
	}
	return h.Contains(currentDirectories, directoryName), nil
}

func addNewDirectoryToDirectory(db db.DataBase, currentDirectoryId int, newDirectoryId int) error {
	req := fmt.Sprintf("update directories set directories = array_append(directories, %d) where id = %d", newDirectoryId, currentDirectoryId)
	
	if err := db.MakeQuery(req); err != nil {
		return err
	}
	return nil
}

func jumpToDirectory(db db.DataBase, userId int64, directoryId int) error {
	req := fmt.Sprintf("update users set currentDirectory = %d where userId = %d", directoryId, userId)
	
	if err := db.MakeQuery(req); err != nil {
		return err
	}
	return nil
}

func getParentDirectory(db db.DataBase, directoryId int) (Directory, error) {
	currentDirectory, err := getDirectory(db, directoryId)
	if err != nil{
		return Directory{}, err
	}
	
	if currentDirectory.ParentId == -1 {
		return Directory{}, fmt.Errorf("this directory is root")
	}

	return getDirectory(db, currentDirectory.ParentId)
}

func createNewUser(db db.DataBase, user User) (string, error) {
	exist, err := userExists(db, int64(user.UserId))
	if err != nil {
		return "", err
	}
	if exist {
		return "", errors.New("user already exists")
	}

	directoryId, err := createRootDirectory(db, Directory{Name : "/", UserId: user.UserId})
	if err != nil {
		return "", err
	}
	id, err := db.Insert("insert into users (Username, UserID, FirstName, LastName, CurrentDirectory) values ($1, $2, $3, $4, $5) returning Id", 
		user.UserName, user.UserId, user.FirstName, user.LastName, directoryId)
	return id, err
}

func userExists(db db.DataBase, userId int64) (bool, error) {
	req := fmt.Sprintf("select id from users where userId = %d", userId)
	id, err := db.SelectRow(req)
	if err != nil {
		return true, err
	}
	return id != "", nil
}

func getIdOfFileByUniqueId(db db.DataBase, fileUniqueId string) (int, error) {
	req := fmt.Sprintf("select id from files where fileUniqueId = '%s'", fileUniqueId)
	id, err := db.SelectRow(req)
	if err != nil {
		return -1, err
	}
	if id == "" {
		return -1, nil
	}
	return strconv.Atoi(id)
}

func getCurrentDirectory(db db.DataBase, userId int64) (int, error) {
	req := fmt.Sprintf("select currentDirectory from users where userId = %d", userId)
	directoryId, err := db.SelectRow(req)
	if err != nil {
		return -1, err
	}
	if directoryId == "" {
		return -1, errors.New("directory or user not found!")
	}
	return strconv.Atoi(directoryId)
}

func getDirectory(db db.DataBase, id int) (Directory, error) {
	d := Directory{}
	req := fmt.Sprintf("select id, parentId, name, userId, size from directories where id = %d", id)
	rows, err := db.Select(req)
	if err != nil {
		return d, err
	}
	
	rows.Next()
	if err := rows.Scan(&d.Id, &d.ParentId, &d.Name, &d.UserId, &d.Size); err != nil {
		return d, err
	}
	return d, nil
}

func getFileById(db db.DataBase, id int) (File, error) {
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

func getFileByName(db db.DataBase, name string) (File, error) { 
	// This function works wrong 
	// bcs it finds all files with requested name 
	// including when user ain't host
	// solution -> add into database column 'hosts' with type []int(better use id`s users)
	f := File{}
	req := fmt.Sprintf("select * from files where name = '%s'", name)
	rows, err := db.Select(req)
	if err != nil {
		return f, err
	}

	rows.Next()
	
	if err = rows.Scan(&f.Id, &f.Name, &f.FileId, &f.FileUniqueId, &f.FileSize); err != nil {
		return f, err
	}
	return f, nil
}

func getAvailableDirectories(db db.DataBase, userId int64) ([]Directory, error) {
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

func getAvailableFiles(db db.DataBase, userId int64) ([]File, error) {
	var res []File

	directoryId, err := getCurrentDirectory(db, userId)
	if err != nil {
		return nil, err
	}

	arr, err := getIdsArray(db, directoryId, "files")

	for _, id := range arr {
		d, err := getFileById(db, id)
		if err != nil {
			return res, err
		}
		res = append(res, d)
	}

	return res, nil
}

func getIdsArray(db db.DataBase, directoryId int, name string) ([]int, error) {
	req := fmt.Sprintf("select %s from directories where id = %d", name, directoryId)
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

func getNamesArray(db db.DataBase, directoryId int, name string) ([]string, error) {
	req := fmt.Sprintf("select %s from directories where id = %d", name, directoryId)
	idsStr, err := db.SelectRow(req)
	if err != nil {
		return nil, err
	}
	ids, err := parseIds(idsStr)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	req = fmt.Sprintf("select name from directories where id in ( %s )", strings.Join(h.IntArrayToStrArray(ids), ", "))
	rows, err := db.Select(req)
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

func resetUserData(db db.DataBase, userId int64) error {
	directoryId, err := createRootDirectory(db, Directory{Name : "/", UserId: int(userId)})
	if err != nil {
		return err
	}

	req := fmt.Sprintf("update users set currentDirectory = %s where userId=%d", directoryId, userId)
	if err := db.MakeQuery(req); err != nil {
		return err
	}

	req = fmt.Sprintf("delete from directories where userId=%d and id != %s", userId, directoryId)
	if err := db.MakeQuery(req); err != nil {
		return err
	}

	return nil
}

func parseIds(jsonBuffer string) ([]int, error) {
	ids := []int{}
	if len(jsonBuffer) == 0 {
		return ids, nil
	}
	jsonBuffer = strings.Replace(jsonBuffer, "{", "[", -1)
	jsonBuffer = strings.Replace(jsonBuffer, "}", "]", -1)

    err := json.Unmarshal([]byte(jsonBuffer), &ids)
    if err != nil {
        return nil, err
    }

    return ids, nil
}
