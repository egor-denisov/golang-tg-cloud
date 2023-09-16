package db

import (
	"errors"
	"fmt"
	"log"
	"main/lib/h"
	. "main/types"
	"strconv"
	"strings"
)

//Function for creating root directory
func (db *DataBase) CreateRootDirectory(directory Directory) (string, error) {
	return db.insert("insert into directories (Name, ParentId, UserId, Path, Created) values ($1, -1, $2, $1, now()) returning Id", 
		directory.Name, directory.UserId)
}
//Function for creating new directory
func (db *DataBase) CreateNewDirectory(directory Directory) (string, error) {
	// Checking the correctness for a directory name
	if !h.IsValidName(directory.Name) {
		return "", fmt.Errorf("wrong folder name")
	}
	// Checking existence of folder with current name
	existence, err := db.FolderExists(int64(directory.UserId), directory.ParentId, directory.Name)
	if err != nil {
		return "", err
	}
	// Return error if directory exists
	if existence {
		return "", fmt.Errorf("directory with this name already exists in current folder")
	}
	// Getting parent directory data
	currentDirectory, err := db.GetDirectory(directory.ParentId)
	if err != nil {
		return "", err
	}
	// Inserting new directory into database
	id, err := db.insert("insert into directories (Name, ParentId, UserId, Path, Created) values ($1, $2, $3, $4, now()) returning Id", 
		directory.Name, directory.ParentId, directory.UserId, currentDirectory.Path + directory.Name + "/")
	if err != nil {
		return "", err
	}
	// Converting id to integer 
	newDirectoryId, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	// Returning id of new directory and result of adding new directory
	return id, db.AddNewDirectoryToDirectory(directory.ParentId, newDirectoryId)
}
// Function for creating a new file
func (db *DataBase) CreateNewFile(userId int64, directoryId int, file File) (int, error) {
	// Getting id by unique id
	id, err := db.GetIdOfFileByUniqueId(file.FileUniqueId)
	if err != nil {
		return -1, err
	}
	// If file already exists then error is returned 
	if id > 0 {
		file.Id = id
		// Getting files id`s from directory
		currentFiles, err := db.GetIdsArray(directoryId, "files")
		if err != nil {
			return -1, err
		}
		// Checking containing file in directory
		if h.Contains(currentFiles, file.Id) {
			return -1, fmt.Errorf("file already exists in folder")
		}
	}else{
		// Inserting file into database
		idStr, err := db.insert("insert into files (Name, FileId, FileUniqueId, FileType, Created, FileSize, ThumbnailFileId) values ($1, $2, $3, $4, now(), $5, $6) returning Id", 
				file.Name, file.FileId, file.FileUniqueId, file.FileType, file.FileSize, file.ThumbnailFileId)
		
		if err != nil {
			return -1, err
		}
		// Converting id to integer
		file.Id, err = strconv.Atoi(idStr)
		if err != nil {
			return -1, err
		}
	}
	// Returning file id and result of adding file to directory
	return file.Id, db.AddFileToDirectory(directoryId, file)
}
// Function for adding file to directory
func (db *DataBase) AddFileToDirectory(directoryId int, file File) error {
	req := fmt.Sprintf("update directories set files = array_append(files, %d), size = size + %d where id = %d", 
		file.Id, file.FileSize, directoryId)
	return db.makeQuery(req)
}
// Function for checking existence of directory by name
func (db *DataBase) FolderExists(userId int64, currentDirectoryId int, directoryName string) (bool, error) {
	// Getting an array of directory names contained in the current directory
	currentDirectories, err := db.GetNamesArray(currentDirectoryId, "directories")
	return h.Contains(currentDirectories, directoryName), err
}
// Function for adding a directory into the current directory
func (db *DataBase) AddNewDirectoryToDirectory(currentDirectoryId int, newDirectoryId int) error {
	req := fmt.Sprintf("update directories set directories = array_append(directories, %d) where id = %d", 
		newDirectoryId, currentDirectoryId)
	return db.makeQuery(req)
}
// Function for updating current directory
func (db *DataBase) JumpToDirectory(userId int64, directoryId int) error {
	req := fmt.Sprintf("update users set currentDirectory = %d where userId = %d", directoryId, userId)
	return db.makeQuery(req)
}
// Getting parent for a directory
func (db *DataBase) GetParentDirectory(directoryId int) (Directory, error) {
	// Getting directory information
	currentDirectory, err := db.GetDirectory(directoryId)
	if err != nil{
		return Directory{}, err
	}
	// If the directory is not root then error is returned
	if currentDirectory.ParentId == -1 {
		return Directory{}, fmt.Errorf("this directory is root")
	}
	// Else information about the parent is returned
	return db.GetDirectory(currentDirectory.ParentId)
}
// Function for creating a new user
func (db *DataBase) CreateNewUser(user User) (string, error) {
	// Checking existence of this userId
	exist, err := db.UserExists(int64(user.UserId))
	if err != nil {
		return "", err
	}
	// If user exists then error is returned
	if exist {
		return "", errors.New("user already exists")
	}
	// Creating a root directory for user
	directoryId, err := db.CreateRootDirectory(Directory{Name : "(root)/", UserId: user.UserId})
	if err != nil {
		return "", err
	}
	// Inserting a new user into the database
	return db.insert("insert into users (Username, UserID, FirstName, LastName, CurrentDirectory) values ($1, $2, $3, $4, $5) returning Id", 
		user.UserName, user.UserId, user.FirstName, user.LastName, directoryId)
}
// Function for checking existence of user by userId
func (db *DataBase) UserExists(userId int64) (bool, error) {
	req := fmt.Sprintf("select id from users where userId = %d", userId)
	id, err := db.selectRow(req)
	return id != "", err
}
// Function for getting id by uniqueId
func (db *DataBase) GetIdOfFileByUniqueId(fileUniqueId string) (int, error) {
	req := fmt.Sprintf("select id from files where fileUniqueId = '%s'", fileUniqueId)
	id, err := db.selectRow(req)
	if err != nil {
		return -1, err
	}
	// If id not found, returning -1
	if id == "" {
		return -1, nil
	}
	// Returning result of conversion id to integer
	return strconv.Atoi(id)
}
// Function for getting current directory
func (db *DataBase) GetCurrentDirectory(userId int64) (int, error) {
	req := fmt.Sprintf("select currentDirectory from users where userId = %d", userId)
	directoryId, err := db.selectRow(req)
	if err != nil {
		return -1, err
	}
	// If id not found, returning -1
	if directoryId == "" {
		return -1, errors.New("directory or user not found!")
	}
	// Returning result of conversion id to integer
	return strconv.Atoi(directoryId)
}
// Function for getting information about directory by id
func (db *DataBase) GetDirectory(id int) (Directory, error) {
	d := Directory{}
	req := fmt.Sprintf("select id, parentId, name, userId, size, path, created from directories where id = %d", id)
	rows, err := db.selectRows(req)
	if err != nil {
		return d, err
	}
	rows.Next()
	// Setting data into directory instance
	if err := rows.Scan(&d.Id, &d.ParentId, &d.Name, &d.UserId, &d.Size, &d.Path, &d.Created); err != nil {
		return d, err
	}
	// Setting files and directories id`s
	d.Files, err = db.GetIdsArray(id, "files")
	if err != nil {
		return d, err
	}
	d.Directories, err = db.GetIdsArray(id, "directories")
	if err != nil {
		return d, err
	}
	return d, err
}
// Function for getting the file information by id
func (db *DataBase) GetFileById(id int) (File, error) {
	f := File{}
	req := fmt.Sprintf("select * from files where id = %d", id)
	rows, err := db.selectRows(req)
	if err != nil {
		return f, err
	}
	rows.Next()
	// Setting data into directory instance
	err = rows.Scan(&f.Id, &f.Name, &f.FileId, &f.FileUniqueId, &f.FileSize, &f.FileType, &f.Created, &f.ThumbnailFileId, &f.ThumbnailSource, &f.FileSource)
	return f, err
}
// Function for getting available directories in directory
func (db *DataBase) GetAvailableDirectoriesInDiretory(userId int64, directoryId int) ([]Directory, error) {
	var res []Directory
	// Adding parrent directory
	d, err := db.GetParentDirectory(directoryId)
	if err == nil {
		d.Name = "../"
		res = append(res, d)
	}
	// Getting an array of directory ids contained in the directory
	arr, err := db.GetIdsArray(directoryId, "directories")
	if err != nil {
		return nil, err
	}
	// Iterating through directory ids and getting the directories inforamation
	for _, id := range arr {
		d, err := db.GetDirectory(id)
		if err != nil {
			return res, err
		}
		res = append(res, d)
	}
	// Returning result array
	return res, nil
}
// Function for getting available files in directory
func (db *DataBase) GetAvailableFilesInDiretory(userId int64, directoryId int) ([]File, error) {
	var res []File
	
	// Getting an array of files ids contained in the directory
	arr, err := db.GetIdsArray(directoryId, "files")
	if err != nil {
		return nil, err
	}
	// Iterating through directory ids and getting the files inforamation
	for _, id := range arr {
		d, err := db.GetFileById(id)
		if err != nil {
			return res, err
		}
		res = append(res, d)
	}
	// Returning result array
	return res, nil
}
// Function for getting available items in directory
func (db *DataBase) GetAvailableItemsInDirectory (userId int64, directoryId int) (DirectoryContent, error) {
	// Getting available directories
	directories, err := db.GetAvailableDirectoriesInDiretory(userId, directoryId)
	if err != nil {
		return DirectoryContent{}, err
	}
	// Getting available files
	files, err := db.GetAvailableFilesInDiretory(userId, directoryId)
	// Returning result
	return DirectoryContent{
		Directories: directories,
		Files: files,
	}, err
}
// Function for getting array of id`s ("files"/"directories" - parameter 'name') from directory
func (db *DataBase) GetIdsArray(directoryId int, name string) ([]int, error) {
	req := fmt.Sprintf("select %s from directories where id = %d", name, directoryId)
	idsStr, err := db.selectRow(req)
	if err != nil {
		return nil, err
	}
	// Parsing JSON and returning an array of integers
	return h.ParseIds(idsStr)
}
// Function for getting array of id`s ("files"/"directories" - parameter 'name') from directory
func (db *DataBase) GetNamesArray(directoryId int, name string) ([]string, error) {
	req := fmt.Sprintf("select %s from directories where id = %d", name, directoryId)
	idsStr, err := db.selectRow(req)
	if err != nil {
		return nil, err
	}
	// Parsing JSON
	ids, err := h.ParseIds(idsStr)
	if err != nil {
		return nil, err
	}
	// Checking length of array
	if len(ids) == 0 {
		return nil, nil
	}
	// Getting name for each id
	req = fmt.Sprintf("select name from directories where id in ( %s )", strings.Join(h.IntArrayToStrArray(ids), ", "))
	rows, err := db.selectRows(req)
	if err != nil {
		return nil, err
	}
	// Iterating through the rows and compilation result array
	var res []string
	for rows.Next() {
		current := ""
		err = rows.Scan(&current)
		if err != nil {
			return nil, err
		}
		res = append(res, current)
	}
	// Returning an array of names
	return res, nil
}
// Reseting the root directory from the database for user
func (db *DataBase) ResetUserData(userId int64) error {
	// Creating new root directory for user
	directoryId, err := db.CreateRootDirectory(Directory{Name : "(root)/", UserId: int(userId)})
	if err != nil {
		return err
	}
	// Updating the root directory for user
	req := fmt.Sprintf("update users set currentDirectory = %s where userId=%d", directoryId, userId)
	if err := db.makeQuery(req); err != nil {
		return err
	}
	// Deleting all user`s directories except the new root directory
	req = fmt.Sprintf("delete from directories where userId=%d and id != %s", userId, directoryId)
	return db.makeQuery(req)
}
// Updating source to file in database
func (db *DataBase) UpdateSource(fileId int, newSource string, isThumbnail bool) error {
	field := "fileSource"
	if isThumbnail {
		field = "thumbnailSource"
	}
	
	req := fmt.Sprintf("update files set %s = '%s' where id=%d", field, newSource, fileId)
	log.Print(req)
	return db.makeQuery(req)
}
// Getting user data from database
func (db *DataBase) GetUserInfo(userId int64) (user User, err error) {
	u := User{}
	req := fmt.Sprintf("select * from users where UserID = %d", userId)
	rows, err := db.selectRows(req)
	if err != nil {
		return u, err
	}
	rows.Next()
	// Setting data into user instance
	err = rows.Scan(&u.Id, &u.UserName, &u.UserId, &u.FirstName, &u.LastName, &u.CurrentDirectory)
	return u, err
}

