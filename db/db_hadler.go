package db

import (
	"errors"
	"fmt"
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
	existence, err := db.FolderExists(directory.UserId, directory.ParentId, directory.Name)
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
	file.UserId = int(userId)
	// Checking containing file in directory
	exists, err := db.FileExists(file.UserId, directoryId, file.Name)
	if err != nil {
		return -1, err
	}
	if exists {
		return -1, fmt.Errorf("file already exists in folder")
	}
	// Getting id by unique id
	id, err := db.GetIdOfFileByUniqueId(file.FileUniqueId)
	if err != nil {
		return -1, err
	}
	// If file already exists then error is returned 
	if id > 0 {
		// Getting source from similar file
		fileInfo, err := db.GetFileById(userId, id, true);
		if err != nil {
			return -1, err
		}
		file.FileSource = fileInfo.FileSource
		file.ThumbnailSource = fileInfo.ThumbnailSource
	}
	// Inserting file into database
	idStr, err := db.insert("insert into files (userId, Name, FileId, FileUniqueId, FileType, Created, FileSize, ThumbnailFileId, FileSource, ThumbnailSource) values ($1, $2, $3, $4, $5, now(), $6, $7, $8, $9) returning Id", 
		file.UserId, file.Name, file.FileId, file.FileUniqueId, file.FileType, file.FileSize, file.ThumbnailFileId, file.FileSource, file.ThumbnailSource)
	if err != nil {
		return -1, err
	}
	// Converting id to integer
	file.Id, err = strconv.Atoi(idStr)
	if err != nil {
		return -1, err
	}
	// Returning file id and result of adding file to directory
	return file.Id, db.AddFileToDirectory(directoryId, file)
}
// Function for adding file to directory
func (db *DataBase) AddFileToDirectory(directoryId int, file File) error {
	if err := db.addSizeFileToDirectory(directoryId, file.FileSize); err != nil {
		return err
	}
	req := fmt.Sprintf(`update directories set files = array_append(files, %d) where id = %d;`, file.Id, directoryId)
	return db.makeQuery(req)
}
// Function for changing size of directory and its parents
func (db *DataBase) addSizeFileToDirectory(directoryId int, fileSize int) error {
	req := fmt.Sprintf(`
	WITH RECURSIVE r AS (SELECT Id, ParentId FROM directories WHERE id = %d UNION
		SELECT directories.Id, directories.ParentId FROM directories JOIN r ON directories.Id = r.ParentId)
   
		update directories set size = size + (%d) where Id IN (SELECT Id FROM r);
		`, directoryId, fileSize)
	return db.makeQuery(req)
}
// Function for checking existence of directory by name
func (db *DataBase) FolderExists(userId int, currentDirectoryId int, directoryName string) (bool, error) {
	// Getting an array of directory names contained in the current directory
	currentDirectories, err := db.GetNamesArray(userId, currentDirectoryId, "directories")
	return h.Contains(currentDirectories, directoryName), err
}
// Function for checking existence of directory by name
func (db *DataBase) FileExists(userId int, currentDirectoryId int, fileName string) (bool, error) {
	// Getting an array of file names contained in the current directory
	currentFiles, err := db.GetNamesArray(userId, currentDirectoryId, "files")
	return h.Contains(currentFiles, fileName), err
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
func (db *DataBase) CreateNewUser(user User) (User, error) {
	var u User
	// Checking existence of this userId
	exist, err := db.UserExists(int64(user.UserId), user.UserName)
	if err != nil {
		return u, err
	}
	// If user exists then error is returned
	if exist {
		return u, errors.New("user already exists")
	}
	// Creating a root directory for user
	currentDirectoryStr, err := db.CreateRootDirectory(Directory{Name : "(root)/", UserId: user.UserId})
	if err != nil {
		return u, err
	}
	user.CurrentDirectory, err = strconv.Atoi(currentDirectoryStr)
	if err != nil {
		return u, err
	}
	// Creating hash for user
	user.Hash, err = h.HashData(strconv.Itoa(user.UserId), user.UserName, user.FirstName, user.LastName)
	if err != nil {
		return u, err
	}
	// Inserting a new user into the database
	idStr, err :=  db.insert("insert into users (Username, UserID, FirstName, LastName, CurrentDirectory, Hash) values ($1, $2, $3, $4, $5, $6) returning id", 
		user.UserName, user.UserId, user.FirstName, user.LastName, user.CurrentDirectory, user.Hash)
	if err != nil {
		return u, err
	}
	user.Id, err = strconv.Atoi(idStr)

	return user, err
}
// Function for checking existence of user by userId
func (db *DataBase) UserExists(userId int64, username string) (bool, error) {
	req := fmt.Sprintf("select id from users where userId = %d or username = '%s'", userId, username)
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
	defer rows.Close()
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
func (db *DataBase) GetFileById(userId int64, id int, unsafe bool) (File, error) {
	f := File{}
	req := ""
	if unsafe {
		req = fmt.Sprintf("select * from files where id = %d", id)
	}else{
		req = fmt.Sprintf("select * from files where id = %d and userId = %d", id, userId)
	}
	rows, err := db.selectRows(req)
	defer rows.Close()
	if err != nil {
		return f, err
	}
	rows.Next()
	// Setting data into directory instance
	err = rows.Scan(&f.Id, &f.UserId, &f.Name, &f.FileId, &f.FileUniqueId, &f.FileSize, &f.FileType, &f.Created, &f.ThumbnailFileId, &f.ThumbnailSource, &f.FileSource, &f.SharedId, &f.IsShared)
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
		d, err := db.GetFileById(userId, id, false)
		if err != nil {
			return res, err
		}
		res = append(res, d)
	}
	// Returning result array
	return res, nil
}
// Function for getting available items in directory
func (db *DataBase) GetAvailableItemsInDirectory(userId int64, directoryId int) (DirectoryContent, error) {
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
func (db *DataBase) GetNamesArray(userId int, directoryId int, name string) ([]string, error) {
	req := fmt.Sprintf("select %s from directories where id = %d and userId = %d", name, directoryId, userId)
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
	req = fmt.Sprintf("select name from %s where id in ( %s )", name, strings.Join(h.IntArrayToStrArray(ids), ", "))
	rows, err := db.selectRows(req)
	defer rows.Close()
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
func (db *DataBase) ResetUserData(userId int64) (int, error) {
	// Creating new root directory for user
	directoryId, err := db.CreateRootDirectory(Directory{Name : "(root)/", UserId: int(userId)})
	if err != nil {
		return -1, err
	}
	// Updating the root directory for user
	req := fmt.Sprintf("update users set currentDirectory = %s where userId=%d", directoryId, userId)
	if err := db.makeQuery(req); err != nil {
		return -1, err
	}
	// Deleting all user`s directories except the new root directory
	req = fmt.Sprintf("delete from directories where userId=%d and id != %s", userId, directoryId)
	if err := db.makeQuery(req); err != nil {
		return -1, err
	}
	// Returning new root directory id
	return strconv.Atoi(directoryId)
}
// Updating source to file in database
func (db *DataBase) UpdateSource(fileId int, newSource string, isThumbnail bool) error {
	field := "fileSource"
	if isThumbnail {
		field = "thumbnailSource"
	}
	
	req := fmt.Sprintf("update files set %s = '%s' where id=%d", field, newSource, fileId)
	return db.makeQuery(req)
}
// Getting user data from database
func (db *DataBase) GetUserInfo(userId int64) (user User, err error) {
	u := User{}
	req := fmt.Sprintf("select * from users where UserID = %d", userId)
	rows, err := db.selectRows(req)
	defer rows.Close()
	if err != nil {
		return u, err
	}
	rows.Next()
	// Setting data into user instance
	err = rows.Scan(&u.Id, &u.UserName, &u.UserId, &u.FirstName, &u.LastName, &u.CurrentDirectory, &u.Hash)
	return u, err
}

// Getting user data from database
func (db *DataBase) GetUserHash(userId int64) (hash string, err error) {
	req := fmt.Sprintf("select hash from users where UserID = %d", userId)
	return db.selectRow(req)
}

// Function for updating item name
func (db *DataBase) UpdateItemName(userId int, id int, directoryId int, newName string, typeItem string) error {
	if typeItem != "directory" {
		// Checking existence of folder with current name
		existence, err := db.FileExists(userId, directoryId, newName)
		if err != nil {
			return err
		}
		// Return error if directory exists
		if existence {
			return fmt.Errorf("file with this name already exists in current folder")
		}
		// make request if its file updating
		req := fmt.Sprintf("update files set name = '%s' where id=%d and userId=%d", newName, id, userId)
		return db.makeQuery(req)
	}
	// Checking existence of folder with current name
	existence, err := db.FolderExists(userId, directoryId, newName)
	if err != nil {
		return err
	}
	// Return error if directory exists
	if existence {
		return fmt.Errorf("directory with this name already exists in current folder")
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
	req := fmt.Sprintf("update directories set name = '%s' where id=%d and userId=%d", newName, id, userId)
	if err := db.makeQuery(req); err != nil {
		return err
	}
	// create new path and update it for all child elements
	p := strings.Split(directory.Path, "/")
	return db.UpdatePath(id, strings.Join(p[:len(p) - 2], "/") + "/")
	
}

// Function for erasing item
func (db *DataBase) DeleteItem(id int, userId int, directoryId int, typeItem string) error {
	if typeItem != "directory" {
		// make request if its filev
		fileInfo, err := db.GetFileById(int64(userId), id, false)
		if err != nil {
			return err
		}
		if err := db.addSizeFileToDirectory(directoryId, -fileInfo.FileSize); err != nil {
			return err
		}
		req := fmt.Sprintf("update directories set files = array_remove(files, %d) where id=%d and userId=%d", id, directoryId, userId)
		return db.makeQuery(req)
	}
	directoryInfo, err := db.GetDirectory(id)
	if err != nil {
		return err
	}
	if err := db.addSizeFileToDirectory(directoryId, -directoryInfo.Size); err != nil {
		return err
	}
	req := fmt.Sprintf(`
		update directories set directories = array_remove(directories, %d) where id=%d and userId=%d;
		delete from directories where id=%d
	`, id, directoryId, userId, id)
	return db.makeQuery(req)
	
}

// Function to change sharing item
func (db *DataBase) ChangeSharingFile(id int, userId int, share bool) (string, error) {
	file, err := db.GetFileById(int64(userId), id, false)
	if err != nil {
		return "", err
	}
	req := fmt.Sprintf("update files set isShared = %t where id=%d and userId=%d returning SharedId", share, file.Id, userId)
	return db.selectRow(req)
}

func (db *DataBase) UpdatePath(id int, path string) error {
	req := fmt.Sprintf("update directories set path = CONCAT('%s', name, '/') where id=%d returning path, directories", path, id)
	rows, err := db.selectRows(req)
	defer rows.Close()
	if err != nil {
		return err
	}
	rows.Next()
	var newPath, directoriesStr string;
	if err := rows.Scan(&newPath, &directoriesStr); err != nil {
		return err
	}
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

func (db *DataBase) GetSharedItemId(sharedId string) (int, error) {
	req := fmt.Sprintf("select id from files where isShared = true and sharedId = '%s'", sharedId)
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

// Function for checking hash
func (db *DataBase) CheckHash(userId int, hash string) (error) {
	userHash, err := db.GetUserHash(int64(userId))
	if err != nil {
		return err
	}
	if userHash != hash {
		return fmt.Errorf("hash mismatch")
	}
	return nil
}