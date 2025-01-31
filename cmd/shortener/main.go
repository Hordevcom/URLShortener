package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	"github.com/Hordevcom/URLShortener/cmd/shortener/config"
	"github.com/gin-gonic/gin"
)

var urlStore = make(map[string]string)
var conf config.Config

func shortenURL(ctx *gin.Context) {
	url, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	if string(url) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "url param required"})
		return
	}

	shortURL := fmt.Sprintf("%x", md5.Sum([]byte(url)))[:8]
	urlStore[shortURL] = string(url)
	ctx.String(http.StatusCreated, "http://"+conf.Host+"/"+shortURL)
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
	conf = config.GetCLParams()

	server := gin.Default()
	server.HandleMethodNotAllowed = true
	server.POST(`/`, shortenURL)
	server.GET(`/:id`, redirect)

	server.Run(conf.ServerAdress)
}