// Function for updating item name
func (db *DataBase) UpdateItemName(id int, newName string, typeItem string) error {
	if typeItem != "directory" {
		// make request if its file updating
		req := fmt.Sprintf("update files set name = '%s' where id=%d", newName, id)
		return db.makeQuery(req)
	}
	// else getting directory info
	directory, err := db.GetDirectory(id)
	if err != nil {
		return err
	}
	// check if it is not root directory
	if directory.ParentId == -1 {
		return fmt.Errorf("you cannot rename a root directory")
	}
	// make directory updating request
	req := fmt.Sprintf("update directories set name = '%s' where id=%d", newName, id)
	if err := db.makeQuery(req); err != nil {
		return err
	}
	// create new path and update it for all child elements
	p := strings.Split(directory.Path, "/")
	return db.UpdatePath(id, strings.Join(p[:len(p) - 2], "/") + "/")
	
}

func (db *DataBase) UpdatePath(id int, path string) error {
	req := fmt.Sprintf("update directories set path = CONCAT('%s', name, '/') where id=%d returning path, directories", path, id)
	rows, err := db.selectRows(req)
	if err != nil {
		return err
	}
	rows.Next()
	var newPath, directoriesStr string;
	if err := rows.Scan(&newPath, &directoriesStr); err != nil {
		return err
	}
	fmt.Print(newPath, directoriesStr)
	directories, err := h.ParseIds(directoriesStr);
	if err != nil {
		return err
	}
	for _, childId := range directories {
		err := db.UpdatePath(childId, newPath)
		if err != nil {
			return err
		}
	}
	return nil
}