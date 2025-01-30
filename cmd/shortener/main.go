package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var urlStore = make(map[string]string)

func shortenURL(ctx *gin.Context) {
	url, _ := io.ReadAll(ctx.Request.Body)

	if string(url) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "url param required"})
		return
	}

	shortURL := fmt.Sprintf("%x", md5.Sum([]byte(url)))[:8]
	urlStore[shortURL] = string(url)
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "http://localhost:8080/" + shortURL})
}

func redirect(ctx *gin.Context) {

	shortURL := ctx.Param("id")

	if urlStore[shortURL] != "" {
		ctx.Header("Location", urlStore[shortURL])
		ctx.IndentedJSON(http.StatusTemporaryRedirect, gin.H{})
	} else {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{})
	}
}

func main() {

	server := gin.Default()
	server.HandleMethodNotAllowed = true
	server.POST(`/`, shortenURL)
	server.GET(`/:id`, redirect)

	server.Run("localhost:8080")
}
