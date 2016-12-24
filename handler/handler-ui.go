package handler

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func (handler *Handler) ServeFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := viewDir + r.URL.Path

	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		http.ServeFile(w, r, path)
		return
	}

	http.NotFound(w, r)
}

func (handler *Handler) ServeIndexPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := checkToken(r); err != nil {
		redirectPage(w, r, "/login")
	}

	path := viewDir + "/index.html"
	http.ServeFile(w, r, path)
}

func (handler *Handler) ServeLoginPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := checkToken(r); err == nil {
		redirectPage(w, r, "/")
	}

	path := viewDir + "/login.html"
	http.ServeFile(w, r, path)
}
