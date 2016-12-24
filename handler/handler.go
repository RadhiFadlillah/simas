package handler

import (
	"database/sql"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

const (
	viewDir     string = "./view"
	tokenSecret string = "mxGXCGMKuaewWjfUQbtJ6vYdHJVLkUUej2YsUnKJhMM4PTqcrvjHb7T27iAsFj4S"
	letters     string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Handler struct {
	DB *sqlx.DB
}

func checkToken(r *http.Request) (jwt.MapClaims, error) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return nil, errors.New("Token tidak tersedia")
	}

	token, err := jwt.Parse(tokenCookie.Value, jwtKeyFunc)
	if err != nil || !token.Valid {
		return nil, errors.New("Token tidak valid")
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims.Valid() != nil {
		return nil, errors.New("Token sudah expired")
	}

	return claims, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("Unexpected signing method")
	}

	return []byte(tokenSecret), nil
}

func redirectPage(w http.ResponseWriter, r *http.Request, url string) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.Redirect(w, r, url, 301)
}

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func checkError(err error) {
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
}

func delay() {
	time.Sleep(2000 * time.Millisecond)
}
