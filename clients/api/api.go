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
	res.router.POST("/upload", res.uploadFile)
	res.router.GET("/available", res.getAvailableItems)

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
	// Getting file as bytes from storage
	fileBytes, err := api.storage.GetFileAsBytes(fileData.FileId)
	if err != nil {
		ProccessError(context, err)
		return 
	}
	// Setting headers and provide file for downloading
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileData.Name))
	http.ServeContent(context.Writer, context.Request, "filename", time.Now(), bytes.NewReader(fileBytes))
}

// Function for uploading new files 
func (api *ApiClient) uploadFile(context *gin.Context) {
	// Getting uploaded file
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
	// Getting user data
	var user User
	userData := context.PostForm("user_data")
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		ProccessError(context, err)
		return
	}
	// Getting directory id
	directoryId, err := strconv.Atoi(context.PostForm("directory"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	// Adding our data in queue for later adding to telegram server
	api.storage.AddToUploadingQueue(path, headers.Filename, user, directoryId)
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
	context.IndentedJSON(http.StatusOK, items)
}
// Function for proccessing errors in working of api
func ProccessError(context *gin.Context, err error) {
	log.Print(err)
	switch strings.Split(err.Error(), ":")[0] {
	case "sql":
		context.IndentedJSON(http.StatusNotFound, "Not found")
	case "strconv.Atoi":
		context.IndentedJSON(http.StatusNotAcceptable, "Wrong format of request")
	default: 
		context.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
	
}