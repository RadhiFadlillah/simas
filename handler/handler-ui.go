package handler

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

func (handler *Handler) ServeFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Read asset path
	path := r.URL.Path
	if path[0:1] == "/" {
		path = path[1:]
	}

	// Load asset
	asset, err := Asset(path)
	checkError(err)

	// Set response header content type
	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}

	// Serve asset
	buffer := bytes.NewBuffer(asset)
	io.Copy(w, buffer)
}

func (handler *Handler) ServeIndexPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := checkToken(r); err != nil {
		redirectPage(w, r, "/login")
		return
	}

	// Load asset
	path := "index.html"
	asset, err := Asset(path)
	checkError(err)

	// Serve asset
	w.Header().Set("Content-Type", "text/html")
	buffer := bytes.NewBuffer(asset)
	io.Copy(w, buffer)
}

func (handler *Handler) ServeLoginPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := checkToken(r); err == nil {
		redirectPage(w, r, "/")
		return
	}

	// Load asset
	path := "login.html"
	asset, err := Asset(path)
	checkError(err)

	// Serve asset
	w.Header().Set("Content-Type", "text/html")
	buffer := bytes.NewBuffer(asset)
	io.Copy(w, buffer)
}
