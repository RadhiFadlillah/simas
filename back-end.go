package main

import (
	"fmt"
	"log"
	"net/http"
	"simas/handler"
	"simas/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type BackEnd struct {
	DB         *sqlx.DB
	Config     model.Configuration
	PortNumber int
}

func NewBackEnd(port int, config model.Configuration) BackEnd {
	dbSource := fmt.Sprintf("%s:%s@/%s",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseName)

	backEnd := BackEnd{
		DB:         sqlx.MustConnect("mysql", dbSource),
		Config:     config,
		PortNumber: port,
	}

	backEnd.generateAdmin()

	return backEnd
}

func (backend *BackEnd) ServeApp() {
	// Create handler
	hdl := handler.Handler{
		DB:     backend.DB,
		Config: backend.Config,
	}

	// Create router
	router := httprouter.New()

	// Handle path to UI
	router.GET("/res/*filepath", hdl.ServeFile)
	router.GET("/style/*filepath", hdl.ServeFile)
	router.GET("/", hdl.ServeIndexPage)
	router.GET("/login", hdl.ServeLoginPage)

	// Handle path to API
	router.POST("/api/login", hdl.Login)

	router.GET("/api/account", hdl.SelectAccount)
	router.PUT("/api/account", hdl.UpdateAccount)
	router.POST("/api/account", hdl.InsertAccount)
	router.POST("/api/account/password", hdl.UpdatePassword)
	router.DELETE("/api/account/:id", hdl.DeleteAccount)

	router.GET("/api/surat", hdl.SelectSurat)
	router.GET("/api/surat/id/:id", hdl.GetSurat)
	router.GET("/api/surat/image/:name", hdl.GetFileSurat)
	router.PUT("/api/surat", hdl.UpdateSurat)
	router.POST("/api/surat", hdl.InsertSurat)
	router.DELETE("/api/surat/:id", hdl.DeleteSurat)

	router.POST("/api/disposisi", hdl.InsertDisposisi)
	router.POST("/api/diarsip", hdl.InsertDiarsip)
	router.POST("/api/ditindak", hdl.InsertDitindak)

	// Set panic handler
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, arg interface{}) {
		http.Error(w, fmt.Sprint(arg), 500)
	}

	// Serve app
	log.Printf("Serve app in port %d\n", backend.PortNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", backend.PortNumber), router))
}

func (backend *BackEnd) Close() {
	backend.DB.Close()
}

func (backend *BackEnd) generateAdmin() {
	// If there are no existing account, create new admin
	var nAccount int
	err := backend.DB.Get(&nAccount, "SELECT COUNT(*) FROM account")
	checkError(err)

	if nAccount == 0 {
		password := []byte("admin")
		hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
		checkError(err)

		backend.DB.MustExec(`INSERT INTO account 
			(email, nama, password, jabatan, admin, penginput) VALUES (?, ?, ?, ?, ?, ?)`,
			"admin@simas", "Administrator", hashedPassword, "Administrator", 1, 1)
	}
}
