package main

import (
	"embed"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/samuelventura/laurelview/pkg/lvsdk"
	"github.com/valyala/fasthttp"
)

//go:embed build/*
var fs embed.FS

var cache map[string]*staticDso

type staticDso struct {
	bytes []byte
	mime  string
}

func EmbedHandler(log lvsdk.Logger, ctx *fasthttp.RequestCtx) {
	if cache == nil {
		cache = make(map[string]*staticDso)
		mime.AddExtensionType(".map", "application/json")
	}
	path := string(ctx.Path())
	if path == "/" {
		path = "/index.html"
	}
	ext := filepath.Ext(path)
	//FIXME add client side standard http file caching
	ct := mime.TypeByExtension(ext)
	log.Trace(path, ct)
	if !strings.HasPrefix(path, "/ws/") {
		path = "build" + path
		static, ok := cache[path]
		if !ok {
			bytes, err := fs.ReadFile(path)
			if err != nil {
				log.Trace(path, err)
				ctx.Response.SetStatusCode(http.StatusNotFound)
				return
			}
			static = &staticDso{}
			cache[path] = static
			static.bytes = bytes
			static.mime = ct
		}
		ctx.SetContentType(static.mime)
		ctx.Write(static.bytes)
		return
	}
}
