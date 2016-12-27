package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"simas/model"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	var request model.LoginRequest
	checkError(json.NewDecoder(r.Body).Decode(&request))

	// Validate input value
	if request.Email == "" {
		panic(errors.New("Email harus diisi"))
	}

	if request.Password == "" {
		panic(errors.New("Password harus diisi"))
	}

	// Get account data from database
	account := model.Account{}
	err := handler.DB.Get(&account, "SELECT * FROM account WHERE email = ?", request.Email)
	if err != nil {
		panic(errors.New("Email tidak terdaftar"))
	}

	// Compare password with database
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.Password))
	if err != nil {
		panic(errors.New("Email dan password tidak cocok"))
	}

	// Calculate expiration time
	nbf := time.Now()
	exp := time.Now().Add(2 * time.Hour)

	if request.Remember {
		exp = time.Date(nbf.Year(), nbf.Month(), nbf.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	}

	// Generate token
	isAdmin := false
	if account.Admin == 1 {
		isAdmin = true
	}

	isPenginput := false
	if account.Penginput == 1 {
		isPenginput = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf":   nbf.Unix(),
		"exp":   exp.Unix(),
		"sub":   account.ID,
		"admin": isAdmin,
		"input": isPenginput,
	})

	tokenString, err := token.SignedString([]byte(tokenSecret))
	checkError(err)

	// Return login result
	result := model.LoginResult{
		Account: account,
		Token:   tokenString,
	}

	delay()
	w.Header().Add("Content-Type", "application/json")
	checkError(json.NewEncoder(w).Encode(&result))
}

func (handler *Handler) SelectAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse URL query
	queries := r.URL.Query()
	page := queries.Get("page")
	keyword := queries.Get("keyword")

	// Parse page number
	pageNumber, _ := strconv.Atoi(page)
	if pageNumber == 0 {
		pageNumber = 1
	}
	offset := (pageNumber - 1) * 20

	// Get max page from database
	var maxPage int
	err := handler.DB.Get(&maxPage, "SELECT FLOOR(COUNT(*) / 20) + 1 FROM account")
	checkError(err)

	// Prepare query
	sqlQuery := "SELECT id, email, nama, jabatan, telepon, admin, penginput FROM account WHERE 1"
	args := []interface{}{}

	// Add keyword to query
	if keyword != "" {
		keyword += "%"
		sqlQuery += " AND (nama LIKE ? OR email LIKE ? OR jabatan LIKE ?)"
		args = append(args, keyword, keyword, "%"+keyword)
	}

	// Add order and limit clause
	sqlQuery += " ORDER BY nama LIMIT 20 OFFSET ?"
	args = append(args, offset)

	// Select all account from database
	listAccount := []model.Account{}
	err = handler.DB.Select(&listAccount, sqlQuery, args...)
	checkError(err)

	// Encode result
	pageListAccount := model.PageListAccount{
		Page:    pageNumber,
		MaxPage: maxPage,
		Item:    listAccount,
	}

	w.Header().Add("Content-Type", "application/json")
	checkError(json.NewEncoder(w).Encode(&pageListAccount))
}

func (handler *Handler) InsertAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	var account model.Account
	checkError(json.NewDecoder(r.Body).Decode(&account))

	// Validate input
	if account.Nama == "" {
		panic(errors.New("Nama harus diisi"))
	}

	if account.Email == "" {
		panic(errors.New("Email harus diisi"))
	}

	if account.Jabatan == "" {
		panic(errors.New("Jabatan harus diisi"))
	}

	// Generate password
	randomPassword := randomString(10)
	fmt.Println(randomPassword)

	// Hash password with bcrypt
	password := []byte(randomPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
	checkError(err)

	// Insert account to database
	res := handler.DB.MustExec(`INSERT INTO account 
		(nama, email, jabatan, telepon, password, admin, penginput) VALUES 
		(?, ?, ?, ?, ?, ?, ?)`,
		account.Nama,
		account.Email,
		account.Jabatan,
		account.Telepon,
		hashedPassword,
		account.Admin,
		account.Penginput)

	// Return inserted ID
	delay()
	id, _ := res.LastInsertId()
	fmt.Fprint(w, id)
}

func (handler *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	var account model.Account
	checkError(json.NewDecoder(r.Body).Decode(&account))

	// Validate input
	if account.Nama == "" {
		panic(errors.New("Nama harus diisi"))
	}

	if account.Email == "" {
		panic(errors.New("Email harus diisi"))
	}

	if account.Jabatan == "" {
		panic(errors.New("Jabatan harus diisi"))
	}

	// Check if requested account is admin
	var admin int
	checkError(handler.DB.Get(&admin, "SELECT admin FROM account WHERE id = ?", account.ID))

	// If it is admin, check number of existing admin in database
	if admin == 1 {
		var nAdmin int
		checkError(handler.DB.Get(&nAdmin, "SELECT COUNT(*) FROM account WHERE admin = 1"))

		if nAdmin == 1 && account.Admin == 0 {
			panic(errors.New("Setidaknya harus ada satu admin"))
		}
	}

	// Update account in database
	handler.DB.MustExec(`UPDATE account SET nama = ?, email = ?, jabatan = ?, 
		telepon = ?, admin = ?, penginput = ? WHERE id = ?`,
		account.Nama,
		account.Email,
		account.Jabatan,
		account.Telepon,
		account.Admin,
		account.Penginput,
		account.ID)

	// Return updated ID
	delay()
	fmt.Fprint(w, account.ID)
}

func (handler *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get id account from URL address
	idAccount := ps.ByName("id")

	// Check if requested account is admin
	var admin int
	checkError(handler.DB.Get(&admin, "SELECT admin FROM account WHERE id = ?", idAccount))

	// If it is admin, check number of existing admin in database
	if admin == 1 {
		var nAdmin int
		checkError(handler.DB.Get(&nAdmin, "SELECT COUNT(*) FROM account WHERE admin = 1"))

		if nAdmin == 1 {
			panic(errors.New("Setidaknya harus ada satu admin"))
		}
	}

	// Delete account in database
	handler.DB.MustExec("DELETE FROM account WHERE id = ?", idAccount)

	// Return ID
	delay()
	fmt.Fprint(w, idAccount)
}
