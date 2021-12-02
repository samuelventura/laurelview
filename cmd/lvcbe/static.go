package main

import (
	"embed"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//go:embed build/*
var build embed.FS

func static(c *gin.Context) {
	path := "build" + c.Request.URL.Path
	if path == "build/" {
		path = "build/index.html"
	}
	//log.Println(path)
	data, err := build.ReadFile(path)
	if err != nil {
		c.Next()
	} else {
		ext := filepath.Ext(path)
		ct := mime.TypeByExtension(ext)
		c.Data(http.StatusOK, ct, data)
	}
}
