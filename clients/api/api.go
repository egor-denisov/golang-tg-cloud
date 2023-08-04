package api

import (
	"main/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ApiClient struct {
	router *gin.Engine
	db db.DataBase
}

func New(db db.DataBase) ApiClient {
	res := ApiClient{
		router:  gin.Default(),
		db:      db,
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

	context.IndentedJSON(http.StatusOK, file)
}

func ProccessError(context *gin.Context, err error) {
	switch strings.Split(err.Error(), ":")[0] {
	case "sql":
		context.IndentedJSON(http.StatusNotFound, "Not found")
	case "strconv.Atoi":
		context.IndentedJSON(http.StatusNotAcceptable, "Wrong format of request")
	default: 
		context.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
	
}