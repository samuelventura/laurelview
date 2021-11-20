package main

import (
	"embed"
	"log"
	"mime"
	"net"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/samuelventura/go-tree"
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

func api(node tree.Node) {
	mime.AddExtensionType(".map", "application/json")
	endpoint := node.GetValue("endpoint").(string)
	gin.SetMode(gin.ReleaseMode) //remove debug warning
	router := gin.New()          //remove default logger
	router.Use(gin.Recovery())   //looks important
	router.Use(static)
	rapi := router.Group("/api")
	rapi.GET("/ok", func(c *gin.Context) {
		c.JSON(200, "ok")
	})
	listen, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
	}
	node.AddCloser("listen", listen.Close)
	port := listen.Addr().(*net.TCPAddr).Port
	log.Println("port", port)
	server := &http.Server{
		Addr:    endpoint,
		Handler: router,
	}
	node.AddProcess("server", func() {
		err = server.Serve(listen)
		if err != nil {
			log.Println(endpoint, port, err)
		}
	})
}
