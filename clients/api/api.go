package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"main/db"
	"main/lib/h"
	"main/storage"
	. "main/types"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiClient struct {
	router *gin.Engine
	db db.DataBase
	storage storage.Storage
}

// Function for creating a new ApiClient instance
func New(db db.DataBase, s storage.Storage) ApiClient {
	// Creating a new router and set max limit of memory for it
	router := gin.Default()
	router.MaxMultipartMemory = 8 * 20
	// Creating a new instance of the ApiClient
	res := ApiClient{
		router:  router,
		db:      db,
		storage: s,
	}
	// Setting route for the router
	res.router.GET("/directory", res.getDirecoryById)
	res.router.GET("/file", res.getFileById)
	res.router.GET("/thumbnail", res.getThumbnailById)
	res.router.POST("/upload", res.uploadFile)
	res.router.GET("/available", res.getAvailableItems)
	res.router.GET("/auth", res.authorization)
	res.router.GET("/createDirectory", res.createDirectory)
	res.router.GET("/edit", res.editItem)

	res.router.OPTIONS("/upload", res.preloader)
	return res
}

// Function for listing server
func (api *ApiClient) Listen() {
	api.router.Run()
}
// Function returns a directory info by id
func (api *ApiClient) getDirecoryById(context *gin.Context) {
	// Getting the directory id from request parameters and convert it to a number
	id, err := strconv.Atoi(context.Query("id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Getting data about the directory by id
	d, err := api.db.GetDirectory(id)
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Return the directory data object
	setHeaders(context)
	context.IndentedJSON(http.StatusOK, d)
}

// Function returns a file by id
func (api *ApiClient) getFileById(context *gin.Context) {
	// Getting the directory id from request parameters and convert it to a number
	id, err := strconv.Atoi(context.Query("id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Getting data about the file by id
	fileData, err := api.db.GetFileById(id)
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Checking hashing value in database
	if len(fileData.FileSource) <= 0 {
		fileData.FileSource, err = api.storage.GetFileURL(fileData.FileId)
		if err != nil {
			ProccessError(context, err)
			return
		}
		if err := api.db.UpdateSource(fileData.Id, fileData.FileSource, false); err != nil {
			ProccessError(context, err)
			return
		}
	}
	// Getting file as bytes from storage
	fileBytes, err := api.storage.GetFileAsBytes(fileData.FileSource)
	if err != nil {
		ProccessError(context, err)
		return 
	}
	// Setting headers and provide file for downloading
	setHeaders(context)
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileData.Name))
	http.ServeContent(context.Writer, context.Request, "filename", time.Now(), bytes.NewReader(fileBytes))
}

// Function returns a thumbnail by id
func (api *ApiClient) getThumbnailById(context *gin.Context) {
	// Getting the directory id from request parameters and convert it to a number
	id, err := strconv.Atoi(context.Query("id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Getting data about the file by id
	fileData, err := api.db.GetFileById(id)
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Checking hashing value in database
	if len(fileData.ThumbnailSource) <= 0 {
		fileData.ThumbnailSource, err = api.storage.GetFileURL(fileData.ThumbnailFileId)
		if err != nil {
			ProccessError(context, err)
			return
		}
		if err := api.db.UpdateSource(fileData.Id, fileData.ThumbnailSource, true); err != nil {
			ProccessError(context, err)
			return
		}
	}
	// Getting file as bytes from storage
	fileBytes, err := api.storage.GetFileAsBytes(fileData.ThumbnailSource)
	if err != nil {
		ProccessError(context, err)
		return 
	}
	// Setting headers and provide file for downloading
	setHeaders(context)
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileData.Name))
	http.ServeContent(context.Writer, context.Request, "filename", time.Now(), bytes.NewReader(fileBytes))
}
// Function for sending headers for uploading files
func (api *ApiClient) preloader(context *gin.Context) {
	setHeaders(context)
}
// Function for uploading new files 
func (api *ApiClient) uploadFile(context *gin.Context) {
	setHeaders(context)
	// Getting uploaded file
	fmt.Printf("%+v\n", context.Request.Header)
	_, headers, err := context.Request.FormFile("file")
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Saving file in temp storage
	path := "./assets/" + h.GenerateUniqueName()
	if err := context.SaveUploadedFile(headers, path); err != nil {
		ProccessError(context, err)
		return
	}
	// Getting user id
	userId, err := strconv.Atoi(context.PostForm("user_id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Getting directory id
	directoryId, err := strconv.Atoi(context.PostForm("directory_id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	
	// Adding our data in queue for later adding to telegram server
	api.storage.AddToUploadingQueue(path, headers.Filename, userId, directoryId)
}

// Function for getting available files and directories
func (api *ApiClient) getAvailableItems(context *gin.Context) {
	// Getting the user id from request parameters and convert it to a number
	userId, err := strconv.Atoi(context.Query("user_id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Getting the directory id from request parameters and convert it to a number
	directoryId, err := strconv.Atoi(context.Query("directory_id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Getting data about the directory by id
	items, err := api.db.GetAvailableItemsInDirectory(int64(userId), directoryId)
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Return the directory data object
	setHeaders(context)
	context.IndentedJSON(http.StatusOK, items)
}

// Function for authorization user
func (api *ApiClient) authorization(context *gin.Context) {
	// Getting user_id from data
	userId, err := strconv.Atoi(context.Query("user_id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	userInfo, err := api.db.GetUserInfo(int64(userId))
	if err != nil {
		ProccessError(context, err)
		return
	}
	setHeaders(context)
	context.IndentedJSON(http.StatusOK, userInfo)
}

// Function for editing the file or directory
func (api *ApiClient) editItem(context *gin.Context) {
	// Getting user_id from data
	id, err := strconv.Atoi(context.Query("id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	newName := context.Query("name")
	typeItem := context.Query("type")

	setHeaders(context)
	// Setting new name
	if err := api.db.UpdateItemName(id, newName, typeItem); err != nil {
		ProccessError(context, err)
		return
	}
}

// Function for creating directory
func (api *ApiClient) createDirectory(context *gin.Context) {
	// Getting the directory id from request parameters and convert it to a number
	var directoryInfo Directory
	if err := json.Unmarshal([]byte(context.Query("directory")), &directoryInfo); err != nil {
		ProccessError(context, err)
		return
	}
	// Creating new folder in the database
	if _, err := api.db.CreateNewDirectory(directoryInfo); err != nil {
		ProccessError(context, err)
		return
	}
	setHeaders(context)
	context.IndentedJSON(http.StatusOK, "ok")
}

// Function for proccessing errors in working of api
func ProccessError(context *gin.Context, err error) {
	log.Print(err)
	setHeaders(context)
	switch strings.Split(err.Error(), ":")[0] {
	case "sql":
		context.IndentedJSON(http.StatusNotFound, "Not found")
	case "strconv.Atoi":
		context.IndentedJSON(http.StatusNotAcceptable, "Wrong format of request")
	case "directory with this name already exists in current folder":
		context.IndentedJSON(http.StatusNotAcceptable, "Directory with this name already exists in current folder")
	case "wrong folder name":
		context.IndentedJSON(http.StatusNotAcceptable, "Directory name is uncorrectly. Don`t use symbols: <, >, :, «, /,\\ , |, ?, *")
	default: 
		context.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
	
}
// Function for creating headers
func setHeaders(context *gin.Context) {
	context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	context.Writer.Header().Set("Access-Control-Allow-Credential", "true")
	context.Writer.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	context.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

}