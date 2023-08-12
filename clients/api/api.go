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
	router := gin.Default()
	router.MaxMultipartMemory = 8 * 20
	
	res := ApiClient{
		router:  router,
		db:      db,
		storage: s,
	}
	res.router.GET("/directory", res.getDirecoryById)
	res.router.GET("/file", res.getFileById)
	res.router.POST("/upload", res.uploadFile)
	return res
}

// Function for listing server
func (api *ApiClient) Listen() {
	api.router.Run()
}

func (api *ApiClient) getDirecoryById(context *gin.Context) {
	
	id, err := strconv.Atoi(context.Query("id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	
	d, err := api.db.GetDirectory(id)
	if err != nil {
		ProccessError(context, err)
		return
	}

	context.IndentedJSON(http.StatusOK, d)
}

func (api *ApiClient) getFileById(context *gin.Context) {
	
	id, err := strconv.Atoi(context.Query("id"))
	if err != nil {
		ProccessError(context, err)
		return
	}
	
	fileData, err := api.db.GetFileById(id)
	if err != nil {
		ProccessError(context, err)
		return
	}

	fileBytes, err := api.storage.GetFileAsBytes(fileData.FileId)
	if err != nil {
		ProccessError(context, err)
		return 
	}
	
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileData.Name))
	http.ServeContent(context.Writer, context.Request, "filename", time.Now(), bytes.NewReader(fileBytes))
}

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
	userData := context.PostForm("user-data")
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		ProccessError(context, err)
		return
	}
	// Adding our data in queue for later adding to telegram server
	api.storage.AddToUploadingQueue(path, headers.Filename, user)
}

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