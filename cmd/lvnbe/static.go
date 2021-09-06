package main

import (
	"embed"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

//go:embed build/*
var build embed.FS

func NewEmbedHandler(log Logger) Handler {
	mime.AddExtensionType(".map", "application/json")
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		if path == "/" {
			path = "/index.html"
		}
		ext := filepath.Ext(path)
		ct := mime.TypeByExtension(ext)
		log.Trace("static", path, ct)
		file, err := build.Open("build" + path)
		if err != nil {
			log.Debug(path, err)
			ctx.Response.SetStatusCode(http.StatusNotFound)
			return
		}
		defer file.Close()
		ctx.SetContentType(ct)
		io.Copy(ctx, file)
	}
}
