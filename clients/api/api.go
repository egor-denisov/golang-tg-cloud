package api

import (
	"log"
	"main/db"
	"main/storage"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ApiClient struct {
	router *gin.Engine
	db db.DataBase
	storage storage.Storage
}

func New(db db.DataBase, s storage.Storage) ApiClient {
	res := ApiClient{
		router:  gin.Default(),
		db:      db,
		storage: s,
	}
	res.router.GET("/directory", res.getDirecoryById)
	res.router.GET("/file", res.getFileById)
	return res
}

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
	
	file, err := api.db.GetFileById(id)
	if err != nil {
		ProccessError(context, err)
		return
	}

	//context.IndentedJSON(http.StatusOK, file)

	url, err := api.storage.GetFile(file.FileId)
	if err != nil {
		ProccessError(context, err)
		return 
	}
	context.File(url)
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