package main

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

//go:embed build/*
var build embed.FS

type staticDso struct {
	bytes []byte
	mime  string
}

func NewEmbedCache(log Logger) map[string]*staticDso {
	mime.AddExtensionType(".map", "application/json")
	cache := make(map[string]*staticDso)
	fs.WalkDir(build, "build", func(path string, entry fs.DirEntry, e error) error {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			return nil
		}
		ext := filepath.Ext(path)
		ct := mime.TypeByExtension(ext)
		log.Trace("cache", path, ct)
		bytes, err := build.ReadFile(path)
		if err == nil {
			static := &staticDso{}
			static.bytes = bytes
			static.mime = ct
			cache[path] = static
		}
		return err
	})
	return cache
}

func NewEmbedHandler(log Logger, cache map[string]*staticDso) Handler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		if path == "/" {
			path = "/index.html"
		}
		static, ok := cache["build"+path]
		if !ok {
			log.Debug(path, "NF404")
			ctx.Response.SetStatusCode(http.StatusNotFound)
			return
		}
		log.Trace(path, static.mime)
		ctx.SetContentType(static.mime)
		ctx.Write(static.bytes)
	}
}
