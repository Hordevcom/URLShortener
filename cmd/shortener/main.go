package main

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var urlStore = make(map[string]string)

func shortenUrl(ctx *gin.Context) {
	url := ctx.PostForm("url")

	if url == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "url param required"})
		return
	}

	shortUrl := fmt.Sprintf("%x", md5.Sum([]byte(url)))[:8]
	urlStore[shortUrl] = url
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "http://localhost:8080/" + shortUrl})
}

func redirect(ctx *gin.Context) {

	shortUrl := ctx.Param("id")

	if urlStore[shortUrl] != "" {
		ctx.IndentedJSON(http.StatusTemporaryRedirect, gin.H{"message": "Location: " + urlStore[shortUrl]})
	} else {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{})
	}
}

func main() {

	server := gin.Default()
	server.HandleMethodNotAllowed = true
	server.POST(`/`, shortenUrl)
	server.GET(`/:id`, redirect)

	server.Run("localhost:8080")
}
